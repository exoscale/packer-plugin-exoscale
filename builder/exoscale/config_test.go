package exoscale

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testConfigAPIKey           = "EXOabcdef0123456789abcdef01"
	testConfigAPISecret        = "ABCDEFGHIJKLMNOPRQSTUVWXYZ0123456789abcdefg"
	testConfigInstanceTemplate = "Linux Ubuntu 20.04 LTS 64-bit"
	testConfigTemplateZone     = "ch-gva-2"
	testConfigTemplateName     = "test-packer"
	testConfigSSHUsername      = "ubuntu"
)

func TestNewConfig(t *testing.T) {
	_, err := NewConfig()
	require.Error(t, err, "incomplete configuration should return an error")

	// Minimal configuration
	config, err := NewConfig([]interface{}{map[string]interface{}{
		"api_key":           testConfigAPIKey,
		"api_secret":        testConfigAPISecret,
		"instance_template": testConfigInstanceTemplate,
		"template_name":     testConfigTemplateName,
		"template_zone":     testConfigTemplateZone,
		"ssh_username":      testConfigSSHUsername,
	}}...)
	require.NoError(t, err)
	require.NotNil(t, config)
	require.Equal(t, defaultAPIEndpoint, config.APIEndpoint)
	require.Equal(t, defaultInstanceType, config.InstanceType)
	require.Equal(t, defaultInstanceDiskSize, config.InstanceDiskSize)
	require.Equal(t, []string{defaultInstanceSecurityGroup}, config.InstanceSecurityGroups)
	require.Equal(t, defaultInstanceTemplateFilter, config.InstanceTemplateFilter)
	require.Equal(t, config.InstanceZone, testConfigTemplateZone)
	require.Equal(t, defaultTemplateBootMode, config.TemplateBootMode)
}
