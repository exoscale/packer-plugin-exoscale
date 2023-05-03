package exoscale

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/stretchr/testify/require"
)

var (
	testAccSnapshotDownload     = true
	testAccSnapshotDownloadPath = os.TempDir()
	testAccTemplateName         = "packer-plugin-test-" + new(testSuite).randomString(6)
	testAccTemplateZones        = []string{"ch-gva-2", "ch-dk-2"}
	testAccTemplateDescription  = new(testSuite).randomString(6)
	testAccTemplateUsername     = "packer"
	testAccTemplateMaintainer   = "Exoscale"
	testAccTemplateVersion      = "0.acceptance"
	testAccTemplateBuild        = new(testSuite).randomString(8)

	testAccSnapshotDownloadFile = filepath.Join(testAccSnapshotDownloadPath, testAccTemplateName+".qcow2")
	testAccSnapshotChecksumFile = filepath.Join(testAccSnapshotDownloadPath, testAccTemplateName+".md5sum")
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

		"instance_template":  "Linux Ubuntu 22.04 LTS 64-bit",
		"instance_disk_size": 10,
		"ssh_username":       "ubuntu",

		"snapshot_download":      testAccSnapshotDownload,
		"snapshot_download_path": testAccSnapshotDownloadPath,

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

	require.FileExists(t, testAccSnapshotDownloadFile)
	require.FileExists(t, testAccSnapshotChecksumFile)
	md5sumContent, err := os.ReadFile(testAccSnapshotChecksumFile)
	require.NoError(t, err)
	md5sumActual, err := md5sum(testAccSnapshotDownloadFile)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%s *%s.qcow2", md5sumActual, testAccTemplateName), string(md5sumContent))
	// Leave this material be in the download path

	require.NoError(t, artifact.Destroy())
}

func md5sum(filePath string) (string, error) {
	var md5sum string
	file, err := os.Open(filePath)
	if err != nil {
		return md5sum, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return md5sum, err
	}
	hashBytes := hash.Sum(nil)[:16]
	md5sum = hex.EncodeToString(hashBytes)
	return md5sum, nil
}
