package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mevansam/cf-cli-api/cfapi"
	"github.com/mevansam/cf-event-resource-type/resource"
	"github.com/mitchellh/colorstring"
)

func main() {

	fmt.Printf(colorstring.Color("[yellow]Running 'check' command to retrieve Cloud Foundry Events...\n"))

	var request resource.CheckRequest
	inputRequest(&request)

	if request.Source.Debug {
		b, err := json.MarshalIndent(request, "", "  ")
		if err != nil {
			resource.Fatalf("[red]Marshalling JSON input for debugging")
		}
		fmt.Printf(colorstring.Color("[green]Check command input:\n%s\n"), string(b))
	}

	cfClient, err := resource.NewCfClient(cfapi.NewCfCliSessionProvider(), request.Source)
	if err != nil {
		resource.Fatalf("[red]connecting to CF endpoint %s: %s\n", request.Source.API, err.Error())
	}

	command := resource.NewCheckCommand(cfClient)
	response, err := command.Run(request)
	if err != nil {
		resource.Fatalf("[red]running command: %s\n", err)
	}

	if request.Source.Debug {
		b, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			resource.Fatalf("[red]Marshalling JSON response for debugging")
		}
		fmt.Printf(colorstring.Color("[green]Check command response:\n%s\n"), string(b))
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
