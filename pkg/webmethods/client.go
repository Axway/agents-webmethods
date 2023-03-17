package webmethods

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

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
	ListAPIs() ([]ListApiResponse, error)
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
	client.httpClient = coreapi.NewClient(webMethodConfig.TLS, webMethodConfig.ProxyURL)
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
	url := c.url + "/rest/apigateway/health"
	fmt.Println(url)
	status = &hc.Status{
		Result: hc.OK,
	}

	request := coreapi.Request{
		Method: coreapi.GET,
		URL:    url,
	}
	response, err := c.httpClient.Send(request)

	if err != nil {
		status = &hc.Status{
			Result:  hc.FAIL,
			Details: fmt.Sprintf("%s Failed. Unable to connect to Boomi, check Boomi configuration. %s", name, err.Error()),
		}
	}

	if response.Code != http.StatusOK {
		status = &hc.Status{
			Result:  hc.FAIL,
			Details: fmt.Sprintf("%s Failed. Unable to connect to Boomi, check Boomi configuration.", name),
		}
	}

	return status

}

// ListAPIs lists webmethods  APIM apis.
func (c *WebMethodClient) ListAPIs() ([]ListApiResponse, error) {
	//webmethodsApis := make([]WebmethodsApi, 0)
	url := fmt.Sprintf("%s/rest/apigateway/apis", c.url)
	query := map[string]string{
		"isActive": "true",
	}
	headers := map[string]string{
		"Authorization": c.createAuthToken(),
		"Accept":        "application/json",
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
	listApi := &ListApi{}
	err = json.Unmarshal(response.Body, &listApi)
	if err != nil {
		return nil, err
	}
	return listApi.ListApiResponse, nil
}

// ListAPIs lists webmethods  APIM apis.
func (c *WebMethodClient) GetApiDetails(id string) (*ApiResponse, error) {
	getApiDetails := &GetApiDetails{}
	url := fmt.Sprintf("%s/rest/apigateway/apis/%s", c.url, id)
	headers := map[string]string{
		"Authorization": c.createAuthToken(),
		"Accept":        "application/json",
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

	err = json.Unmarshal(response.Body, getApiDetails)
	if err != nil {
		return nil, err
	}
	return &getApiDetails.ApiResponse, nil
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
	return "Basic :" + base64.StdEncoding.EncodeToString([]byte(credential))
}
