package client

import (
	"github.com/jarcoal/httpmock"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var c *RollbarApiClient


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
