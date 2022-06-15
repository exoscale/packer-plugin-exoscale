package exoscale

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/stretchr/testify/require"
)

var (
	testAccTemplateName        = "packer-plugin-test-" + new(testSuite).randomString(6)
	testAccTemplateZones       = []string{"ch-gva-2", "ch-dk-2"}
	testAccTemplateDescription = new(testSuite).randomString(6)
	testAccTemplateUsername    = "packer"
	testAccTemplateMaintainer  = "Exoscale"
	testAccTemplateVersion     = "0.acceptance"
	testAccTemplateBuild       = new(testSuite).randomString(8)
)

func TestAccBuilder(t *testing.T) {
	var builder Builder

	if v := os.Getenv(acctest.TestEnvVar); v == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", acctest.TestEnvVar))
		return
	}

	require.NotEmpty(t, os.Getenv("EXOSCALE_API_KEY"),
		"EXOSCALE_API_KEY environment variable must be set for acceptance tests")
	require.NotEmpty(t, os.Getenv("EXOSCALE_API_SECRET"),
		"EXOSCALE_API_SECRET environment variable must be set for acceptance tests")

	_, _, err := builder.Prepare([]interface{}{map[string]interface{}{
		"api_key":    os.Getenv("EXOSCALE_API_KEY"),
		"api_secret": os.Getenv("EXOSCALE_API_SECRET"),

		"instance_template":  "Linux Ubuntu 20.04 LTS 64-bit",
		"instance_disk_size": 10,
		"ssh_username":       "ubuntu",

		"template_zones":       testAccTemplateZones,
		"template_name":        testAccTemplateName,
		"template_description": testAccTemplateDescription,
		"template_username":    testAccTemplateUsername,
		"template_maintainer":  testAccTemplateMaintainer,
		"template_version":     testAccTemplateVersion,
		"template_build":       testAccTemplateBuild,
	}}...)
	require.NoError(t, err)

	artifact, err := builder.Run(context.Background(), packer.TestUi(t), &packer.MockHook{})
	require.NoError(t, err)
	require.NotNil(t, artifact)

	a := artifact.(*Artifact)
	require.Equal(t, len(a.templates), len(testAccTemplateZones))
	for _, template := range a.templates {
		require.NotNil(t, template.ID)
		require.Equal(t, testAccTemplateName, *template.Name)
		require.Equal(t, testAccTemplateDescription, *template.Description)
		require.Equal(t, defaultTemplateBootMode, *template.BootMode)
		require.Equal(t, testAccTemplateUsername, *template.DefaultUser)
		require.Equal(t, testAccTemplateMaintainer, *template.Maintainer)
		require.Equal(t, testAccTemplateVersion, *template.Version)
		require.Equal(t, testAccTemplateBuild, *template.Build)
	}

	require.NoError(t, artifact.Destroy())
}
