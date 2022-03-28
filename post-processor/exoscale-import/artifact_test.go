package exoscaleimport

import (
	egoscale "github.com/exoscale/egoscale/v2"
	"github.com/stretchr/testify/mock"
)

func (ts *testSuite) TestArtifact_Destroy() {
	var (
		testConfig = Config{
			TemplateZones: testTemplateZones,
		}
		templateDeleted = 0
	)

	testTemplates := []*egoscale.Template{}
	testTemplatesMap := map[string]*egoscale.Template{}
	for i := 0; i < len(testTemplateZones); i++ {
		testTemplates = append(
			testTemplates,
			&egoscale.Template{
				ID:   &testTemplateID,
				Zone: &testTemplateZones[i],
			},
		)
		testTemplatesMap[testTemplateZones[i]] = testTemplates[i]
	}

	testArtifact := &Artifact{
		postProcessor: &PostProcessor{
			config: &testConfig,
			exo:    ts.exo,
		},
		state:     ts.state,
		templates: testTemplates,
	}

	for i := 0; i < len(testTemplateZones); i++ {
		ts.exo.(*exoscaleClientMock).
			On(
				"DeleteTemplate",
				mock.Anything,        // ctx
				testTemplateZones[i], // zone
				mock.Anything,        // template
			).
			Run(func(args mock.Arguments) {
				ts.Require().Equal(testTemplatesMap[args.Get(1).(string)], args.Get(2))
				templateDeleted++
			}).
			Return(nil)
	}

	ts.Require().NoError(testArtifact.Destroy())
	ts.Require().Equal(templateDeleted, 1) // NB: DeleteTemplate needs be called only once

}
