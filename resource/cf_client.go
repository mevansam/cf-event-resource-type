package resource

import (
	"os"

	"github.com/mevansam/cf-cli-api/cfapi"
)

// CfClient -
type CfClient struct {
	session cfapi.CfSession
}

// NewCfClient -
func NewCfClient(provider cfapi.CfSessionProvider, src Source) (*CfClient, error) {

	cfClient := &CfClient{}

	var (
		err error

		logger *cfapi.Logger
	)

	if src.Trace {
		logger = cfapi.NewFileLogger(src.Debug, "true", os.Stdin, os.Stderr)
	} else {
		logger = cfapi.NewFileLogger(src.Debug, "", os.Stdin, os.Stderr)
	}

	cfClient.session, err = provider.NewCfSession(
		src.API, src.User, src.Password,
		src.Org, src.Space,
		src.SkipSSLValidation, logger)

	if err != nil {
		return nil, err
	}
	return cfClient, nil
}
