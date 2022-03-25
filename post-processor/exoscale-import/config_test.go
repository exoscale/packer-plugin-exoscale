package exoscaleimport

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testConfigAPIKey        = "EXOabcdef0123456789abcdef01"
	testConfigAPISecret     = "ABCDEFGHIJKLMNOPRQSTUVWXYZ0123456789abcdefg"
	testConfigImageBucket   = "test-template-images"
	testConfigTemplateZones = []string{"ch-gva-2", "ch-dk-2"}
	testConfigTemplateName  = "test-packer"
	// Deprecated
	testConfigTemplateZone  = "ch-dk-2"
)

func TestNewConfig(t *testing.T) {
	_, err := NewConfig()
	require.Error(t, err, "incomplete configuration should return an error")

	config, err := NewConfig([]interface{}{map[string]interface{}{
		// Minimal configuration
		"api_key":        testConfigAPIKey,
		"api_secret":     testConfigAPISecret,
		"image_bucket":   testConfigImageBucket,
		"template_name":  testConfigTemplateName,
		"template_zones": testConfigTemplateZones,
	}}...)
	require.NoError(t, err)
	require.NotNil(t, config)
	require.Equal(t, defaultAPIEnvironment, config.APIEnvironment)
	require.Equal(t, testConfigTemplateZones[0], config.ImageZone)
	require.Equal(t, "https://sos-"+testConfigTemplateZones[0]+".exo.io", config.SOSEndpoint)
	require.Equal(t, defaultTemplateBootMode, config.TemplateBootMode)
}

func TestNewConfigDeprecated(t *testing.T) {
	config, err := NewConfig([]interface{}{map[string]interface{}{
		// Minimal configuration
		"api_key":       testConfigAPIKey,
		"api_secret":    testConfigAPISecret,
		"image_bucket":  testConfigImageBucket,
		"template_name": testConfigTemplateName,
		// Deprecated
		"template_zone": testConfigTemplateZone,
	}}...)
	require.NoError(t, err)
	require.NotNil(t, config)
	require.Equal(t, 1, len(config.TemplateZones))
	require.Equal(t, testConfigTemplateZone, config.TemplateZones[0])
	require.Equal(t, testConfigTemplateZone, config.ImageZone)
	require.Equal(t, "https://sos-"+testConfigTemplateZone+".exo.io", config.SOSEndpoint)

	config, err = NewConfig([]interface{}{map[string]interface{}{
		// Minimal configuration
		"api_key":        testConfigAPIKey,
		"api_secret":     testConfigAPISecret,
		"image_bucket":   testConfigImageBucket,
		"template_name":  testConfigTemplateName,
		"template_zones": testConfigTemplateZones,
		// Deprecated
		"template_zone": testConfigTemplateZone,
	}}...)
	require.NoError(t, err)
	require.NotNil(t, config)
	require.Equal(t, testConfigTemplateZones, config.TemplateZones)
}
