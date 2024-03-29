package webmethods

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	coreapi "github.com/Axway/agent-sdk/pkg/api"
	"github.com/Axway/agent-sdk/pkg/filter"
	agenterrors "github.com/Axway/agent-sdk/pkg/util/errors"
	hc "github.com/Axway/agent-sdk/pkg/util/healthcheck"

	"github.com/Axway/agents-webmethods/pkg/config"
)

const HealthCheckEndpoint = "health"
const searchURL = "%s/rest/apigateway/search"
const getApplicationURL = "%s/rest/apigateway/applications/%s"

// Page describes the page query parameter
type Page struct {
	Offset   int
	PageSize int
}

var (
	token string
)

// Client interface to gateway
type Client interface {
	createAuthToken() string
	ListAPIs() ([]ListApiResponse, error)
	SearchAPIs() (*Apis, error)
	GetApiDetails(id string) (*ApiResponse, error)
	IsAllowedTags(tags []Tag) bool
	GetApiSpec(id string) ([]byte, error)
	GetWsdl(gatewayEndpoint string) ([]byte, error)
	FindApplicationByName(applicationName string) (*SearchApplicationResponse, error)
	CreateApplication(application *Application) (*Application, error)
	UpdateApplication(application *Application) (*Application, error)
	SubscribeApplication(applicationId string, ApplicationApiSubscription *ApplicationApiSubscription) error
	GetApplication(applicationId string) (*ApplicationResponse, error)
	RotateApplicationApikey(applicationId string) error
	CreateOauth2Strategy(strategy *Strategy) (*StrategyResponse, error)
	DeleteStrategy(strategyId string) error
	RefereshOauth2Credential(strategyId string) (*StrategyResponse, error)
	GetStrategy(strategyId string) (*StrategyResponse, error)
	DeleteApplication(applicationId string) error
	OnConfigChange(webMethodConfig *config.WebMethodConfig) error
	DeleteApplicationAccessTokens(applicationId string) error
	UnsubscribeApplication(applicationId string, apiId string) error
	ListOauth2Servers() (*OauthServers, error)
	GetTransactionsWindow(startDate, endDate string) ([]byte, error)
	Healthcheck(name string) (status *hc.Status)
}

// WebMethodClient is the client for interacting with Webmethods APIM.
type WebMethodClient struct {
	url string

	username        string
	password        string
	httpClient      coreapi.Client
	discoveryFilter filter.Filter
}

// NewClient creates a new client for interacting with Webmethods APIM.
func NewClient(webMethodConfig *config.WebMethodConfig, httpClient coreapi.Client) (*WebMethodClient, error) {
	client := &WebMethodClient{}
	client.httpClient = httpClient
	err := client.OnConfigChange(webMethodConfig)
	hc.RegisterHealthcheck("Webmethods API Gateway", HealthCheckEndpoint, client.Healthcheck)
	return client, err
}

func (c *WebMethodClient) OnConfigChange(webMethodConfig *config.WebMethodConfig) error {
	c.url = webMethodConfig.WebmethodsApimUrl
	c.username = webMethodConfig.Username
	c.password = webMethodConfig.Password
	c.discoveryFilter = nil
	if strings.TrimSpace(webMethodConfig.Filter) != "" {

		newFilter, err := filter.NewFilter(webMethodConfig.Filter)
		if err != nil {
			return err
		}
		c.discoveryFilter = newFilter
	}
	return nil
}

