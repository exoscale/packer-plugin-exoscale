package exoscale

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

var (
	testAccTemplateName        = "test-packer-builder-exoscale"
	testAccTemplateZone        = "ch-gva-2"
	testAccTemplateDescription = "Built with Packer"
	testAccTemplateUsername    = "packer"
	testAccBuilderTemplate     = fmt.Sprintf(`{
  "variables": {
    "api_key": "{{env `+"`EXOSCALE_API_KEY`"+`}}",
    "api_secret": "{{env `+"`EXOSCALE_API_SECRET`"+`}}"
  },

  "builders": [{
    "type": "test",
    "api_key": "{{user `+"`api_key`"+`}}",
    "api_secret": "{{user `+"`api_secret`"+`}}",
    "instance_template": "Linux Ubuntu 20.04 LTS 64-bit",
    "template_zone": "%s",
    "template_name": "%s",
    "template_description": "%s",
    "template_username": "%s",
    "ssh_username": "ubuntu"
  }]
}`,
		testAccTemplateZone,
		testAccTemplateName,
		testAccTemplateDescription,
		testAccTemplateUsername,
	)
)

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("EXOSCALE_API_KEY"); v == "" {
		t.Fatal("EXOSCALE_API_KEY must be set for acceptance tests")
	}

	if v := os.Getenv("EXOSCALE_API_SECRET"); v == "" {
		t.Fatal("EXOSCALE_API_SECRET must be set for acceptance tests")
	}
}

func testAccCheckBuilderArtifact(artifacts []packer.Artifact) error {
	if len(artifacts) != 1 {
		return fmt.Errorf("expected 1 artifact, got %d", len(artifacts))
	}
	artifact := artifacts[0].(*Artifact)

	if artifact.template.ID == nil {
		return fmt.Errorf("artifact template ID is not set")
	}

	if artifact.template.Name != testAccTemplateName {
		return fmt.Errorf("expected template name %q, got %q",
			testAccTemplateName,
			artifact.template.Name)
	}

	if artifact.template.ZoneName != testAccTemplateZone {
		return fmt.Errorf("expected template zone %q, got %q",
			testAccTemplateZone,
			artifact.template.ZoneName)
	}

	if artifact.template.BootMode != defaultTemplateBootMode {
		return fmt.Errorf("expected template boot mode %q, got %q",
			defaultTemplateBootMode,
			artifact.template.BootMode)
	}

	if username, ok := artifact.template.Details["username"]; !ok {
		return errors.New("artifact username not set")
	} else if username != testAccTemplateUsername {
		return fmt.Errorf("expected template username %q, got %q",
			testAccTemplateUsername,
			username)
	}

	return nil
}

func TestAccBuilder(t *testing.T) {
	acctest.Test(t, acctest.TestCase{
		Builder:  &Builder{},
		Template: testAccBuilderTemplate,
		PreCheck: func() { testAccPreCheck(t) },
		Check:    testAccCheckBuilderArtifact,
	})
}
