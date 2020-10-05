package client

import (
	"github.com/brianvoe/gofakeit/v5"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

/*
 * Ginkgo setup
 */

var c *RollbarApiClient

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Suite")
}

var _ = BeforeSuite(func() {
	// Seed gofakeit random generator
	gofakeit.Seed(0) // Setting seed to 0 will use time.Now().UnixNano()

	// Setup RollbarApiClient and enable mocking
	var err error
	c, err = NewClient("fakeTokenString")
	Expect(err).NotTo(HaveOccurred())
	httpmock.ActivateNonDefault(c.GetHttpClient())
})

var _ = BeforeEach(func() {
	httpmock.Reset()
})

var _ = AfterSuite(func() {
	httpmock.DeactivateAndReset()
})

// fixture loads a JSON file from the fixtures folder and returns it as a string
func fixture(path string) string {
	const fixPath = "testdata/fixtures/"
	b, err := ioutil.ReadFile(fixPath + path)
	if err != nil {
		panic(err)
	}
	return string(b)
}

/*
 * Testify setup
 */

// ClientTestSuite is a Testify test suite for the Rollbar API client
type ClientTestSuite struct {
	suite.Suite
	client *RollbarApiClient
}

func (s *ClientTestSuite) SetupSuite() {
	// Seed gofakeit random generator
	gofakeit.Seed(0) // Setting seed to 0 will use time.Now().UnixNano()

	// Setup RollbarApiClient and enable mocking
	c, err := NewClient("fakeTokenString")
	s.Nil(err)
	httpmock.ActivateNonDefault(c.GetHttpClient())
	s.client = c
}

func (s *ClientTestSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func (s *ClientTestSuite) BeforeTest() {
	httpmock.Reset()
}

// TestRollbarClientTestSuite connects the Testify test suite to the 'go test'
// built-in testing framework.
func TestRollbarClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
