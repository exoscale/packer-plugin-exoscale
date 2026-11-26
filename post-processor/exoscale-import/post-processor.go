package exoscaleimport

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/version"

	egoscale "github.com/exoscale/egoscale/v2"
	exoapi "github.com/exoscale/egoscale/v2/api"
	egoscalev3 "github.com/exoscale/egoscale/v3"
	exocredentials "github.com/exoscale/egoscale/v3/credentials"
)

const (
	qemuBuilderID     = "transcend.qemu"
	fileBuilderID     = "packer.file"
	artificeBuilderID = "packer.post-processor.artifice"
)

var userAgent = "Exoscale-Packer-Post-Processor/" + version.SDKVersion.FormattedVersion()

func init() {
	egoscale.UserAgent = userAgent + " " + egoscale.UserAgent
}

type exoscaleClient interface {
	CopyTemplate(context.Context, string, *egoscale.Template, string) (*egoscale.Template, error)
	DeleteTemplate(context.Context, string, *egoscale.Template) error
	GetTemplate(context.Context, string, string) (*egoscale.Template, error)
}

type s3Client interface {
	s3manager.UploadAPIClient

	DeleteObject(context.Context, *s3.DeleteObjectInput, ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

type PostProcessor struct {
	config       *Config
	runner       multistep.Runner
	exo          exoscaleClient
	registration *egoscalev3.Client
	sos          s3Client
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	config, err := NewConfig(raws...)
	if err != nil {
		return err
	}
	p.config = config

	packer.LogSecretFilter.Set(p.config.APIKey, p.config.APISecret)

	return nil
}

func (p *PostProcessor) newSOSClient(ctx context.Context) (*s3.Client, error) {
	cfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(p.config.ImageZone),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			p.config.APIKey,
			p.config.APISecret,
			"")),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(p.config.SOSEndpoint)
		o.UsePathStyle = true
		// Preserve the existing SOS request format instead of adding checksum trailers.
		o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
	}), nil
}

func (p *PostProcessor) PostProcess(
	ctx context.Context,
	ui packer.Ui,
	a packer.Artifact,
) (packer.Artifact, bool, bool, error) {
	switch a.BuilderId() {
	case qemuBuilderID, fileBuilderID, artificeBuilderID:
		break
	default:
		err := fmt.Errorf("unsupported artifact type %q: this post-processor only supports "+
			"artifacts from QEMU/file builders and Artifice post-processor", a.BuilderId())
		return nil, false, false, err
	}

	exo, err := egoscale.NewClient(
		p.config.APIKey,
		p.config.APISecret,
		egoscale.ClientOptWithTimeout(time.Duration(p.config.APITimeout*int64(time.Second))),
	)
	if err != nil {
		return nil, false, false, fmt.Errorf("unable to initialize Exoscale client: %w", err)
	}
	p.exo = exo

	registerEndpoint := exoapi.NewReqEndpoint(p.config.APIEnvironment, p.config.TemplateZones[0])
	endpoint := url.URL{Scheme: "https", Host: registerEndpoint.Host(), Path: "/" + exoapi.Prefix}
	p.registration, err = egoscalev3.NewClient(
		exocredentials.NewStaticCredentials(p.config.APIKey, p.config.APISecret),
		egoscalev3.ClientOptWithEndpoint(egoscalev3.Endpoint(endpoint.String())),
		egoscalev3.ClientOptWithUserAgent(userAgent),
	)
	if err != nil {
		return nil, false, false, fmt.Errorf("unable to initialize template registration client: %w", err)
	}

	sos, err := p.newSOSClient(ctx)
	if err != nil {
		return nil, false, false, fmt.Errorf("unable to initialize SOS client: %w", err)
	}
	p.sos = sos

	state := new(multistep.BasicStateBag)
	state.Put("ui", ui)
	state.Put("artifact", a)
	// A single template is registered (once) in the first zone and then copied as many times
	// as they are additional zones. We must keep track of each template-zone, as final artifact
	state.Put("templates", []*egoscale.Template{})

	steps := []multistep.Step{
		&stepUploadImage{postProcessor: p},
		&stepRegisterTemplate{postProcessor: p},
		&stepCopyTemplate{postProcessor: p},
	}

	ctx = exoapi.WithEndpoint(ctx, exoapi.NewReqEndpoint(p.config.APIEnvironment, p.config.ImageZone))

	p.runner = commonsteps.NewRunnerWithPauseFn(steps, p.config.PackerConfig, ui, state)
	p.runner.Run(ctx, state)

	if rawErr, ok := state.GetOk("error"); ok {
		return nil, false, false, rawErr.(error)
	}

	if _, ok := state.GetOk(multistep.StateCancelled); ok {
		return nil, false, false, errors.New("post-processing cancelled")
	}

	if _, ok := state.GetOk(multistep.StateHalted); ok {
		return nil, false, false, errors.New("post-processing halted")
	}

	templates, ok := state.GetOk("templates")
	if !ok {
		return nil, false, false, errors.New("unable to find templates in state")
	}

	return &Artifact{
		StateData: map[string]interface{}{"generated_data": state.Get("generated_data")},

		postProcessor: p,
		state:         state,
		templates:     templates.([]*egoscale.Template),
	}, false, false, nil
}
