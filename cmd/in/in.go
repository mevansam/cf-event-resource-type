package main

import (
	"encoding/json"
	"os"

	"github.com/mevansam/cf-event-resource-type/resource"
)

func main() {
	var request resource.InRequest
	inputRequest(&request)

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
