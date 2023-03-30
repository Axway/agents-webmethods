package webmethods

import (
	"testing"

	"git.ecd.axway.org/apigov/agents-webmethods/pkg/config"
	coreapi "github.com/Axway/agent-sdk/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestListApi(t *testing.T) {

	response := `{
		"apiResponse": [
			{
				"api": {
					"apiName": "petstore",
					"apiVersion": "1.0.17",
					"apiDescription": "This is a sample Pet Store Server based on the OpenAPI 3.0 specification.  You can find out more about\nSwagger at [http://swagger.io](http://swagger.io). In the third iteration of the pet store, we've switched to the design first approach!\nYou can now help us improve the API whether it's by making changes to the definition itself or to the code.\nThat way, with time, we can improve the API in general, and expose some of the new features in OAS3.\n\nSome useful links:\n- [The Pet Store repository](https://github.com/swagger-api/swagger-petstore)\n- [The source API definition for the Pet Store](https://github.com/swagger-api/swagger-petstore/blob/master/src/main/resources/openapi.yaml)",
					"isActive": true,
					"type": "REST",
					"tracingEnabled": false,
					"publishedPortals": [
						"69ca5c6b-8b0a-4bae-8021-4ccf0ffd136f"
					],
					"systemVersion": 1,
					"id": "2b598e47-3e0c-4b0f-8a72-da7fdf6c5ea2"
				},
				"responseStatus": "SUCCESS"
			}
		]
	}`
	mc := &MockClient{}
	cfg := &config.WebMethodConfig{
		WebmethodsApimUrl: "",
		CachePath:         "/tmp",
		Password:          "abc",
		PollInterval:      10,
		ProxyURL:          "",
		Username:          "123",
	}
	webMethodsClient := NewClient(cfg, mc)

	mc.SendFunc = func(request coreapi.Request) (*coreapi.Response, error) {
		return &coreapi.Response{
			Code: 200,
			Body: []byte(response),
		}, nil
	}

	apiResponse, err := webMethodsClient.ListAPIs()

	assert.Nil(t, err)
	webmethodApi := apiResponse[0].WebmethodsApi
	assert.Equal(t, webmethodApi.Id, "2b598e47-3e0c-4b0f-8a72-da7fdf6c5ea2")
	assert.Equal(t, webmethodApi.ApiName, "petstore")
	assert.Equal(t, webmethodApi.ApiVersion, "1.0.17")
	assert.True(t, webmethodApi.IsActive)
}
