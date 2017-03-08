package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mevansam/cf-event-resource-type/resource"
	"github.com/mitchellh/colorstring"
)

func main() {

	fmt.Fprintf(os.Stderr, colorstring.Color(
		"[yellow]Running 'in' command to create Cloud Foundry event content for concourse job...\n"))

	var request resource.InRequest
	inputRequest(&request)

	if request.Source.Debug {
		b, err := json.MarshalIndent(request, "", "  ")
		if err != nil {
			resource.Fatalf("[red]Marshalling JSON input for debugging")
		}
		fmt.Fprintf(os.Stderr, colorstring.Color("[green]In command input:\n%s\n"), string(b))
	}

	if len(os.Args) < 2 {
		resource.Fatalf("[red]the destination directory was not passed in as the first argument")
	}

	downloadDir := os.Args[1]
	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		resource.Fatalf("[red]the destination directory '%s' does not exist", downloadDir)
	}

	command := resource.NewInCommand(downloadDir)
	response, err := command.Run(request)
	if err != nil {
		resource.Fatalf("[red]running command: %s\n", err)
	}

	if request.Source.Debug {
		b, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			resource.Fatalf("[red]Marshalling JSON response for debugging")
		}
		fmt.Fprintf(os.Stderr, colorstring.Color("[green]In command response:\n%s\n"), string(b))
	}

	outputResponse(response)
}

func inputRequest(request *resource.InRequest) {
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
		resource.Fatalf("[red]error reading request from stdin: %s\n", err)
	}
}

func outputResponse(response resource.InResponse) {
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		resource.Fatalf("[red]error writing response to stdout: %s\n", err)
	}
}
