package webmethods

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	hc "github.com/Axway/agent-sdk/pkg/util/healthcheck"
	"github.com/Axway/agent-sdk/pkg/util/log"

	"git.ecd.axway.org/apigov/agents-webmethods/pkg/config"
)

const HealthCheckEndpoint = "webmethods"

// Page describes the page query parameter
type Page struct {
	Offset   int
	PageSize int
}

// Client interface to gateway
type Client interface {
	GetAPI(id string) (*API, error)
	ListAPIs() ([]API, error)
	OnConfigChange(webMethodConfig *config.WebMethodConfig)
}

type ListAPIClient interface {
	ListAPIs() ([]API, error)
}

// WebMethodClient is the client for interacting with Webmethods APIM.
type WebMethodClient struct {
	specPath string
}

// NewClient creates a new client for interacting with Webmethods APIM.
func NewClient(webMethodConfig *config.WebMethodConfig) *WebMethodClient {
	client := &WebMethodClient{}

	client.OnConfigChange(webMethodConfig)

	hc.RegisterHealthcheck("Webmethods API Gateway", HealthCheckEndpoint, client.healthcheck)

	return client
}

func (c *WebMethodClient) OnConfigChange(webMethodConfig *config.WebMethodConfig) {
	c.specPath = ""
}

func (c *WebMethodClient) healthcheck(name string) (status *hc.Status) {
	status = &hc.Status{
		Result: hc.OK,
	}

	return status
}

func (c *WebMethodClient) listSpecFiles() []string {
	var files []string
	filepath.Walk(c.specPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files
}

// ListAPIs lists the API.
func (c *WebMethodClient) ListAPIs() ([]API, error) {
	specFiles := c.listSpecFiles()
	apis := make([]API, 0)
	for _, specFile := range specFiles {
		apiName, apiSpec, err := c.getSpec(specFile)
		if err != nil {
			log.Infof("Failed to load sample API specification from %s: %s ", c.specPath, err.Error())
		}
		api := API{
			ID:            apiName,
			Name:          apiName,
			Description:   specFile,
			Version:       "1.0.0",
			Url:           "",
			Documentation: []byte(specFile),
			ApiSpec:       apiSpec,
		}
		apis = append(apis, api)
	}

	return apis, nil
}

func (a *WebMethodClient) getSpec(specFile string) (string, []byte, error) {
	fileName := filepath.Base(specFile)
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

	bytes, err := ioutil.ReadFile(specFile)
	if err != nil {
		return "", nil, err
	}
	return fileName, bytes, nil
}

// GetAPI gets a single api by id
func (c *WebMethodClient) GetAPI(id string) (*API, error) {

	return nil, nil
}
