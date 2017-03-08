package resource_test

import (
	"github.com/mevansam/cf-cli-api/cfapi"

	"github.com/mevansam/cf-event-resource-type/resource"
	. "github.com/mevansam/cf-event-resource-type/resource/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cloud Foundry Client", func() {

	var (
		cfSessionProvider cfapi.CfSessionProvider
	)

	BeforeEach(func() {
	})

	AfterEach(func() {
	})

	Context("verify client initialization", func() {

		It("validates source arguments are passed in correctly", func() {

			src := resource.Source{
				API:               "https://api.local.pcfdev.io",
				User:              "admin",
				Password:          "admin",
				Org:               "pcfdev-org",
				Space:             "pcfdev-space",
				SkipSSLValidation: true,
				Debug:             true,
				Trace:             true,
			}

			cfSessionProvider = &FakeCfSessionProvider{
				NewCfSessionStub: func(
					apiEndPoint string, userName string, password string,
					orgName string, spaceName string,
					sslDisabled bool, logger *cfapi.Logger) (cfSession cfapi.CfSession, err error) {

					Expect(apiEndPoint).To(Equal(src.API))
					Expect(userName).To(Equal(src.User))
					Expect(password).To(Equal(src.Password))
					Expect(orgName).To(Equal(src.Org))
					Expect(spaceName).To(Equal(src.Space))
					Expect(sslDisabled).To(Equal(src.SkipSSLValidation))

					return &FakeCfSession{}, nil
				},
			}

			_, err := resource.NewCfClient(cfSessionProvider, src)
			Expect(err).Should(BeNil())
		})
	})
})
