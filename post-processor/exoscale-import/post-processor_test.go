package exoscaleimport

import (
	"math/rand"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/stretchr/testify/suite"
)

var (
	testTemplateID = new(testSuite).randomID()
	testZone       = "ch-gva-2"

	testSeededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type testSuite struct {
	suite.Suite

	exo   exoscaleClient
	sos   s3Client
	state *multistep.BasicStateBag
	ui    packer.Ui
}

func (ts *testSuite) SetupTest() {
	ts.exo = new(exoscaleClientMock)
	ts.sos = new(s3ClientMock)
	ts.ui = new(packer.MockUi)

	ts.state = new(multistep.BasicStateBag)
	ts.state.Put("ui", ts.ui)
}

func (ts *testSuite) TearDownTest() {
	ts.exo = nil
	ts.sos = nil
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

func TestSuiteExoscaleImportPostProcessor(t *testing.T) {
	suite.Run(t, new(testSuite))
}
