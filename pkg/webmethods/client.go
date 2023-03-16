package webmethods

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	coreapi "github.com/Axway/agent-sdk/pkg/api"
	hc "github.com/Axway/agent-sdk/pkg/util/healthcheck"

	"git.ecd.axway.org/apigov/agents-webmethods/pkg/config"
)

const HealthCheckEndpoint = "health"

// Page describes the page query parameter
type Page struct {
	Offset   int
	PageSize int
}

// Client interface to gateway
type Client interface {
	createAuthToken() string
	ListAPIs() ([]WebmethodsApi, error)
	GetApiDetails(id string) (*ApiResponse, error)
	GetApiSpec(id string) ([]byte, error)
	GetWsdl(gatewayEndpoint string) ([]byte, error)
	OnConfigChange(webMethodConfig *config.WebMethodConfig)
}

// WebMethodClient is the client for interacting with Webmethods APIM.
type WebMethodClient struct {
	url        string
	username   string
	password   string
	httpClient coreapi.Client
}

// NewClient creates a new client for interacting with Webmethods APIM.
func NewClient(webMethodConfig *config.WebMethodConfig) *WebMethodClient {
	client := &WebMethodClient{}
	client.OnConfigChange(webMethodConfig)
	hc.RegisterHealthcheck("Webmethods API Gateway", HealthCheckEndpoint, client.healthcheck)
	return client
}

func (c *WebMethodClient) OnConfigChange(webMethodConfig *config.WebMethodConfig) {
	c.url = webMethodConfig.WebmethodsApimUrl
	c.username = webMethodConfig.Username
	c.password = webMethodConfig.Password
}

func (c *WebMethodClient) healthcheck(name string) (status *hc.Status) {
	status = &hc.Status{
		Result: hc.OK,
	}
	return status
}

// ListAPIs lists webmethods  APIM apis.
func (c *WebMethodClient) ListAPIs() ([]WebmethodsApi, error) {
	webmethodsApis := make([]WebmethodsApi, 0)
	url := fmt.Sprintf("%s/rest/apigateway/apis", c.url)
	query := map[string]string{
		"isActive": "true",
	}
	headers := map[string]string{
		"Authorization": c.createAuthToken(),
	}
	request := coreapi.Request{
		Method:      coreapi.GET,
		URL:         url,
		Headers:     headers,
		QueryParams: query,
	}

	response, err := c.httpClient.Send(request)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response.Body, webmethodsApis)
	if err != nil {
		return nil, err
	}
	return webmethodsApis, nil
}

// ListAPIs lists webmethods  APIM apis.
func (c *WebMethodClient) GetApiDetails(id string) (*ApiResponse, error) {
	apiResponse := &ApiResponse{}
	url := fmt.Sprintf("%s/rest/apigateway/apis/%s", c.url, id)
	headers := map[string]string{
		"Authorization": c.createAuthToken(),
	}
	request := coreapi.Request{
		Method:  coreapi.GET,
		URL:     url,
		Headers: headers,
	}

	response, err := c.httpClient.Send(request)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response.Body, apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, nil
}

// GetAPI gets a single api by id
func (c *WebMethodClient) GetApiSpec(id string) ([]byte, error) {

	url := fmt.Sprintf("%s/rest/apigateway/apis/%s", c.url, id)
	query := map[string]string{
		"format": "openapi",
	}
	headers := map[string]string{
		"Authorization": c.createAuthToken(),
	}
	request := coreapi.Request{
		Method:      coreapi.GET,
		URL:         url,
		Headers:     headers,
		QueryParams: query,
	}

	response, err := c.httpClient.Send(request)
	if err != nil {
		return nil, err
	}

	return response.Body, nil

}

func (c *WebMethodClient) GetWsdl(gatewayEndpoint string) ([]byte, error) {

	url := gatewayEndpoint + "?wsdl"
	headers := map[string]string{
		"Authorization": c.createAuthToken(),
	}
	request := coreapi.Request{
		Method:  coreapi.GET,
		URL:     url,
		Headers: headers,
	}

	response, err := c.httpClient.Send(request)
	if err != nil {
		return nil, err
	}

	return response.Body, nil

}

func (c *WebMethodClient) createAuthToken() string {
	credential := c.username + ":" + c.password
	return base64.StdEncoding.EncodeToString([]byte(credential))
}
