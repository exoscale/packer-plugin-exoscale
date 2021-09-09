package exoscaleimport

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testConfigAPIKey       = "EXOabcdef0123456789abcdef01"
	testConfigAPISecret    = "ABCDEFGHIJKLMNOPRQSTUVWXYZ0123456789abcdefg"
	testConfigImageBucket  = "test-template-images"
	testConfigTemplateZone = "ch-gva-2"
	testConfigTemplateName = "test-packer"
)

func TestNewConfig(t *testing.T) {
	_, err := NewConfig()
	require.Error(t, err, "incomplete configuration should return an error")

	// Minimal configuration
	config, err := NewConfig([]interface{}{map[string]interface{}{
		"api_key":       testConfigAPIKey,
		"api_secret":    testConfigAPISecret,
		"image_bucket":  testConfigImageBucket,
		"template_name": testConfigTemplateName,
		"template_zone": testConfigTemplateZone,
	}}...)
	require.NoError(t, err)
	require.NotNil(t, config)
	require.Equal(t, defaultAPIEnvironment, config.APIEnvironment)
	require.Equal(t, "https://sos-"+config.TemplateZone+".exo.io", config.SOSEndpoint)
	require.Equal(t, defaultTemplateBootMode, config.TemplateBootMode)
}