func (c *WebMethodClient) Healthcheck(name string) (status *hc.Status) {
	url := c.url + "/rest/apigateway/health"
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
			Details: fmt.Sprintf("%s Failed. Unable to connect to Webmethods, check Webmethods configuration. %s", name, err.Error()),
		}
		return status
	}
	if response.Code != http.StatusOK {
		status = &hc.Status{
			Result:  hc.FAIL,
			Details: fmt.Sprintf("%s Failed. Unable to connect to Webmethods, check Boomi configuration.", name),
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

func (c *WebMethodClient) SearchAPIs() (*Apis, error) {
	//webmethodsApis := make([]WebmethodsApi, 0)
	url := fmt.Sprintf(searchURL, c.url)
	requestStr := `{
		"types": [
			"api"
		],
		"condition": "and",
		"scope": [
			{
				"attributeName": "isActive",
				"keyword": true
			}
		],
		"sortByField": "apiName",
		"sortOrder": "ASC"
	}`
	apis := &Apis{}
	headers := map[string]string{
		"Authorization": c.createAuthToken(),
		"Content-Type":  "application/json",
	}

	request := coreapi.Request{
		Method:  coreapi.POST,
		URL:     url,
		Headers: headers,
		Body:    []byte(requestStr),
	}
	response, err := c.httpClient.Send(request)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response.Body, apis)
	if err != nil {
		return nil, err
	}
	return apis, nil
}

// GetApiDetails ListAPIs lists webmethods  APIM apis.
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

func (c *WebMethodClient) IsAllowedTags(tags []Tag) bool {
	if c.discoveryFilter != nil {
		tagsMap := make(map[string]interface{})
		for _, value := range tags {
			tagsMap[value.Name] = ""
		}
		return c.discoveryFilter.Evaluate(tagsMap)
	}
	return true
}

// GetApiSpec  gets a single api by id
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

func (c *WebMethodClient) FindApplicationByName(applicationName string) (*SearchApplicationResponse, error) {
	searchRequest := &Search{}
	searchRequest.Types = []string{"APPLICATION"}
	searchRequest.ResponseFields = []string{"applicationID", "name"}
	applicationResponse := &SearchApplicationResponse{}
	scope := Scope{}
	scope.AttributeName = "name"
	scope.Keyword = applicationName
	searchRequest.Scope = []Scope{scope}
	url := fmt.Sprintf(searchURL, c.url)
	headers := map[string]string{
		"Authorization": c.createAuthToken(),
		"Content-Type":  "application/json",
	}
	buffer, err := json.Marshal(searchRequest)
	if err != nil {
		return nil, agenterrors.Newf(2000, err.Error())
	}
	request := coreapi.Request{
		Method:  coreapi.POST,
		URL:     url,
		Headers: headers,
		Body:    buffer,
	}
	response, err := c.httpClient.Send(request)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response.Body, applicationResponse)
	if err != nil {
		return nil, err
	}
	return applicationResponse, nil
}

func (c *WebMethodClient) GetApplication(applicationId string) (*ApplicationResponse, error) {
	applicationResponse := &ApplicationResponse{}
	url := fmt.Sprintf(getApplicationURL, c.url, applicationId)
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

	err = json.Unmarshal(response.Body, applicationResponse)
	if err != nil {
		return nil, err
	}
	return applicationResponse, nil
}

func (c *WebMethodClient) CreateApplication(application *Application) (*Application, error) {
	responseApplication := &Application{}
	url := fmt.Sprintf("%s/rest/apigateway/applications", c.url)
	headers := map[string]string{
		"Authorization": c.createAuthToken(),
		"Content-Type":  "application/json",
	}
	buffer, err := json.Marshal(application)
	if err != nil {
		return nil, agenterrors.Newf(2000, err.Error())
	}
	request := coreapi.Request{
		Method:  coreapi.POST,
		URL:     url,
		Headers: headers,
		Body:    buffer,
	}

	response, err := c.httpClient.Send(request)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response.Body, responseApplication)
	return responseApplication, nil
}

func (c *WebMethodClient) UpdateApplication(application *Application) (*Application, error) {
	responseApplication := &Application{}
	url := fmt.Sprintf(getApplicationURL, c.url, application.Id)
	headers := map[string]string{
		"Authorization": c.createAuthToken(),
		"Content-Type":  "application/json",
	}
	buffer, err := json.Marshal(application)
	if err != nil {
		return nil, agenterrors.Newf(2000, err.Error())
	}
	request := coreapi.Request{
		Method:  coreapi.PUT,
		URL:     url,
		Headers: headers,
		Body:    buffer,
	}

	response, err := c.httpClient.Send(request)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response.Body, responseApplication)
	return responseApplication, nil
}

