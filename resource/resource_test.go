package resource_test

import (
	"strings"
	"time"

	"code.cloudfoundry.org/cli/cf/models"
	"github.com/mevansam/cf-cli-api/cfapi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Resource Commands Test Suite")
}

func assertVersionEqual(v1, v2 string) bool {
	fields := strings.Split(v2, "|")
	ts, _ := time.Parse(time.RFC3339, fields[4])
	return v1 == strings.Join(append(fields[:4], time.Unix(ts.Unix(), 0).Format(time.RFC3339)), "|")
}

var srcApps = []models.Application{
	models.Application{
		ApplicationFields: models.ApplicationFields{
			Name: "app1",
			GUID: "19b9d70b-6ebe-47d7-9313-f0c213445036",
		},
	},
	models.Application{
		ApplicationFields: models.ApplicationFields{
			Name: "app2",
			GUID: "62aa7e42-f047-46de-98f1-ff4818a82d6e",
		},
	},
}

var testEvents = map[string]cfapi.CfEvent{
	"d558dc89-92c8-409b-8d37-62c31671933b": {
		GUID: "d558dc89-92c8-409b-8d37-62c31671933b",
		Name: "mksapp2-1",
		Type: "route",
		EventList: []models.EventFields{
			{
				GUID:        "1615ed27-2773-4a8c-80f0-aa9a5c074ef2",
				Name:        "audit.route.create",
				Timestamp:   time.Unix(1488561298, 0),
				Description: "",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
		},
	},
	"88f1ae1c-130b-4afb-a9e4-5e8fa5fc03ed": {
		GUID: "88f1ae1c-130b-4afb-a9e4-5e8fa5fc03ed",
		Name: "mksapp1-1",
		Type: "route",
		EventList: []models.EventFields{
			{
				GUID:        "5a633540-6c19-478c-b700-fff68787ea52",
				Name:        "audit.route.create",
				Timestamp:   time.Unix(1488561421, 0),
				Description: "",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
		},
	},
	"19b9d70b-6ebe-47d7-9313-f0c213445036": {
		GUID: "19b9d70b-6ebe-47d7-9313-f0c213445036",
		Name: "app1",
		Type: "app",
		EventList: []models.EventFields{
			{
				GUID:        "eeaff937-5776-4ba3-9555-a2cba60daacc",
				Name:        "audit.app.create",
				Timestamp:   time.Unix(1488561136, 0),
				Description: "instances: 1, memory: 512, state: STOPPED, environment_json: PRIVATE DATA HIDDEN",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "4ca9d804-b61f-406f-81fd-d1e0577bf7f6",
				Name:        "audit.app.map-route",
				Timestamp:   time.Unix(1488561141, 0),
				Description: "",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "544290c5-4485-4ef2-93cf-42599c1df72c",
				Name:        "audit.app.update",
				Timestamp:   time.Unix(1488561141, 0),
				Description: "",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "4adf8988-ed73-42dd-942f-c7966b5cc36d",
				Name:        "audit.app.update",
				Timestamp:   time.Unix(1488561168, 0),
				Description: "state: STARTED",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "814a279f-da11-43ee-b7e5-653f15a4ce25",
				Name:        "audit.app.update",
				Timestamp:   time.Unix(1488561260, 0),
				Description: "instances: 1, memory: 512, environment_json: PRIVATE DATA HIDDEN",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "825d3adc-007b-4c04-8c93-c5196f367fa4",
				Name:        "audit.app.update",
				Timestamp:   time.Unix(1488561280, 0),
				Description: "state: STOPPED",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "1f3cceab-cfab-4d14-8b42-a5ade9158df8",
				Name:        "audit.app.update",
				Timestamp:   time.Unix(1488561284, 0),
				Description: "state: STARTED",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "2caaebff-a694-44a0-8fe8-255ef1d566f8",
				Name:        "audit.app.process.crash",
				Timestamp:   time.Unix(1488561314, 0),
				Description: "index: 0, reason: CRASHED, exit_description: 2 error(s) occurred:\n\n* 2 error(s) occurred:\n\n* Exited with status 148\n* cancelled\n* cancelled",
				Actor:       "19b9d70b-6ebe-47d7-9313-f0c213445036",
				ActorName:   "web",
			},
			{
				GUID:        "24084929-5120-4692-af91-8cb316cdfef6",
				Name:        "app.crash",
				Timestamp:   time.Unix(1488561314, 0),
				Description: "index: 0, reason: CRASHED, exit_description: 2 error(s) occurred:\n\n* 2 error(s) occurred:\n\n* Exited with status 148\n* cancelled\n* cancelled",
				Actor:       "19b9d70b-6ebe-47d7-9313-f0c213445036",
				ActorName:   "app1",
			},
			{
				GUID:        "1b03c2d4-c6d3-4411-9360-f8a5577f831d",
				Name:        "audit.app.process.crash",
				Timestamp:   time.Unix(1488561322, 0),
				Description: "index: 0, reason: CRASHED, exit_description: 2 error(s) occurred:\n\n* 2 error(s) occurred:\n\n* Exited with status 148\n* cancelled\n* cancelled",
				Actor:       "19b9d70b-6ebe-47d7-9313-f0c213445036",
				ActorName:   "web",
			},
			{
				GUID:        "594fb8b0-3ba8-4eeb-aab0-d490f8ce7683",
				Name:        "app.crash",
				Timestamp:   time.Unix(1488561322, 0),
				Description: "index: 0, reason: CRASHED, exit_description: 2 error(s) occurred:\n\n* 2 error(s) occurred:\n\n* Exited with status 148\n* cancelled\n* cancelled",
				Actor:       "19b9d70b-6ebe-47d7-9313-f0c213445036",
				ActorName:   "app1",
			},
			{
				GUID:        "1ab1a14d-fbe9-4d72-98f5-931026401271",
				Name:        "audit.app.process.crash",
				Timestamp:   time.Unix(1488561329, 0),
				Description: "index: 0, reason: CRASHED, exit_description: 2 error(s) occurred:\n\n* 2 error(s) occurred:\n\n* Exited with status 148\n* cancelled\n* cancelled",
				Actor:       "19b9d70b-6ebe-47d7-9313-f0c213445036",
				ActorName:   "web",
			},
			{
				GUID:        "0c79fd84-1b80-402d-9878-39f71e9836d7",
				Name:        "app.crash",
				Timestamp:   time.Unix(1488561329, 0),
				Description: "index: 0, reason: CRASHED, exit_description: 2 error(s) occurred:\n\n* 2 error(s) occurred:\n\n* Exited with status 148\n* cancelled\n* cancelled",
				Actor:       "19b9d70b-6ebe-47d7-9313-f0c213445036",
				ActorName:   "app1",
			},
			{
				GUID:        "9a684e79-acc4-47ca-b943-24ff66697289",
				Name:        "audit.app.process.crash",
				Timestamp:   time.Unix(1488561388, 0),
				Description: "index: 0, reason: CRASHED, exit_description: 2 error(s) occurred:\n\n* 2 error(s) occurred:\n\n* Exited with status 148\n* cancelled\n* cancelled",
				Actor:       "19b9d70b-6ebe-47d7-9313-f0c213445036",
				ActorName:   "web",
			},
			{
				GUID:        "fa2bc385-a523-4f48-a718-68ac9a0a4657",
				Name:        "app.crash",
				Timestamp:   time.Unix(1488561388, 0),
				Description: "index: 0, reason: CRASHED, exit_description: 2 error(s) occurred:\n\n* 2 error(s) occurred:\n\n* Exited with status 148\n* cancelled\n* cancelled",
				Actor:       "19b9d70b-6ebe-47d7-9313-f0c213445036",
				ActorName:   "app1",
			},
			{
				GUID:        "2dc01aa7-2cfa-466c-8307-4b4166ba8b8d",
				Name:        "audit.app.update",
				Timestamp:   time.Unix(1488561406, 0),
				Description: "instances: 2",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "9c3f59c0-1fdd-470d-91ac-6edbd7a8272d",
				Name:        "audit.app.process.crash",
				Timestamp:   time.Unix(1488561411, 0),
				Description: "index: 1, reason: CRASHED, exit_description: 2 error(s) occurred:\n\n* 2 error(s) occurred:\n\n* Exited with status 148\n* cancelled\n* cancelled",
				Actor:       "19b9d70b-6ebe-47d7-9313-f0c213445036",
				ActorName:   "web",
			},
			{
				GUID:        "5ca51bea-6a32-49ed-94eb-8fa56de9e9bd",
				Name:        "app.crash",
				Timestamp:   time.Unix(1488561411, 0),
				Description: "index: 1, reason: CRASHED, exit_description: 2 error(s) occurred:\n\n* 2 error(s) occurred:\n\n* Exited with status 148\n* cancelled\n* cancelled",
				Actor:       "19b9d70b-6ebe-47d7-9313-f0c213445036",
				ActorName:   "app1",
			},
			{
				GUID:        "1d1d253a-2fbb-4ee6-9a86-6e1beb740178",
				Name:        "audit.app.process.crash",
				Timestamp:   time.Unix(1488561417, 0),
				Description: "index: 1, reason: CRASHED, exit_description: 2 error(s) occurred:\n\n* 2 error(s) occurred:\n\n* Exited with status 148\n* cancelled\n* cancelled",
				Actor:       "19b9d70b-6ebe-47d7-9313-f0c213445036",
				ActorName:   "web",
			},
			{
				GUID:        "47432f69-4a29-44f6-b68b-1d833b4a619d",
				Name:        "app.crash",
				Timestamp:   time.Unix(1488561417, 0),
				Description: "index: 1, reason: CRASHED, exit_description: 2 error(s) occurred:\n\n* 2 error(s) occurred:\n\n* Exited with status 148\n* cancelled\n* cancelled",
				Actor:       "19b9d70b-6ebe-47d7-9313-f0c213445036",
				ActorName:   "app1",
			},
			{
				GUID:        "060eeba0-46ee-4adc-9a42-0c3af57b3510",
				Name:        "audit.app.map-route",
				Timestamp:   time.Unix(1488561422, 0),
				Description: "",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "42538b45-0446-4b6c-9f25-dbbda3c3d0d7",
				Name:        "audit.app.update",
				Timestamp:   time.Unix(1488561422, 0),
				Description: "",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "6e475d88-e229-4605-9dc8-54a0981d4aa2",
				Name:        "audit.app.process.crash",
				Timestamp:   time.Unix(1488561423, 0),
				Description: "index: 1, reason: CRASHED, exit_description: 2 error(s) occurred:\n\n* 2 error(s) occurred:\n\n* Exited with status 148\n* cancelled\n* cancelled",
				Actor:       "19b9d70b-6ebe-47d7-9313-f0c213445036",
				ActorName:   "web",
			},
			{
				GUID:        "e8d4f178-6007-47f6-acbe-d4b903b152a4",
				Name:        "app.crash",
				Timestamp:   time.Unix(1488561423, 0),
				Description: "index: 1, reason: CRASHED, exit_description: 2 error(s) occurred:\n\n* 2 error(s) occurred:\n\n* Exited with status 148\n* cancelled\n* cancelled",
				Actor:       "19b9d70b-6ebe-47d7-9313-f0c213445036",
				ActorName:   "app1",
			},
			{
				GUID:        "72c47ca0-ccb1-44ad-a0b4-eea6559217f2",
				Name:        "audit.app.update",
				Timestamp:   time.Unix(1488561445, 0),
				Description: "instances: 1, memory: 512, environment_json: PRIVATE DATA HIDDEN",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "c398071e-f3aa-4ddc-a049-eeed04cdc855",
				Name:        "audit.app.update",
				Timestamp:   time.Unix(1488561467, 0),
				Description: "state: STOPPED",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "1cd55491-bcb8-4cbc-96b9-cd7476ba4db3",
				Name:        "audit.app.update",
				Timestamp:   time.Unix(1488561468, 0),
				Description: "instances: 1, memory: 512, environment_json: PRIVATE DATA HIDDEN",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "eb3ed27f-8dbd-4a5d-8b4d-12de1f15f9f7",
				Name:        "audit.app.update",
				Timestamp:   time.Unix(1488561472, 0),
				Description: "state: STARTED",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "6a1db406-5fd9-47ce-827e-bedabf312edf",
				Name:        "audit.app.update",
				Timestamp:   time.Unix(1488561494, 0),
				Description: "state: STARTED",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
		},
	},
	"e02a2b91-3b12-4153-971c-73c0243fe8e4": {
		GUID: "e02a2b91-3b12-4153-971c-73c0243fe8e4",
		Name: "mksapp1",
		Type: "route",
		EventList: []models.EventFields{
			{
				GUID:        "6af58a8c-632a-40a5-bcab-10a6841d2d08",
				Name:        "audit.route.create",
				Timestamp:   time.Unix(1488561140, 0),
				Description: "",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
		},
	},
	"62aa7e42-f047-46de-98f1-ff4818a82d6e": {
		GUID: "62aa7e42-f047-46de-98f1-ff4818a82d6e",
		Name: "app2",
		Type: "app",
		EventList: []models.EventFields{
			{
				GUID:        "cf1cb714-e4ba-4894-b8ff-62ac9a62b53c",
				Name:        "audit.app.create",
				Timestamp:   time.Unix(1488561176, 0),
				Description: "instances: 1, memory: 512, state: STOPPED, environment_json: PRIVATE DATA HIDDEN",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "7625787c-a0d8-43e2-9815-3da0d4606543",
				Name:        "audit.app.map-route",
				Timestamp:   time.Unix(1488561181, 0),
				Description: "",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "2ee48170-f40a-41d1-b5a8-aa2d79d8305c",
				Name:        "audit.app.update",
				Timestamp:   time.Unix(1488561182, 0),
				Description: "",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "d6234621-97b3-41ed-ae44-ab01c3b3bf66",
				Name:        "audit.app.update",
				Timestamp:   time.Unix(1488561199, 0),
				Description: "state: STARTED",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "3a2664ef-0613-464b-9b25-bed259e7b62e",
				Name:        "audit.app.map-route",
				Timestamp:   time.Unix(1488561299, 0),
				Description: "",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "32291e3a-9d3f-45b6-8baf-0d1b550726ef",
				Name:        "audit.app.update",
				Timestamp:   time.Unix(1488561299, 0),
				Description: "",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "4b39d8c4-c54e-4816-ad8b-ba3541d05538",
				Name:        "audit.app.update",
				Timestamp:   time.Unix(1488561314, 0),
				Description: "instances: 1, memory: 512, environment_json: PRIVATE DATA HIDDEN",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "1466976d-2c1a-4ddd-a08b-b54323175bc2",
				Name:        "audit.app.update",
				Timestamp:   time.Unix(1488561336, 0),
				Description: "state: STOPPED",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
			{
				GUID:        "d3daff16-fc03-4b2b-bf43-a3d9efb2e491",
				Name:        "audit.app.update",
				Timestamp:   time.Unix(1488561339, 0),
				Description: "state: STARTED",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
		},
	},
	"06f95896-45aa-4178-91e4-b25c9c9c882b": {
		GUID: "06f95896-45aa-4178-91e4-b25c9c9c882b",
		Name: "mksapp2",
		Type: "route",
		EventList: []models.EventFields{
			{
				GUID:        "bc1f8e4a-1d51-4ba8-be40-0e3c975c1aa1",
				Name:        "audit.route.create",
				Timestamp:   time.Unix(1488561180, 0),
				Description: "",
				Actor:       "5b3fedc4-22e9-4276-85e1-f16d60330adc",
				ActorName:   "msamaratunga@pivotal.io",
			},
		},
	},
}
