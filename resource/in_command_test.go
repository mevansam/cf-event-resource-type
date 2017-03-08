package resource_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mevansam/cf-cli-api/cfapi"
	"github.com/mevansam/cf-cli-api/filters"
	"github.com/mevansam/cf-event-resource-type/resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("In Command", func() {

	var (
		logger *cfapi.Logger

		err    error
		tmpDir string
	)

	BeforeEach(func() {

		logger = cfapi.NewLogger(testing.Verbose(), strconv.FormatBool(testing.Verbose()))

		tmpDir, err = ioutil.TempDir("", "")
		logger.DebugMessage("Creating download content in temp dir: %s", tmpDir)

		os.MkdirAll(tmpDir, os.ModePerm)
		if err != nil {
			panic(err.Error())
		}
	})

	AfterEach(func() {
		os.RemoveAll(tmpDir)
	})

	Context("verify running of in command", func() {

		It("validates the download directory content is created correctly", func() {
			request := resource.NewInRequest()
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
			version := resource.Version{
				"app1": "19b9d70b-6ebe-47d7-9313-f0c213445036|app1|app|created|2017-03-03T18:12:48+01:00",
				"app2": "62aa7e42-f047-46de-98f1-ff4818a82d6e|app2|app|created|2017-03-03T18:13:19+01:00",
			}
			request.Version = &version

			command := resource.NewInCommand(tmpDir)
			response, err := command.Run(request)
			Expect(err).Should(BeNil())

			fileList := []string{tmpDir + "/ENV.sh", tmpDir + "/version.json", tmpDir + "/metadata.json"}
			for _, app := range request.Source.Apps {
				fileList = append(fileList, tmpDir+"/"+app+".event")
				fileList = append(fileList, tmpDir+"/"+app+".timestamp")
			}

			for _, f := range fileList {
				_, err := os.Stat(f)
				Expect(os.IsNotExist(err)).ShouldNot(BeTrue())
			}

			envContent, err := ioutil.ReadFile(tmpDir + "/ENV.sh")
			Expect(err).Should(BeNil())
			Expect(regexp.MatchString("export CF_API=\""+request.Source.API+"\"\n", string(envContent))).Should(BeTrue())
			Expect(regexp.MatchString("export CF_USER=\""+request.Source.User+"\"\n", string(envContent))).Should(BeTrue())
			Expect(regexp.MatchString("export CF_PASSWORD=\""+request.Source.Password+"\"\n", string(envContent))).Should(BeTrue())
			Expect(regexp.MatchString("export CF_ORG=\""+request.Source.Org+"\"\n", string(envContent))).Should(BeTrue())
			Expect(regexp.MatchString("export CF_SPACE=\""+request.Source.Space+"\"\n", string(envContent))).Should(BeTrue())

			re := regexp.MustCompile("export CF_APPS=\"(.*) \"")
			matches := re.FindAllStringSubmatch("export CF_APPS=\"app1 app2 \"", -1)
			Expect(matches).ShouldNot(BeNil())
			Expect(len(matches)).To(Equal(1))
			Expect(len(matches[0])).To(Equal(2))
			apps := strings.Split(matches[0][1], " ")
			Expect(apps).To(Equal(request.Source.Apps))

			for _, app := range apps {
				appEvent, err := filters.NewAppEvent(response.Version[app])
				Expect(err).Should(BeNil())
				event, err := ioutil.ReadFile(tmpDir + "/" + app + ".event")
				Expect(err).Should(BeNil())
				Expect(string(event)).To(Equal(string(appEvent.EventType)))
				timestamp, err := ioutil.ReadFile(tmpDir + "/" + app + ".timestamp")
				Expect(err).Should(BeNil())
				ts, err := strconv.ParseInt(string(timestamp), 10, 64)
				Expect(err).Should(BeNil())
				Expect(time.Unix(ts, 0)).To(Equal(appEvent.Timestamp))
			}

			fmt.Printf("Response: %# v\n", response)
		})
	})
})
