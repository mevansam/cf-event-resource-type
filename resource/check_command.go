package resource

import (
	"fmt"
	"time"

	"code.cloudfoundry.org/cli/cf/models"

	"github.com/mevansam/cf-cli-api/filters"
	"github.com/mevansam/cf-cli-api/utils"
)

// CheckCommand -
type CheckCommand struct {
	client *CfClient
}

// epoch -
var epoch, _ = time.Parse(time.RFC3339, "1970-01-01T00:00:00Z")

// NewCheckCommand -
func NewCheckCommand(client *CfClient) *CheckCommand {
	return &CheckCommand{client}
}

// Run -
func (c *CheckCommand) Run(request CheckRequest) (versions []Version, err error) {

	var (
		getAppEvents, allApps, exists bool

		appEventVersion string
		from            time.Time

		appsInSpace []models.Application

		appEvents    []filters.AppEvent
		lastAppEvent filters.AppEvent
		eventFilter  filters.EventFilter
	)

	logger := c.client.session.GetSessionLogger()

	allApps = len(request.Source.Apps) == 0
	eventFilter = filters.NewAppEventFilter(c.client.session)
	newVersion := make(Version)

	if appsInSpace, err = c.client.session.AppSummary().GetSummariesInCurrentSpace(); err != nil {
		return
	}

	for appName, appEventVersion := range request.Version {

		if lastAppEvent, err = filters.NewAppEvent(appEventVersion); err != nil {
			return
		}
		if lastAppEvent.EventType != filters.EtDeleted {
			if _, exists = utils.ContainsApp(appName, appsInSpace); !exists {
				appsInSpace = append(appsInSpace, models.Application{
					ApplicationFields: models.ApplicationFields{
						GUID: lastAppEvent.SourceGUID,
						Name: lastAppEvent.SourceName,
					},
				})
			}
		}
	}

	if len(appsInSpace) > 0 {
		for _, app := range appsInSpace {

			getAppEvents = true
			if !allApps {
				_, getAppEvents = utils.ContainsInStrings([]string{app.Name}, request.Source.Apps)
			}

			if getAppEvents {

				appEventVersion, exists = request.Version[app.Name]
				if exists {
					if lastAppEvent, err = filters.NewAppEvent(appEventVersion); err != nil {
						return
					}
					from = lastAppEvent.Timestamp
				} else {
					from = epoch
				}

				logger.DebugMessage("Retrieving new events for app '%s' after timestamp '%s'.",
					app.Name, from.Format(time.RFC3339))

				if appEvents, err = eventFilter.GetEventsForApp(app.GUID, from, true); err != nil {
					return
				}
				if exists {
					if len(appEvents) > 0 {
						for _, ae := range appEvents {
							if ae.EventType != lastAppEvent.EventType ||
								ae.Timestamp.After(lastAppEvent.Timestamp) {
								newVersion[app.Name] = fmt.Sprintf("%s", ae)
								break
							}
						}
					} else {
						newVersion[app.Name] = fmt.Sprintf("%s", lastAppEvent)
					}
				} else {
					// Iterate back in time until last deploy event is found and start from there
					i := len(appEvents) - 1
					for i >= 0 {
						ae := appEvents[i]
						if ae.EventType == filters.EtCreated || ae.EventType == filters.EtModified {
							break
						}
						i--
					}
					if i >= 0 {
						newVersion[app.Name] = fmt.Sprintf("%s", appEvents[i])
					}
				}
			}
		}
		versions = []Version{newVersion}
	} else {
		versions = []Version{request.Version}
	}
	return
}
