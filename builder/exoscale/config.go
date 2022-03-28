//go:generate packer-sdc mapstructure-to-hcl2 -type Config

package exoscale

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	pkrconfig "github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

const (
	defaultAPIEnvironment                   = "api"
	defaultInstanceType                     = "medium"
	defaultInstanceDiskSize           int64 = 50
	defaultInstanceSecurityGroup            = "default"
	defaultInstanceTemplateVisibility       = "public"
	defaultTemplateBootMode                 = "legacy"
)

type Config struct {
	APIEnvironment             string   `mapstructure:"api_environment"`
	APIKey                     string   `mapstructure:"api_key"`
	APISecret                  string   `mapstructure:"api_secret"`
	APITimeout                 int64    `mapstructure:"api_timeout"`
	InstanceName               string   `mapstructure:"instance_name"`
	InstanceZone               string   `mapstructure:"instance_zone"`
	InstanceTemplate           string   `mapstructure:"instance_template"`
	InstanceTemplateVisibility string   `mapstructure:"instance_template_visibility"`
	InstanceType               string   `mapstructure:"instance_type"`
	InstanceDiskSize           int64    `mapstructure:"instance_disk_size"`
	InstanceSecurityGroups     []string `mapstructure:"instance_security_groups"`
	InstancePrivateNetworks    []string `mapstructure:"instance_private_networks"`
	InstanceSSHKey             string   `mapstructure:"instance_ssh_key"`
	TemplateZones              []string `mapstructure:"template_zones"`
	TemplateName               string   `mapstructure:"template_name"`
	TemplateDescription        string   `mapstructure:"template_description"`
	TemplateUsername           string   `mapstructure:"template_username"`
	TemplateBootMode           string   `mapstructure:"template_boot_mode"`
	TemplateDisablePassword    bool     `mapstructure:"template_disable_password"`
	TemplateDisableSSHKey      bool     `mapstructure:"template_disable_sshkey"`
	// Deprecated
	TemplateZone string `mapstructure:"template_zone"`

	ctx interpolate.Context

	common.PackerConfig `mapstructure:",squash"`
	Comm                communicator.Config `mapstructure:",squash"`
}

func NewConfig(raws ...interface{}) (*Config, []string, error) {
	var config Config

	err := pkrconfig.Decode(
		&config,
		&pkrconfig.DecodeOpts{
			Interpolate:        true,
			InterpolateContext: &config.ctx,
			InterpolateFilter: &interpolate.RenderFilter{
				Exclude: []string{},
			},
		},
		raws...)
	if err != nil {
		return nil, nil, err
	}

	// Deprecated arguments
	warnings := []string{}
	// (template_zones <-> template_zone)
	if config.TemplateZone != "" {
		if len(config.TemplateZones) == 0 {
			config.TemplateZones = []string{config.TemplateZone}
		} else {
			warnings = append(warnings, "Both template_zones and template_zone are defined; ignoring the latter")
		}
	}

	requiredArgs := map[string]interface{}{
		"api_key":           config.APIKey,
		"api_secret":        config.APISecret,
		"instance_template": config.InstanceTemplate,
		"template_name":     config.TemplateName,
		"template_zones":    config.TemplateZones,
	}

	errs := new(packer.MultiError)
	for k, v := range requiredArgs {
		if reflect.ValueOf(v).IsZero() || reflect.ValueOf(v).Len() == 0 {
			errs = packer.MultiErrorAppend(errs, fmt.Errorf("%s must be set", k))
		}
	}

	if es := config.Comm.Prepare(&config.ctx); len(es) > 0 {
		errs = packer.MultiErrorAppend(errs, es...)
	}

	if len(errs.Errors) > 0 {
		return nil, warnings, errs
	}

	if config.APITimeout == 0 {
		config.APITimeout = 3600
	}

	if config.InstanceZone == "" {
		config.InstanceZone = config.TemplateZones[0]
	}

	if config.TemplateBootMode == "" {
		config.TemplateBootMode = defaultTemplateBootMode
	}

	if config.APIEnvironment == "" {
		config.APIEnvironment = defaultAPIEnvironment
	}

	if config.InstanceType == "" {
		config.InstanceType = defaultInstanceType
	}

	if config.InstanceTemplateVisibility == "" {
		config.InstanceTemplateVisibility = defaultInstanceTemplateVisibility
	}

	if config.InstanceDiskSize == 0 {
		config.InstanceDiskSize = defaultInstanceDiskSize
	}

	if len(config.InstanceSecurityGroups) == 0 {
		config.InstanceSecurityGroups = []string{defaultInstanceSecurityGroup}
	}

	return &config, warnings, nil
}

// ConfigSpec returns HCL object spec
func (b *Builder) ConfigSpec() hcldec.ObjectSpec {
	return b.config.FlatMapstructure().HCL2Spec()
}