func (c *WebMethodClient) CreateOauth2Strategy(strategy *Strategy) (*StrategyResponse, error) {
	strategyResponse := &StrategyResponse{}

	url := fmt.Sprintf("%s/rest/apigateway/strategies", c.url)
	headers := map[string]string{
		"Authorization": c.createAuthToken(),
		"Content-Type":  "application/json",
	}
	buffer, err := json.Marshal(strategy)
	if err != nil {
		return nil, agenterrors.Newf(2000, err.Error())
	}
	request := coreapi.Request{
		Method:  coreapi.POST,
		URL:     url,
		Headers: headers,
		Body:    buffer,
	}

	response, err := c.httpClient.Send(request)
	if err != nil {
		return nil, err
	}

	if response.Code == 201 {
		err = json.Unmarshal(response.Body, strategyResponse)
		return strategyResponse, nil
	}

	return nil, agenterrors.Newf(2000, "Unable to create strategy")
}

func (c *WebMethodClient) GetStrategy(strategyId string) (*StrategyResponse, error) {
	strategyResponse := &StrategyResponse{}
	url := fmt.Sprintf("%s/rest/apigateway/strategies/%s", c.url, strategyId)
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

	err = json.Unmarshal(response.Body, strategyResponse)
	if err != nil {
		return nil, err
	}
	return strategyResponse, nil
}

func (c *WebMethodClient) SubscribeApplication(applicationId string, ApplicationApiSubscription *ApplicationApiSubscription) error {
	url := fmt.Sprintf(getApplicationURL+"/apis", c.url, applicationId)
	headers := map[string]string{
		"Authorization": c.createAuthToken(),
		"Content-Type":  "application/json",
	}
	buffer, err := json.Marshal(ApplicationApiSubscription)
	if err != nil {
		return agenterrors.Newf(2000, err.Error())
	}
	request := coreapi.Request{
		Method:  coreapi.POST,
		URL:     url,
		Headers: headers,
		Body:    buffer,
	}

	response, err := c.httpClient.Send(request)
	if err != nil {
		return err
	}
	if response.Code != 201 {
		return agenterrors.Newf(2001, "Unable to assicate API to Application")
	}
	return nil
}

func (c *WebMethodClient) RotateApplicationApikey(applicationId string) error {
	url := fmt.Sprintf(getApplicationURL+"/accessTokens", c.url, applicationId)
	headers := map[string]string{
		"Authorization": c.createAuthToken(),
		"Content-Type":  "application/json",
	}
	var jsonBody = []byte(`{ "type": "apiAccessKeyCredentials"}`)
	request := coreapi.Request{
		Method:  coreapi.POST,
		URL:     url,
		Headers: headers,
		Body:    jsonBody,
	}

	response, err := c.httpClient.Send(request)
	if err != nil {
		return err
	}
	if response.Code != 201 {
		return agenterrors.Newf(2001, "Unable to Rotate API Key")
	}
	return nil
}

func (c *WebMethodClient) DeleteStrategy(strategyId string) error {
	url := fmt.Sprintf("%s/rest/apigateway/strategies/%s", c.url, strategyId)
	headers := map[string]string{
		"Authorization": c.createAuthToken(),
		"Content-Type":  "application/json",
	}
	request := coreapi.Request{
		Method:  coreapi.DELETE,
		URL:     url,
		Headers: headers,
	}

	response, err := c.httpClient.Send(request)
	if err != nil {
		return err
	}
	if response.Code != 204 {
		return agenterrors.Newf(2001, "Unable to Delete Stratgey")
	}
	return nil
}

