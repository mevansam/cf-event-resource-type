package resource

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/mevansam/cf-cli-api/filters"
)

// InCommand -
type InCommand struct {
	downloadDir string
}

// NewInCommand -
func NewInCommand(downloadDir string) *InCommand {
	return &InCommand{downloadDir}
}

// Run -
func (c *InCommand) Run(request InRequest) (response InResponse, err error) {

	/*
		The "in" command will create the following content within the download directory

		/downloadDir/
		   ENV.sh
		   [APP_NAME].event
		   versions
		   metadata
	*/

	var (
		envFile, appGUIDFile, appEventTypeFile, appEventTimestampFile, versionFile, metaDataFile *os.File

		appNames bytes.Buffer
		appEvent filters.AppEvent
	)

	if request.Version == nil {
		Fatal("no version information passed to the 'in' command processor")
	}

	if envFile, err = os.OpenFile(c.downloadDir+"/ENV.sh",
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0744); err != nil {
		return
	}
	defer envFile.Close()
	envFile.WriteString(fmt.Sprintf("export CF_API=\"%s\"\n", request.Source.API))
	envFile.WriteString(fmt.Sprintf("export CF_USER=\"%s\"\n", request.Source.User))
	if len(request.Source.Password) > 0 {
		envFile.WriteString(fmt.Sprintf("export CF_PASSWORD=\"%s\"\n", request.Source.Password))
	}
	if len(request.Source.SSOToken) > 0 {
		envFile.WriteString(fmt.Sprintf("export CF_SSO_TOKEN=\"%s\"\n", request.Source.SSOToken))
	}
	envFile.WriteString(fmt.Sprintf("export CF_ORG=\"%s\"\n", request.Source.Org))
	envFile.WriteString(fmt.Sprintf("export CF_SPACE=\"%s\"\n", request.Source.Space))
	if request.Source.SkipSSLValidation {
		envFile.WriteString("export CF_SKIP_SSL_VALIDATION=\"--skip-ssl-validation\"\n")
	}
	appNames.WriteString("export CF_APPS=\"")

	for appName, appEventData := range *request.Version {

		appNames.WriteString(fmt.Sprintf("%s ", appName))

		if appEvent, err = filters.NewAppEvent(appEventData); err != nil {
			return
		}

		if appGUIDFile, err = os.OpenFile(c.downloadDir+fmt.Sprintf("/%s.guid", appName),
			os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
			return
		}
		appGUIDFile.WriteString(fmt.Sprintf("%s", appEvent.SourceGUID))
		appGUIDFile.Close()

		if appEventTypeFile, err = os.OpenFile(c.downloadDir+fmt.Sprintf("/%s.event", appName),
			os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
			return
		}
		appEventTypeFile.WriteString(fmt.Sprintf("%s", appEvent.EventType))
		appEventTypeFile.Close()

		if appEventTimestampFile, err = os.OpenFile(c.downloadDir+fmt.Sprintf("/%s.timestamp", appName),
			os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
			return
		}
		appEventTimestampFile.WriteString(fmt.Sprintf("%d", appEvent.Timestamp.Unix()))
		appEventTimestampFile.Close()

		response.Metadata = append(response.Metadata, MetadataPair{
			Name:  appName,
			Value: fmt.Sprintf("%s @ %s", appEvent.EventType, appEvent.Timestamp.Format(time.RFC3339)),
		})
	}
	appNames.WriteString("\"")
	envFile.WriteString(appNames.String())

	response.Version = *request.Version

	if versionFile, err = os.OpenFile(c.downloadDir+"/version.json",
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		return
	}
	if err = json.NewEncoder(versionFile).Encode(response.Version); err != nil {
		return
	}
	defer versionFile.Close()

	if metaDataFile, err = os.OpenFile(c.downloadDir+"/metadata.json",
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		return
	}
	if err = json.NewEncoder(metaDataFile).Encode(response.Metadata); err != nil {
		return
	}
	defer metaDataFile.Close()

	return
}
