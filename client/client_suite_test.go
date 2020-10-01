package client

import (
	"github.com/jarcoal/httpmock"
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var c *RollbarApiClient

const (
	fixPath = "testdata/fixtures/"
)


func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Suite")
}


var _ = BeforeSuite(func() {
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
	b, err := ioutil.ReadFile(fixPath + path)
	if err != nil {
		panic(err)
	}
	return string(b)
}