func (c *WebMethodClient) RefereshOauth2Credential(strategyId string) (*StrategyResponse, error) {
	strategyResponse := &StrategyResponse{}

	url := fmt.Sprintf("%s/rest/apigateway/strategies/%s/refreshCredentials", c.url, strategyId)
	headers := map[string]string{
		"Authorization": c.createAuthToken(),
		"Content-Type":  "application/json",
	}
	request := coreapi.Request{
		Method:  coreapi.PUT,
		URL:     url,
		Headers: headers,
	}

	response, err := c.httpClient.Send(request)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response.Body, strategyResponse)
	if err != nil {
		return nil, err
	}
	return strategyResponse, nil
}

func (c *WebMethodClient) DeleteApplication(applicationId string) error {
	url := fmt.Sprintf(getApplicationURL, c.url, applicationId)
	headers := map[string]string{
		"Authorization": c.createAuthToken(),
		"Content-Type":  "application/json",
	}
	request := coreapi.Request{
		Method:  coreapi.DELETE,
		URL:     url,
		Headers: headers,
	}

	response, err := c.httpClient.Send(request)
	if err != nil {
		return err
	}
	if response.Code != 204 {
		return agenterrors.Newf(2001, "Unable to Delete Application")
	}
	return nil
}

func (c *WebMethodClient) DeleteApplicationAccessTokens(applicationId string) error {
	url := fmt.Sprintf(getApplicationURL+"/accessTokens", c.url, applicationId)
	headers := map[string]string{
		"Authorization": c.createAuthToken(),
		"Content-Type":  "application/json",
	}
	var jsonBody = []byte(`{ "type": "accessTokens"}`)
	request := coreapi.Request{
		Method:  coreapi.DELETE,
		URL:     url,
		Headers: headers,
		Body:    jsonBody,
	}

	response, err := c.httpClient.Send(request)
	if err != nil {
		return err
	}
	if response.Code != 204 {
		return agenterrors.Newf(2001, "Unable to Delete Api Key / Oauth tokens")
	}
	return nil
}

func (c *WebMethodClient) UnsubscribeApplication(applicationId string, apiId string) error {
	url := fmt.Sprintf(getApplicationURL+"/apis", c.url, applicationId)
	headers := map[string]string{
		"Authorization": c.createAuthToken(),
		"Content-Type":  "application/json",
	}
	queryString := map[string]string{
		"apiIDs": apiId,
	}

	request := coreapi.Request{
		Method:      coreapi.DELETE,
		URL:         url,
		Headers:     headers,
		QueryParams: queryString,
	}

	response, err := c.httpClient.Send(request)
	if err != nil {
		return err
	}
	if response.Code != 204 {
		return agenterrors.Newf(2001, "Unable to remove API to Application")
	}
	return nil
}

func (c *WebMethodClient) ListOauth2Servers() (*OauthServers, error) {
	requestStr := `{
		"types": [
			"alias"
		],
		"scope": [
			{
				"attributeName": "type",
				"keyword": "authServerAlias"
			}
		],
		"responseFields": [
			"id",
			"name",
			"type",
			"description",
			"scopes"
		],
		"condition": "or",
		"sortByField": "name"
	}`
	oauthServers := &OauthServers{}
	url := fmt.Sprintf(searchURL, c.url)
	headers := map[string]string{
		"Authorization": c.createAuthToken(),
		"Content-Type":  "application/json",
	}

	request := coreapi.Request{
		Method:  coreapi.POST,
		URL:     url,
		Headers: headers,
		Body:    []byte(requestStr),
	}
	response, err := c.httpClient.Send(request)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response.Body, oauthServers)
	if err != nil {
		return nil, err
	}
	return oauthServers, nil
}

func (c *WebMethodClient) GetTransactionsWindow(startDate, endDate string) ([]byte, error) {
	query := map[string]string{
		"eventType": "transactionalEvents",
		"startDate": startDate,
		"endDate":   endDate,
	}
	headers := map[string]string{
		"Authorization": c.createAuthToken(),
	}
	url := fmt.Sprintf("%s/rest/apigateway/apitransactions", c.url)
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

func (c *WebMethodClient) createAuthToken() string {
	if token == "" {
		credential := c.username + ":" + c.password
		token = "Basic :" + base64.StdEncoding.EncodeToString([]byte(credential))
	}
	return token
}
