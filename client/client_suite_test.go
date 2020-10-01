package client_test

import (
	"github.com/jarcoal/httpmock"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var c *client.RollbarApiClient

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Suite")
}


var _ = BeforeSuite(func() {
	var err error
	c, err = client.NewClient("fakeTokenString")
	Expect(err).NotTo(HaveOccurred())
	httpmock.ActivateNonDefault(c.GetHttpClient())
})

var _ = BeforeEach(func() {
	// remove any mocks
	httpmock.Reset()
})

var _ = AfterSuite(func() {
	httpmock.DeactivateAndReset()
})
