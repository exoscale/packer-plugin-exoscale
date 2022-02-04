//go:generate packer-sdc mapstructure-to-hcl2 -type Config

package exoscaleimport

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	pkrconfig "github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

const (
	defaultAPIEnvironment   = "api"
	defaultTemplateBootMode = "legacy"
)

type Config struct {
	SOSEndpoint             string `mapstructure:"sos_endpoint"`
	APIEnvironment          string `mapstructure:"api_environment"`
	APIKey                  string `mapstructure:"api_key"`
	APISecret               string `mapstructure:"api_secret"`
	APITimeout              int64  `mapstructure:"api_timeout"`
	ImageBucket             string `mapstructure:"image_bucket"`
	TemplateZone            string `mapstructure:"template_zone"`
	TemplateName            string `mapstructure:"template_name"`
	TemplateDescription     string `mapstructure:"template_description"`
	TemplateUsername        string `mapstructure:"template_username"`
	TemplateBootMode        string `mapstructure:"template_boot_mode"`
	TemplateDisablePassword bool   `mapstructure:"template_disable_password"`
	TemplateDisableSSHKey   bool   `mapstructure:"template_disable_sshkey"`
	SkipClean               bool   `mapstructure:"skip_clean"`

	ctx interpolate.Context

	common.PackerConfig `mapstructure:",squash"`
}

func NewConfig(raws ...interface{}) (*Config, error) {
	var config Config

	err := pkrconfig.Decode(&config, &pkrconfig.DecodeOpts{
		PluginType:         BuilderId,
		Interpolate:        true,
		InterpolateContext: &config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, raws...)
	if err != nil {
		return nil, err
	}

	requiredArgs := map[string]interface{}{
		"api_key":       config.APIKey,
		"api_secret":    config.APISecret,
		"image_bucket":  config.ImageBucket,
		"template_name": config.TemplateName,
		"template_zone": config.TemplateZone,
	}

	errs := new(packer.MultiError)
	for k, v := range requiredArgs {
		if reflect.ValueOf(v).IsZero() {
			errs = packer.MultiErrorAppend(errs, fmt.Errorf("%s must be set", k))
		}
	}

	if len(errs.Errors) > 0 {
		return nil, errs
	}

	if config.APIEnvironment == "" {
		config.APIEnvironment = defaultAPIEnvironment
	}

	// Template registration can take a _long time_, set the default
	// Exoscale API client timeout to 1h as a precaution.
	if config.APITimeout == 0 {
		config.APITimeout = 3600
	}

	if config.TemplateBootMode == "" {
		config.TemplateBootMode = defaultTemplateBootMode
	}

	if config.SOSEndpoint == "" {
		config.SOSEndpoint = "https://sos-" + config.TemplateZone + ".exo.io"
	}

	return &config, nil
}

// ConfigSpec returns HCL object spec
func (p *PostProcessor) ConfigSpec() hcldec.ObjectSpec {
	return p.config.FlatMapstructure().HCL2Spec()
}
