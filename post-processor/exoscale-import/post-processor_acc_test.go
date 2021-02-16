package exoscaleimport

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
	testAccImageBucket         = "eat-template-images"
	testAccTemplateName        = "test-packer-builder-exoscale"
	testAccTemplateZone        = "ch-dk-2"
	testAccTemplateDescription = "Built with Packer"
	testAccTemplateUsername    = "packer"
	testAccImageFile           = "./testdata/test-packer-post-processor-exoscale-import.qcow2"
)

type testMockArtifact struct {
	files []string
}

func (a *testMockArtifact) BuilderId() string          { return qemuBuilderID }
func (a *testMockArtifact) Files() []string            { return a.files }
func (a *testMockArtifact) Id() string                 { return "" }
func (a *testMockArtifact) String() string             { return "" }
func (a *testMockArtifact) State(_ string) interface{} { return nil }
func (a *testMockArtifact) Destroy() error             { return nil }

func TestAccPostProcessor(t *testing.T) {
	var postProcessor PostProcessor

	if v := os.Getenv(acctest.TestEnvVar); v == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", acctest.TestEnvVar))
		return
	}

	require.NotEmpty(t, os.Getenv("EXOSCALE_API_KEY"),
		"EXOSCALE_API_KEY environment variable must be set for acceptance tests")
	require.NotEmpty(t, os.Getenv("EXOSCALE_API_SECRET"),
		"EXOSCALE_API_SECRET environment variable must be set for acceptance tests")

	err := postProcessor.Configure([]interface{}{map[string]interface{}{
		"api_key":              os.Getenv("EXOSCALE_API_KEY"),
		"api_secret":           os.Getenv("EXOSCALE_API_SECRET"),
		"image_bucket":         testAccImageBucket,
		"template_zone":        testAccTemplateZone,
		"template_name":        testAccTemplateName,
		"template_description": testAccTemplateDescription,
		"template_username":    testAccTemplateUsername,
	}}...)
	require.NoError(t, err)

	artifact, _, _, err := postProcessor.PostProcess(
		context.Background(),
		packer.TestUi(t),
		&testMockArtifact{files: []string{testAccImageFile}})
	require.NoError(t, err)
	require.NotNil(t, artifact)

	a := artifact.(*Artifact)
	require.NotNil(t, a.template.ID)
	require.Equal(t, testAccTemplateZone, a.template.ZoneName)
	require.Equal(t, testAccTemplateName, a.template.Name)
	require.Equal(t, testAccTemplateDescription, a.template.DisplayText)
	require.Equal(t, defaultTemplateBootMode, a.template.BootMode)
	require.Equal(t, testAccTemplateUsername, a.template.Details["username"])

	require.NoError(t, artifact.Destroy())
}
