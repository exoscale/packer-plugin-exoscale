package exoscale

var (
	testConfigAPIKey               = "EXOabcdef0123456789abcdef01"
	testConfigAPISecret            = "ABCDEFGHIJKLMNOPRQSTUVWXYZ0123456789abcdefg"
	testConfigInstanceTemplate     = "Linux Ubuntu 20.04 LTS 64-bit"
	testConfigSnapshotDownload     = true
	testConfigSnapshotDownloadPath = "./output.test"
	testConfigTemplateZones        = []string{"ch-gva-2", "ch-dk-2"}
	testConfigTemplateName         = "test-packer"
	testConfigSSHUsername          = "ubuntu"
	// Deprecated
	testConfigTemplateZone = "ch-dk-2"
)

func (ts *testSuite) TestNewConfig() {
	_, _, err := NewConfig()
	ts.Require().Error(err, "incomplete configuration should return an error")

	config, _, err := NewConfig([]interface{}{map[string]interface{}{
		// Minimal configuration
		"api_key":                testConfigAPIKey,
		"api_secret":             testConfigAPISecret,
		"instance_template":      testConfigInstanceTemplate,
		"snapshot_download":      testConfigSnapshotDownload,
		"snapshot_download_path": testConfigSnapshotDownloadPath,
		"template_name":          testConfigTemplateName,
		"template_zones":         testConfigTemplateZones,
		"ssh_username":           testConfigSSHUsername,
	}}...)
	ts.Require().NoError(err)
	ts.Require().NotNil(config)
	ts.Require().Equal(defaultAPIEnvironment, config.APIEnvironment)
	ts.Require().Equal(defaultInstanceType, config.InstanceType)
	ts.Require().Equal(defaultInstanceDiskSize, config.InstanceDiskSize)
	ts.Require().Equal([]string{defaultInstanceSecurityGroup}, config.InstanceSecurityGroups)
	ts.Require().Equal(defaultInstanceTemplateVisibility, config.InstanceTemplateVisibility)
	ts.Require().Equal(testConfigSnapshotDownload, config.SnapshotDownload)
	ts.Require().Equal(testConfigSnapshotDownloadPath, config.SnapshotDownloadPath)
	ts.Require().Equal(testConfigTemplateZones[0], config.InstanceZone)
	ts.Require().Equal(defaultTemplateBootMode, config.TemplateBootMode)
}

func (ts *testSuite) TestNewConfigDeprecated() {
	config, _, err := NewConfig([]interface{}{map[string]interface{}{
		// Minimal configuration
		"api_key":           testConfigAPIKey,
		"api_secret":        testConfigAPISecret,
		"instance_template": testConfigInstanceTemplate,
		"template_name":     testConfigTemplateName,
		"ssh_username":      testConfigSSHUsername,
		// Deprecated
		"template_zone": testConfigTemplateZone,
	}}...)
	ts.Require().NoError(err)
	ts.Require().NotNil(config)
	ts.Require().Equal(1, len(config.TemplateZones))
	ts.Require().Equal(testConfigTemplateZone, config.TemplateZones[0])
	ts.Require().Equal(testConfigTemplateZone, config.InstanceZone)

	config, warnings, err := NewConfig([]interface{}{map[string]interface{}{
		// Minimal configuration
		"api_key":           testConfigAPIKey,
		"api_secret":        testConfigAPISecret,
		"instance_template": testConfigInstanceTemplate,
		"template_name":     testConfigTemplateName,
		"ssh_username":      testConfigSSHUsername,
		"template_zones":    testConfigTemplateZones,
		// Deprecated
		"template_zone": testConfigTemplateZone,
	}}...)
	ts.Require().NoError(err)
	ts.Require().NotNil(config)
	ts.Require().Equal(testConfigTemplateZones, config.TemplateZones)
	ts.Require().Equal(1, len(warnings))
	ts.Require().Equal("Both template_zones and template_zone are defined; ignoring the latter", warnings[0])
}
