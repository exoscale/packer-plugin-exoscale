package exoscale

import (
	egoscale "github.com/exoscale/egoscale/v2"
	"github.com/stretchr/testify/mock"
)

func (ts *testSuite) TestArtifact_Destroy() {
	var (
		testConfig = Config{
			TemplateZone: testZone,
		}
		templateDeleted bool
	)

	testTemplate := egoscale.Template{
		ID:   &testTemplateID,
		Zone: &testZone,
	}

	testArtifact := &Artifact{
		builder: &Builder{
			buildID: ts.randomID(),
			config:  &testConfig,
			exo:     ts.exo,
		},
		state:    ts.state,
		template: &testTemplate,
	}

	ts.exo.(*exoscaleClientMock).
		On(
			"DeleteTemplate",
			mock.Anything, // ctx
			testZone,      // zone
			mock.Anything, // template
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(&testTemplate, args.Get(2))
			templateDeleted = true
		}).
		Return(nil)

	ts.Require().NoError(testArtifact.Destroy())
	ts.Require().True(templateDeleted)
}
