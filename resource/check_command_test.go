package resource_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"code.cloudfoundry.org/cli/cf/api"
	"code.cloudfoundry.org/cli/cf/models"

	"github.com/kr/pretty"
	"github.com/mevansam/cf-cli-api/cfapi"
	"github.com/mitchellh/colorstring"

	"github.com/mevansam/cf-event-resource-type/resource"
	. "github.com/mevansam/cf-event-resource-type/resource/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Check Command", func() {

	var (
		logger *cfapi.Logger

		cfSessionProvider *FakeCfSessionProvider
		cfSession         *FakeCfSession

		appSummaryRepo *FakeAppSummaryRepository

		src resource.Source
		to  time.Time
	)

	BeforeEach(func() {

		logger = cfapi.NewLogger(testing.Verbose(), strconv.FormatBool(testing.Verbose()))

		appSummaryRepo = &FakeAppSummaryRepository{
			GetSummariesInCurrentSpaceStub: func() (apps []models.Application, apiErr error) {
				return srcApps, nil
			},
		}

		cfSession = &FakeCfSession{
			GetSessionLoggerStub: func() *cfapi.Logger {
				return logger
			},
			AppSummaryStub: func() api.AppSummaryRepository {
				return appSummaryRepo
			},
			GetAllEventsInSpaceStub: func(from time.Time, inclusive bool) (events map[string]cfapi.CfEvent, err error) {
				events = make(map[string]cfapi.CfEvent)
				for guid, cfEvent := range testEvents {
					eventList := []models.EventFields{}
					for _, event := range cfEvent.EventList {
						if (inclusive && event.Timestamp.Equal(from)) ||
							event.Timestamp.After(from) && event.Timestamp.Before(to) {
							eventList = append(eventList, event)
						}
					}
					cfEvent.EventList = eventList
					events[guid] = cfEvent
				}
				return
			},
			GetAllEventsForAppStub: func(appGUID string, from time.Time, inclusive bool) (cfEvent cfapi.CfEvent, err error) {
				cfEvent, ok := testEvents[appGUID]
				if ok {
					eventList := []models.EventFields{}
					for _, event := range cfEvent.EventList {
						if (inclusive && event.Timestamp.Equal(from)) ||
							event.Timestamp.After(from) && event.Timestamp.Before(to) {
							eventList = append(eventList, event)
						}
					}
					cfEvent.EventList = eventList
				}
				return
			},
		}

		cfSessionProvider = &FakeCfSessionProvider{
			NewCfSessionStub: func(
				apiEndPoint string, userName string, password string,
				orgName string, spaceName string,
				sslDisabled bool, logger *cfapi.Logger) (cfapi.CfSession, error) {

				Expect(apiEndPoint).To(Equal(src.API))
				Expect(userName).To(Equal(src.User))
				Expect(password).To(Equal(src.Password))
				Expect(orgName).To(Equal(src.Org))
				Expect(spaceName).To(Equal(src.Space))
				Expect(sslDisabled).To(Equal(src.SkipSSLValidation))

				return cfSession, nil
			},
		}

		for _, cfEvent := range testEvents {
			for _, event := range cfEvent.EventList {
				fmt.Printf("Test Data: Event: %s / %s, Source %s / %s: \n",
					event.Timestamp.Format(time.RFC3339), event.Name,
					cfEvent.Name, cfEvent.Type)
			}
		}
	})

	AfterEach(func() {
	})

	Context("verify running of check command", func() {

		It("validates retrieval of versions", func() {

			var (
				err      error
				versions []resource.Version
			)

			request := resource.NewCheckRequest()
			request.Source = resource.Source{
				API:               "https://api.local.pcfdev.io",
				User:              "admin",
				Password:          "admin",
				Org:               "pcfdev-org",
				Space:             "pcfdev-space",
				Apps:              []string{"app1", "app2"},
				SkipSSLValidation: true,
				Debug:             true,
				Trace:             true,
			}

			cfClient, err := resource.NewCfClient(cfSessionProvider, src)
			Expect(err).Should(BeNil())
			command := resource.NewCheckCommand(cfClient)

			to, _ = time.Parse(time.RFC3339, "2017-03-03T18:12:00+01:00")
			fmt.Printf(colorstring.Color("\n[green]===> Detecting event versions in test data to: %s\n"), to)

			versions, err = command.Run(request)
			Expect(err).Should(BeNil())
			Expect(len(versions)).To(Equal(1))
			Expect(versions[0]["app1"]).To(Equal("19b9d70b-6ebe-47d7-9313-f0c213445036|app1|app|created|1970-01-01T00:00:00Z"))
			Expect(versions[0]["app2"]).To(Equal("62aa7e42-f047-46de-98f1-ff4818a82d6e|app2|app|created|1970-01-01T00:00:00Z"))

			to, _ = time.Parse(time.RFC3339, "2017-03-03T18:13:00+01:00")
			fmt.Printf(colorstring.Color("\n[green]===> Detecting event versions in test data to: %s\n"), to)

			request.Version = versions[0]
			versions, err = command.Run(request)
			Expect(err).Should(BeNil())
			Expect(len(versions)).To(Equal(1))
			Expect(assertVersionEqual(versions[0]["app1"], "19b9d70b-6ebe-47d7-9313-f0c213445036|app1|app|created|2017-03-03T18:12:48+01:00")).To(BeTrue())
			Expect(versions[0]["app2"]).To(Equal("62aa7e42-f047-46de-98f1-ff4818a82d6e|app2|app|created|1970-01-01T00:00:00Z"))

			to, _ = time.Parse(time.RFC3339, "2017-03-03T18:14:00+01:00")
			fmt.Printf(colorstring.Color("\n[green]===> Detecting event versions in test data  to: %s\n"), to)

			request.Version = versions[0]
			versions, err = command.Run(request)
			Expect(err).Should(BeNil())
			Expect(len(versions)).To(Equal(1))
			Expect(assertVersionEqual(versions[0]["app1"], "19b9d70b-6ebe-47d7-9313-f0c213445036|app1|app|created|2017-03-03T18:12:48+01:00")).To(BeTrue())
			Expect(assertVersionEqual(versions[0]["app2"], "62aa7e42-f047-46de-98f1-ff4818a82d6e|app2|app|created|2017-03-03T18:13:19+01:00")).To(BeTrue())

			to, _ = time.Parse(time.RFC3339, "2017-03-03T18:15:00+01:00")
			fmt.Printf(colorstring.Color("\n[green]===> Detecting event versions in test data  to: %s\n"), to)

			request.Version = versions[0]
			versions, err = command.Run(request)
			Expect(err).Should(BeNil())
			Expect(len(versions)).To(Equal(1))
			Expect(assertVersionEqual(versions[0]["app1"], "19b9d70b-6ebe-47d7-9313-f0c213445036|app1|app|modified|2017-03-03T18:14:44+01:00")).To(BeTrue())
			Expect(assertVersionEqual(versions[0]["app2"], "62aa7e42-f047-46de-98f1-ff4818a82d6e|app2|app|routed-added|2017-03-03T18:14:59+01:00")).To(BeTrue())

			to, _ = time.Parse(time.RFC3339, "2017-03-03T18:16:00+01:00")
			fmt.Printf(colorstring.Color("\n[green]===> Detecting event versions in test data  to: %s\n"), to)

			request.Version = versions[0]
			versions, err = command.Run(request)
			Expect(err).Should(BeNil())
			Expect(len(versions)).To(Equal(1))
			Expect(assertVersionEqual(versions[0]["app1"], "19b9d70b-6ebe-47d7-9313-f0c213445036|app1|app|modified|2017-03-03T18:14:44+01:00")).To(BeTrue())
			Expect(assertVersionEqual(versions[0]["app2"], "62aa7e42-f047-46de-98f1-ff4818a82d6e|app2|app|modified|2017-03-03T18:15:39+01:00")).To(BeTrue())

			to, _ = time.Parse(time.RFC3339, "2017-03-03T18:17:00+01:00")
			fmt.Printf(colorstring.Color("\n[green]===> Detecting event versions in test data  to: %s\n"), to)

			request.Version = versions[0]
			versions, err = command.Run(request)
			Expect(err).Should(BeNil())
			Expect(len(versions)).To(Equal(1))
			Expect(assertVersionEqual(versions[0]["app1"], "19b9d70b-6ebe-47d7-9313-f0c213445036|app1|app|scaled|2017-03-03T18:16:46+01:00")).To(BeTrue())
			Expect(assertVersionEqual(versions[0]["app2"], "62aa7e42-f047-46de-98f1-ff4818a82d6e|app2|app|modified|2017-03-03T18:15:39+01:00")).To(BeTrue())

			fmt.Printf("Versions: %# v", pretty.Formatter(versions))

			to, _ = time.Parse(time.RFC3339, "2017-03-03T18:18:00+01:00")
			fmt.Printf(colorstring.Color("\n[green]===> Detecting event versions in test data  to: %s\n"), to)

			request.Version = versions[0]
			versions, err = command.Run(request)
			Expect(err).Should(BeNil())
			Expect(len(versions)).To(Equal(1))
			Expect(assertVersionEqual(versions[0]["app1"], "19b9d70b-6ebe-47d7-9313-f0c213445036|app1|app|routed-added|2017-03-03T18:17:02+01:00")).To(BeTrue())
			Expect(assertVersionEqual(versions[0]["app2"], "62aa7e42-f047-46de-98f1-ff4818a82d6e|app2|app|modified|2017-03-03T18:15:39+01:00")).To(BeTrue())

			to, _ = time.Parse(time.RFC3339, "2017-03-03T18:19:00+01:00")
			fmt.Printf(colorstring.Color("\n[green]===> Detecting event versions in test data  to: %s\n"), to)

			request.Version = versions[0]
			versions, err = command.Run(request)
			Expect(err).Should(BeNil())
			Expect(len(versions)).To(Equal(1))
			Expect(assertVersionEqual(versions[0]["app1"], "19b9d70b-6ebe-47d7-9313-f0c213445036|app1|app|modified|2017-03-03T18:17:52+01:00")).To(BeTrue())
			Expect(assertVersionEqual(versions[0]["app2"], "62aa7e42-f047-46de-98f1-ff4818a82d6e|app2|app|modified|2017-03-03T18:15:39+01:00")).To(BeTrue())

			fmt.Printf("Versions: %# v", pretty.Formatter(versions))

			to, _ = time.Parse(time.RFC3339, "2017-03-03T18:20:00+01:00")
			fmt.Printf(colorstring.Color("\n[green]===> Detecting event versions in test data  to: %s\n"), to)

			request.Version = versions[0]
			versions, err = command.Run(request)
			Expect(err).Should(BeNil())
			Expect(len(versions)).To(Equal(1))
			Expect(assertVersionEqual(versions[0]["app1"], "19b9d70b-6ebe-47d7-9313-f0c213445036|app1|app|modified|2017-03-03T18:17:52+01:00")).To(BeTrue())
			Expect(assertVersionEqual(versions[0]["app2"], "62aa7e42-f047-46de-98f1-ff4818a82d6e|app2|app|modified|2017-03-03T18:15:39+01:00")).To(BeTrue())
		})
	})
})
