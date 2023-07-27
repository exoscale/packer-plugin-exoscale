package exoscale

import (
	"encoding/base64"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/stretchr/testify/suite"

	egoscale "github.com/exoscale/egoscale/v2"
)

var (
	testInstanceDiskSize           = defaultInstanceDiskSize
	testInstanceID                 = new(testSuite).randomID()
	testInstanceIPAddress          = net.ParseIP("1.2.3.4")
	testInstanceName               = new(testSuite).randomString(10)
	testInstancePrivateNetworkID   = new(testSuite).randomID()
	testInstancePrivateNetworkName = new(testSuite).randomString(10)
	testInstanceSecurityGroupID    = new(testSuite).randomID()
	testInstanceSecurityGroupName  = new(testSuite).randomString(10)
	testInstanceSnapshotID         = new(testSuite).randomID()
	testInstanceTypeID             = new(testSuite).randomID()
	testInstanceTypeName           = defaultInstanceType
	testInstanceZone               = "ch-gva-2"
	testSnapshotDownload           = true
	testSnapshotDownloadPath       = filepath.Join(os.TempDir(), "packer-plugin-test-"+new(testSuite).randomString(6))
	testTemplateID                 = new(testSuite).randomID()
	testTemplateName               = "packer-plugin-test-" + new(testSuite).randomString(6)
	testTemplateZones              = []string{"ch-gva-2", "ch-dk-2"}
	testUserData                   = "echo test > /etc/test.txt"
	testUserDataBase64             = base64.StdEncoding.EncodeToString([]byte(testUserData))
	testUserDataFile               = "userdata.txt"

	testSeededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type testSuite struct {
	suite.Suite

	exo   exoscaleClient
	state *multistep.BasicStateBag
	ui    packer.Ui
}

func (ts *testSuite) SetupTest() {
	ts.exo = new(exoscaleClientMock)
	ts.ui = new(packer.MockUi)

	ts.state = new(multistep.BasicStateBag)
	ts.state.Put("ui", ts.ui)
	ts.state.Put("templates", []*egoscale.Template{})
}

func (ts *testSuite) TearDownTest() {
	ts.exo = nil
	ts.state = nil
	ts.ui = nil
}

func (ts *testSuite) randomID() string {
	id, err := uuid.NewV4()
	if err != nil {
		ts.T().Fatalf("unable to generate a new UUID: %s", err)
	}
	return id.String()
}

func (ts *testSuite) randomStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[testSeededRand.Intn(len(charset))]
	}
	return string(b)
}

func (ts *testSuite) randomString(length int) string {
	const defaultCharset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	return ts.randomStringWithCharset(length, defaultCharset)
}

func (ts *testSuite) TestBuilder_Prepare() {
	type args struct {
		raws []interface{}
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "fail missing required config args",
			args:    args{},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				raws: []interface{}{map[string]interface{}{
					"api_key":           testConfigAPIKey,
					"api_secret":        testConfigAPISecret,
					"instance_template": testConfigInstanceTemplate,
					"template_name":     testConfigTemplateName,
					"template_zones":    testConfigTemplateZones,
					"ssh_username":      testConfigSSHUsername,
				}},
			},
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			b := new(Builder)

			_, _, err := b.Prepare(tt.args.raws...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Prepare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			ts.Require().NotNil(b.config)
		})
	}
}

func TestSuiteExoscaleBuilder(t *testing.T) {
	suite.Run(t, new(testSuite))
}
