package main

import (
	"encoding/json"
	"os"

	"github.com/mevansam/cf-cli-api/cfapi"
	"github.com/mevansam/cf-event-resource-type/resource"
)

func main() {
	var request resource.CheckRequest
	inputRequest(&request)

	cfClient, err := resource.NewCfClient(cfapi.NewCfCliSessionProvider(), request.Source)
	if err != nil {
		resource.Fatalf("[red]connecting to CF endpoint %s: %s\n", request.Source.API, err.Error())
	}

	command := resource.NewCheckCommand(cfClient)
	response, err := command.Run(request)
	if err != nil {
		resource.Fatalf("[red]running command: %s\n", err)
	}

	outputResponse(response)
}

func inputRequest(request *resource.CheckRequest) {
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
		resource.Fatalf("[red]error reading request from stdin: %s\n", err)
	}
}

func outputResponse(response []resource.Version) {
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		resource.Fatalf("[red]error writing response to stdout: %s\n", err)
	}
}
