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
	testAccTemplateName        = "test-packer-builder-exoscale"
	testAccTemplateZone        = "ch-gva-2"
	testAccTemplateDescription = "Built with Packer"
	testAccTemplateUsername    = "packer"
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

		"template_zone":        testAccTemplateZone,
		"template_name":        testAccTemplateName,
		"template_description": testAccTemplateDescription,
		"template_username":    testAccTemplateUsername,
	}}...)
	require.NoError(t, err)

	artifact, err := builder.Run(context.Background(), packer.TestUi(t), &packer.MockHook{})
	require.NoError(t, err)
	require.NotNil(t, artifact)

	a := artifact.(*Artifact)
	require.NotNil(t, a.template.ID)
	require.Equal(t, testAccTemplateName, *a.template.Name)
	require.Equal(t, testAccTemplateDescription, *a.template.Description)
	require.Equal(t, defaultTemplateBootMode, *a.template.BootMode)
	require.Equal(t, testAccTemplateUsername, *a.template.DefaultUser)

	require.NoError(t, artifact.Destroy())
}
