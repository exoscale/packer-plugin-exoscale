// Code generated by "packer-sdc mapstructure-to-hcl2"; DO NOT EDIT.

package exoscaleimport

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

// FlatConfig is an auto-generated flat version of Config.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatConfig struct {
	SOSEndpoint             *string           `mapstructure:"sos_endpoint" cty:"sos_endpoint" hcl:"sos_endpoint"`
	APIEnvironment          *string           `mapstructure:"api_environment" cty:"api_environment" hcl:"api_environment"`
	APIKey                  *string           `mapstructure:"api_key" cty:"api_key" hcl:"api_key"`
	APISecret               *string           `mapstructure:"api_secret" cty:"api_secret" hcl:"api_secret"`
	APITimeout              *int64            `mapstructure:"api_timeout" cty:"api_timeout" hcl:"api_timeout"`
	ImageBucket             *string           `mapstructure:"image_bucket" cty:"image_bucket" hcl:"image_bucket"`
	ImageZone               *string           `mapstructure:"image_zone" cty:"image_zone" hcl:"image_zone"`
	TemplateZones           []string          `mapstructure:"template_zones" cty:"template_zones" hcl:"template_zones"`
	TemplateName            *string           `mapstructure:"template_name" cty:"template_name" hcl:"template_name"`
	TemplateDescription     *string           `mapstructure:"template_description" cty:"template_description" hcl:"template_description"`
	TemplateUsername        *string           `mapstructure:"template_username" cty:"template_username" hcl:"template_username"`
	TemplateBootMode        *string           `mapstructure:"template_boot_mode" cty:"template_boot_mode" hcl:"template_boot_mode"`
	TemplateDisablePassword *bool             `mapstructure:"template_disable_password" cty:"template_disable_password" hcl:"template_disable_password"`
	TemplateDisableSSHKey   *bool             `mapstructure:"template_disable_sshkey" cty:"template_disable_sshkey" hcl:"template_disable_sshkey"`
	TemplateMaintainer      *string           `mapstructure:"template_maintainer" cty:"template_maintainer" hcl:"template_maintainer"`
	TemplateVersion         *string           `mapstructure:"template_version" cty:"template_version" hcl:"template_version"`
	TemplateBuild           *string           `mapstructure:"template_build" cty:"template_build" hcl:"template_build"`
	SkipClean               *bool             `mapstructure:"skip_clean" cty:"skip_clean" hcl:"skip_clean"`
	TemplateZone            *string           `mapstructure:"template_zone" cty:"template_zone" hcl:"template_zone"`
	PackerBuildName         *string           `mapstructure:"packer_build_name" cty:"packer_build_name" hcl:"packer_build_name"`
	PackerBuilderType       *string           `mapstructure:"packer_builder_type" cty:"packer_builder_type" hcl:"packer_builder_type"`
	PackerCoreVersion       *string           `mapstructure:"packer_core_version" cty:"packer_core_version" hcl:"packer_core_version"`
	PackerDebug             *bool             `mapstructure:"packer_debug" cty:"packer_debug" hcl:"packer_debug"`
	PackerForce             *bool             `mapstructure:"packer_force" cty:"packer_force" hcl:"packer_force"`
	PackerOnError           *string           `mapstructure:"packer_on_error" cty:"packer_on_error" hcl:"packer_on_error"`
	PackerUserVars          map[string]string `mapstructure:"packer_user_variables" cty:"packer_user_variables" hcl:"packer_user_variables"`
	PackerSensitiveVars     []string          `mapstructure:"packer_sensitive_variables" cty:"packer_sensitive_variables" hcl:"packer_sensitive_variables"`
}

// FlatMapstructure returns a new FlatConfig.
// FlatConfig is an auto-generated flat version of Config.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*Config) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatConfig)
}

// HCL2Spec returns the hcl spec of a Config.
// This spec is used by HCL to read the fields of Config.
// The decoded values from this spec will then be applied to a FlatConfig.
func (*FlatConfig) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"sos_endpoint":               &hcldec.AttrSpec{Name: "sos_endpoint", Type: cty.String, Required: false},
		"api_environment":            &hcldec.AttrSpec{Name: "api_environment", Type: cty.String, Required: false},
		"api_key":                    &hcldec.AttrSpec{Name: "api_key", Type: cty.String, Required: false},
		"api_secret":                 &hcldec.AttrSpec{Name: "api_secret", Type: cty.String, Required: false},
		"api_timeout":                &hcldec.AttrSpec{Name: "api_timeout", Type: cty.Number, Required: false},
		"image_bucket":               &hcldec.AttrSpec{Name: "image_bucket", Type: cty.String, Required: false},
		"image_zone":                 &hcldec.AttrSpec{Name: "image_zone", Type: cty.String, Required: false},
		"template_zones":             &hcldec.AttrSpec{Name: "template_zones", Type: cty.List(cty.String), Required: false},
		"template_name":              &hcldec.AttrSpec{Name: "template_name", Type: cty.String, Required: false},
		"template_description":       &hcldec.AttrSpec{Name: "template_description", Type: cty.String, Required: false},
		"template_username":          &hcldec.AttrSpec{Name: "template_username", Type: cty.String, Required: false},
		"template_boot_mode":         &hcldec.AttrSpec{Name: "template_boot_mode", Type: cty.String, Required: false},
		"template_disable_password":  &hcldec.AttrSpec{Name: "template_disable_password", Type: cty.Bool, Required: false},
		"template_disable_sshkey":    &hcldec.AttrSpec{Name: "template_disable_sshkey", Type: cty.Bool, Required: false},
		"template_maintainer":        &hcldec.AttrSpec{Name: "template_maintainer", Type: cty.String, Required: false},
		"template_version":           &hcldec.AttrSpec{Name: "template_version", Type: cty.String, Required: false},
		"template_build":             &hcldec.AttrSpec{Name: "template_build", Type: cty.String, Required: false},
		"skip_clean":                 &hcldec.AttrSpec{Name: "skip_clean", Type: cty.Bool, Required: false},
		"template_zone":              &hcldec.AttrSpec{Name: "template_zone", Type: cty.String, Required: false},
		"packer_build_name":          &hcldec.AttrSpec{Name: "packer_build_name", Type: cty.String, Required: false},
		"packer_builder_type":        &hcldec.AttrSpec{Name: "packer_builder_type", Type: cty.String, Required: false},
		"packer_core_version":        &hcldec.AttrSpec{Name: "packer_core_version", Type: cty.String, Required: false},
		"packer_debug":               &hcldec.AttrSpec{Name: "packer_debug", Type: cty.Bool, Required: false},
		"packer_force":               &hcldec.AttrSpec{Name: "packer_force", Type: cty.Bool, Required: false},
		"packer_on_error":            &hcldec.AttrSpec{Name: "packer_on_error", Type: cty.String, Required: false},
		"packer_user_variables":      &hcldec.AttrSpec{Name: "packer_user_variables", Type: cty.Map(cty.String), Required: false},
		"packer_sensitive_variables": &hcldec.AttrSpec{Name: "packer_sensitive_variables", Type: cty.List(cty.String), Required: false},
	}
	return s
}
