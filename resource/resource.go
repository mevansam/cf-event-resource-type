package resource

// CheckRequest -
type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

// InRequest -
type InRequest struct {
	Source  Source   `json:"source"`
	Version *Version `json:"version"`
	Params  InParams `json:"params"`
}

// InParams -
type InParams struct {
}

// InResponse -
type InResponse struct {
	Version  Version        `json:"version"`
	Metadata []MetadataPair `json:"metadata"`
}

// OutRequest -
type OutRequest struct {
	Source Source    `json:"source"`
	Params OutParams `json:"params"`
}

// OutParams -
type OutParams struct {
}

// OutResponse -
type OutResponse struct {
	Version  Version        `json:"version"`
	Metadata []MetadataPair `json:"metadata"`
}

// Source -
type Source struct {
	API      string `json:"api"`
	User     string `json:"user"`
	Password string `json:"password"`
	SSOToken string `json:"sso-token"`

	Org   string   `json:"org"`
	Space string   `json:"space"`
	Apps  []string `json:"apps"`

	SkipSSLValidation bool `json:"skip-ssl-validation"`

	Debug bool `json:"debug"`
	Trace bool `json:"trace"`
}

// Version - { "app name": (serialized filters.AppEvent), ... }
type Version map[string]string

// MetadataPair -
type MetadataPair struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	URL      string `json:"url"`
	Markdown bool   `json:"markdown"`
}

// NewCheckRequest -
func NewCheckRequest() CheckRequest {
	res := CheckRequest{
		Source: Source{
			SkipSSLValidation: false,
			Debug:             false,
			Trace:             false,
		},
	}
	return res
}

// NewInRequest -
func NewInRequest() InRequest {
	res := InRequest{}
	return res
}

// NewOutRequest -
func NewOutRequest() InRequest {
	res := InRequest{}
	return res
}
