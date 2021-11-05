package exoscale

var (
	testConfigAPIKey           = "EXOabcdef0123456789abcdef01"
	testConfigAPISecret        = "ABCDEFGHIJKLMNOPRQSTUVWXYZ0123456789abcdefg"
	testConfigInstanceTemplate = "Linux Ubuntu 20.04 LTS 64-bit"
	testConfigTemplateZone     = "ch-gva-2"
	testConfigTemplateName     = "test-packer"
	testConfigSSHUsername      = "ubuntu"
)

func (ts *testSuite) TestNewConfig() {
	_, err := NewConfig()
	ts.Require().Error(err, "incomplete configuration should return an error")

	// Minimal configuration
	config, err := NewConfig([]interface{}{map[string]interface{}{
		"api_key":           testConfigAPIKey,
		"api_secret":        testConfigAPISecret,
		"instance_template": testConfigInstanceTemplate,
		"template_name":     testConfigTemplateName,
		"template_zone":     testConfigTemplateZone,
		"ssh_username":      testConfigSSHUsername,
	}}...)
	ts.Require().NoError(err)
	ts.Require().NotNil(config)
	ts.Require().Equal(defaultAPIEnvironment, config.APIEnvironment)
	ts.Require().Equal(defaultInstanceType, config.InstanceType)
	ts.Require().Equal(defaultInstanceDiskSize, config.InstanceDiskSize)
	ts.Require().Equal([]string{defaultInstanceSecurityGroup}, config.InstanceSecurityGroups)
	ts.Require().Equal(defaultInstanceTemplateVisibility, config.InstanceTemplateVisibility)
	ts.Require().Equal(config.InstanceZone, testConfigTemplateZone)
	ts.Require().Equal(defaultTemplateBootMode, config.TemplateBootMode)
}
