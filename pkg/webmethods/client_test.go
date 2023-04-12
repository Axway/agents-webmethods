package webmethods

import (
	"encoding/json"
	"testing"

	"git.ecd.axway.org/apigov/agents-webmethods/pkg/config"
	coreapi "github.com/Axway/agent-sdk/pkg/api"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var cfg = &config.WebMethodConfig{
	WebmethodsApimUrl: "",
	CachePath:         "/tmp",
	Password:          "abc",
	PollInterval:      10,
	ProxyURL:          "",
	Username:          "123",
}

func TestListApi(t *testing.T) {

	response := `{
		"apiResponse": [
			{
				"api": {
					"apiName": "Calculator",
					"apiVersion": "1",
					"isActive": true,
					"type": "SOAP",
					"tracingEnabled": false,
					"publishedPortals": [],
					"systemVersion": 1,
					"id": "1178fcb2-faae-4fe4-94fa-f5efb0316285"
				},
				"responseStatus": "SUCCESS"
			},
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
	webMethodsClient, _ := NewClient(cfg, mc)
	mc.SendFunc = func(request coreapi.Request) (*coreapi.Response, error) {
		return &coreapi.Response{
			Code: 200,
			Body: []byte(response),
		}, nil
	}
	apiResponse, err := webMethodsClient.ListAPIs()
	assert.Nil(t, err)
	webmethodApi := apiResponse[1].WebmethodsApi
	assert.Equal(t, webmethodApi.Id, "2b598e47-3e0c-4b0f-8a72-da7fdf6c5ea2")
	assert.Equal(t, webmethodApi.ApiName, "petstore")
	assert.Equal(t, webmethodApi.ApiVersion, "1.0.17")
	assert.Equal(t, webmethodApi.ApiType, "REST")

	assert.True(t, webmethodApi.IsActive)

	calculatorApi := apiResponse[0].WebmethodsApi
	assert.Equal(t, calculatorApi.Id, "1178fcb2-faae-4fe4-94fa-f5efb0316285")
	assert.Equal(t, calculatorApi.ApiType, "SOAP")
	assert.True(t, calculatorApi.IsActive)
}

func TestApiDetails(t *testing.T) {

	response := `{
		"apiResponse": {
			"api": {
				"apiDefinition": {
					"info": {
						"description": "This is a sample Pet Store Server based on the OpenAPI 3.0 specification.  You can find out more about\nSwagger at [http://swagger.io](http://swagger.io). In the third iteration of the pet store, we've switched to the design first approach!\nYou can now help us improve the API whether it's by making changes to the definition itself or to the code.\nThat way, with time, we can improve the API in general, and expose some of the new features in OAS3.\n\nSome useful links:\n- [The Pet Store repository](https://github.com/swagger-api/swagger-petstore)\n- [The source API definition for the Pet Store](https://github.com/swagger-api/swagger-petstore/blob/master/src/main/resources/openapi.yaml)",
						"version": "1.0.17",
						"title": "Swagger Petstore - OpenAPI 3.0",
						"termsOfService": "http://swagger.io/terms/",
						"contact": {
							"email": "apiteam@swagger.io"
						},
						"license": {
							"name": "Apache 2.0",
							"url": "http://www.apache.org/licenses/LICENSE-2.0.html"
						}
					},
					"serviceRegistryDisplayName": "petstore_1.0.17",
					"tags": [
						{
							"name": "customrathna"
						},
						{
							"name": "custom2"
						},
						{
							"name": "pet",
							"description": "Everything about your Pets",
							"externalDocs": {
								"description": "Find out more",
								"url": "http://swagger.io"
							}
						},
						{
							"name": "store",
							"description": "Access to Petstore orders",
							"externalDocs": {
								"description": "Find out more about our store",
								"url": "http://swagger.io"
							}
						},
						{
							"name": "user",
							"description": "Operations about user"
						}
					],
					"schemes": [],
					"security": [],
					"paths": {
						"/pet": {
							"put": {
								"tags": [
									"pet",
									"customrathna"
								],
								"summary": "Update an existing pet",
								"description": "Update an existing pet by Id",
								"operationId": "updatePet",
								"consumes": [
									"application/xml",
									"application/json",
									"application/x-www-form-urlencoded"
								],
								"produces": [
									"application/xml",
									"application/json"
								],
								"parameters": [],
								"responses": {
									"200": {
										"description": "Successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"$ref": "#/components/schemas/Pet"
												},
												"examples": {}
											},
											"application/xml": {
												"schema": {
													"$ref": "#/components/schemas/Pet"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"400": {
										"description": "Invalid ID supplied",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"404": {
										"description": "Pet not found",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"405": {
										"description": "Validation exception",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"security": [
									{
										"requirements": {
											"petstore_auth": [
												"write:pets",
												"read:pets"
											]
										}
									}
								],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {
										"application/json": {
											"schema": {
												"$ref": "#/components/schemas/Pet"
											},
											"examples": {}
										},
										"application/x-www-form-urlencoded": {
											"schema": {
												"$ref": "#/components/schemas/Pet"
											},
											"examples": {}
										},
										"application/xml": {
											"schema": {
												"$ref": "#/components/schemas/Pet"
											},
											"examples": {}
										}
									},
									"name": "updatePet"
								}
							},
							"post": {
								"tags": [
									"pet",
									"custom2"
								],
								"summary": "Add a new pet to the store",
								"description": "Add a new pet to the store",
								"operationId": "addPet",
								"consumes": [
									"application/xml",
									"application/json",
									"application/x-www-form-urlencoded"
								],
								"produces": [
									"application/xml",
									"application/json"
								],
								"parameters": [],
								"responses": {
									"200": {
										"description": "Successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"$ref": "#/components/schemas/Pet"
												},
												"examples": {}
											},
											"application/xml": {
												"schema": {
													"$ref": "#/components/schemas/Pet"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"405": {
										"description": "Invalid input",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"security": [
									{
										"requirements": {
											"petstore_auth": [
												"write:pets",
												"read:pets"
											]
										}
									}
								],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {
										"application/json": {
											"schema": {
												"$ref": "#/components/schemas/Pet"
											},
											"examples": {}
										},
										"application/x-www-form-urlencoded": {
											"schema": {
												"$ref": "#/components/schemas/Pet"
											},
											"examples": {}
										},
										"application/xml": {
											"schema": {
												"$ref": "#/components/schemas/Pet"
											},
											"examples": {}
										}
									},
									"name": "addPet"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/pet",
							"enabled": true
						},
						"/pet/findByStatus": {
							"get": {
								"tags": [
									"pet"
								],
								"summary": "Finds Pets by status",
								"description": "Multiple status values can be provided with comma separated strings",
								"operationId": "findPetsByStatus",
								"produces": [
									"application/xml",
									"application/json"
								],
								"parameters": [
									{
										"default": "available",
										"description": "Status values that need to be considered for filter",
										"explode": true,
										"in": "query",
										"name": "status",
										"parameterSchema": {
											"default": "available",
											"enum": [
												"available",
												"pending",
												"sold"
											],
											"type": "string"
										},
										"required": false,
										"style": "FORM",
										"type": "string"
									}
								],
								"responses": {
									"200": {
										"description": "successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"type": "gateway",
													"schema": "{\"type\":\"array\",\"items\":{\"type\":\"object\",\"$ref\":\"#/components/schemas/Pet\"}}"
												},
												"examples": {}
											},
											"application/xml": {
												"schema": {
													"type": "gateway",
													"schema": "{\"type\":\"array\",\"items\":{\"type\":\"object\",\"$ref\":\"#/components/schemas/Pet\"}}"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"400": {
										"description": "Invalid status value",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"security": [
									{
										"requirements": {
											"petstore_auth": [
												"write:pets",
												"read:pets"
											]
										}
									}
								],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "findPetsByStatus"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/pet/findByStatus",
							"enabled": true
						},
						"/pet/findByTags": {
							"get": {
								"tags": [
									"pet"
								],
								"summary": "Finds Pets by tags",
								"description": "Multiple tags can be provided with comma separated strings. Use tag1, tag2, tag3 for testing.",
								"operationId": "findPetsByTags",
								"produces": [
									"application/xml",
									"application/json"
								],
								"parameters": [
									{
										"description": "Tags to filter by",
										"explode": true,
										"in": "query",
										"name": "tags",
										"parameterSchema": {
											"items": {
												"type": "string"
											},
											"type": "array"
										},
										"required": false,
										"style": "FORM"
									}
								],
								"responses": {
									"200": {
										"description": "successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"type": "gateway",
													"schema": "{\"type\":\"array\",\"items\":{\"type\":\"object\",\"$ref\":\"#/components/schemas/Pet\"}}"
												},
												"examples": {}
											},
											"application/xml": {
												"schema": {
													"type": "gateway",
													"schema": "{\"type\":\"array\",\"items\":{\"type\":\"object\",\"$ref\":\"#/components/schemas/Pet\"}}"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"400": {
										"description": "Invalid tag value",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"security": [
									{
										"requirements": {
											"petstore_auth": [
												"write:pets",
												"read:pets"
											]
										}
									}
								],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "findPetsByTags"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/pet/findByTags",
							"enabled": true
						},
						"/pet/{petId}": {
							"get": {
								"tags": [
									"pet"
								],
								"summary": "Find pet by ID",
								"description": "Returns a single pet",
								"operationId": "getPetById",
								"produces": [
									"application/xml",
									"application/json"
								],
								"parameters": [
									{
										"description": "ID of pet to return",
										"explode": false,
										"format": "int64",
										"in": "path",
										"name": "petId",
										"parameterSchema": {
											"format": "int64",
											"type": "integer"
										},
										"required": true,
										"style": "SIMPLE",
										"type": "integer"
									}
								],
								"responses": {
									"200": {
										"description": "successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"$ref": "#/components/schemas/Pet"
												},
												"examples": {}
											},
											"application/xml": {
												"schema": {
													"$ref": "#/components/schemas/Pet"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"400": {
										"description": "Invalid ID supplied",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"404": {
										"description": "Pet not found",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"security": [
									{
										"requirements": {
											"api_key": []
										}
									},
									{
										"requirements": {
											"petstore_auth": [
												"write:pets",
												"read:pets"
											]
										}
									}
								],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "getPetById"
								}
							},
							"post": {
								"tags": [
									"pet"
								],
								"summary": "Updates a pet in the store with form data",
								"description": "",
								"operationId": "updatePetWithForm",
								"parameters": [
									{
										"description": "ID of pet that needs to be updated",
										"explode": false,
										"format": "int64",
										"in": "path",
										"name": "petId",
										"parameterSchema": {
											"format": "int64",
											"type": "integer"
										},
										"required": true,
										"style": "SIMPLE",
										"type": "integer"
									},
									{
										"description": "Name of pet that needs to be updated",
										"explode": true,
										"in": "query",
										"name": "name",
										"parameterSchema": {
											"type": "string"
										},
										"required": false,
										"style": "FORM",
										"type": "string"
									},
									{
										"description": "Status of pet that needs to be updated",
										"explode": true,
										"in": "query",
										"name": "status",
										"parameterSchema": {
											"type": "string"
										},
										"required": false,
										"style": "FORM",
										"type": "string"
									}
								],
								"responses": {
									"405": {
										"description": "Invalid input",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"security": [
									{
										"requirements": {
											"petstore_auth": [
												"write:pets",
												"read:pets"
											]
										}
									}
								],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "updatePetWithForm"
								}
							},
							"delete": {
								"tags": [
									"pet"
								],
								"summary": "Deletes a pet",
								"description": "",
								"operationId": "deletePet",
								"parameters": [
									{
										"description": "",
										"explode": false,
										"in": "header",
										"name": "api_key",
										"parameterSchema": {
											"type": "string"
										},
										"required": false,
										"style": "SIMPLE",
										"type": "string"
									},
									{
										"description": "Pet id to delete",
										"explode": false,
										"format": "int64",
										"in": "path",
										"name": "petId",
										"parameterSchema": {
											"format": "int64",
											"type": "integer"
										},
										"required": true,
										"style": "SIMPLE",
										"type": "integer"
									}
								],
								"responses": {
									"400": {
										"description": "Invalid pet value",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"security": [
									{
										"requirements": {
											"petstore_auth": [
												"write:pets",
												"read:pets"
											]
										}
									}
								],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "deletePet"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/pet/{petId}",
							"enabled": true
						},
						"/pet/{petId}/uploadImage": {
							"post": {
								"tags": [
									"pet"
								],
								"summary": "uploads an image",
								"description": "",
								"operationId": "uploadFile",
								"consumes": [
									"application/octet-stream"
								],
								"produces": [
									"application/json"
								],
								"parameters": [
									{
										"description": "ID of pet to update",
										"explode": false,
										"format": "int64",
										"in": "path",
										"name": "petId",
										"parameterSchema": {
											"format": "int64",
											"type": "integer"
										},
										"required": true,
										"style": "SIMPLE",
										"type": "integer"
									},
									{
										"description": "Additional Metadata",
										"explode": true,
										"in": "query",
										"name": "additionalMetadata",
										"parameterSchema": {
											"type": "string"
										},
										"required": false,
										"style": "FORM",
										"type": "string"
									}
								],
								"responses": {
									"200": {
										"description": "successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"$ref": "#/components/schemas/ApiResponse"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"security": [
									{
										"requirements": {
											"petstore_auth": [
												"write:pets",
												"read:pets"
											]
										}
									}
								],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {
										"application/octet-stream": {
											"schema": {
												"type": "gateway",
												"schema": "{\"type\":\"string\",\"format\":\"binary\"}"
											},
											"examples": {}
										}
									},
									"name": "uploadFile"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/pet/{petId}/uploadImage",
							"enabled": true
						},
						"/store/inventory": {
							"get": {
								"tags": [
									"store"
								],
								"summary": "Returns pet inventories by status",
								"description": "Returns a map of status codes to quantities",
								"operationId": "getInventory",
								"produces": [
									"application/json"
								],
								"parameters": [],
								"responses": {
									"200": {
										"description": "successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"type": "gateway",
													"schema": "{\"type\":\"object\",\"additionalProperties\":{\"type\":\"integer\",\"format\":\"int32\"}}"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"security": [
									{
										"requirements": {
											"api_key": []
										}
									}
								],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "getInventory"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/store/inventory",
							"enabled": true
						},
						"/store/order": {
							"post": {
								"tags": [
									"store"
								],
								"summary": "Place an order for a pet",
								"description": "Place a new order in the store",
								"operationId": "placeOrder",
								"consumes": [
									"application/xml",
									"application/json",
									"application/x-www-form-urlencoded"
								],
								"produces": [
									"application/json"
								],
								"parameters": [],
								"responses": {
									"200": {
										"description": "successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"$ref": "#/components/schemas/Order"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"405": {
										"description": "Invalid input",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {
										"application/json": {
											"schema": {
												"$ref": "#/components/schemas/Order"
											},
											"examples": {}
										},
										"application/x-www-form-urlencoded": {
											"schema": {
												"$ref": "#/components/schemas/Order"
											},
											"examples": {}
										},
										"application/xml": {
											"schema": {
												"$ref": "#/components/schemas/Order"
											},
											"examples": {}
										}
									},
									"name": "placeOrder"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/store/order",
							"enabled": true
						},
						"/store/order/{orderId}": {
							"get": {
								"tags": [
									"store"
								],
								"summary": "Find purchase order by ID",
								"description": "For valid response try integer IDs with value <= 5 or > 10. Other values will generate exceptions.",
								"operationId": "getOrderById",
								"produces": [
									"application/xml",
									"application/json"
								],
								"parameters": [
									{
										"description": "ID of order that needs to be fetched",
										"explode": false,
										"format": "int64",
										"in": "path",
										"name": "orderId",
										"parameterSchema": {
											"format": "int64",
											"type": "integer"
										},
										"required": true,
										"style": "SIMPLE",
										"type": "integer"
									}
								],
								"responses": {
									"200": {
										"description": "successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"$ref": "#/components/schemas/Order"
												},
												"examples": {}
											},
											"application/xml": {
												"schema": {
													"$ref": "#/components/schemas/Order"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"400": {
										"description": "Invalid ID supplied",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"404": {
										"description": "Order not found",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "getOrderById"
								}
							},
							"delete": {
								"tags": [
									"store"
								],
								"summary": "Delete purchase order by ID",
								"description": "For valid response try integer IDs with value < 1000. Anything above 1000 or nonintegers will generate API errors",
								"operationId": "deleteOrder",
								"parameters": [
									{
										"description": "ID of the order that needs to be deleted",
										"explode": false,
										"format": "int64",
										"in": "path",
										"name": "orderId",
										"parameterSchema": {
											"format": "int64",
											"type": "integer"
										},
										"required": true,
										"style": "SIMPLE",
										"type": "integer"
									}
								],
								"responses": {
									"400": {
										"description": "Invalid ID supplied",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"404": {
										"description": "Order not found",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "deleteOrder"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/store/order/{orderId}",
							"enabled": true
						},
						"/user": {
							"post": {
								"tags": [
									"user"
								],
								"summary": "Create user",
								"description": "This can only be done by the logged in user.",
								"operationId": "createUser",
								"consumes": [
									"application/xml",
									"application/json",
									"application/x-www-form-urlencoded"
								],
								"produces": [
									"application/xml",
									"application/json"
								],
								"parameters": [],
								"responses": {
									"default": {
										"description": "successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"$ref": "#/components/schemas/User"
												},
												"examples": {}
											},
											"application/xml": {
												"schema": {
													"$ref": "#/components/schemas/User"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {
										"application/json": {
											"schema": {
												"$ref": "#/components/schemas/User"
											},
											"examples": {}
										},
										"application/x-www-form-urlencoded": {
											"schema": {
												"$ref": "#/components/schemas/User"
											},
											"examples": {}
										},
										"application/xml": {
											"schema": {
												"$ref": "#/components/schemas/User"
											},
											"examples": {}
										}
									},
									"name": "createUser"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/user",
							"enabled": true
						},
						"/user/createWithList": {
							"post": {
								"tags": [
									"user"
								],
								"summary": "Creates list of users with given input array",
								"description": "Creates list of users with given input array",
								"operationId": "createUsersWithListInput",
								"consumes": [
									"application/json"
								],
								"produces": [
									"application/xml",
									"application/json"
								],
								"parameters": [],
								"responses": {
									"200": {
										"description": "Successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"$ref": "#/components/schemas/User"
												},
												"examples": {}
											},
											"application/xml": {
												"schema": {
													"$ref": "#/components/schemas/User"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"default": {
										"description": "successful operation",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {
										"application/json": {
											"schema": {
												"type": "gateway",
												"schema": "{\"type\":\"array\",\"items\":{\"type\":\"object\",\"$ref\":\"#/components/schemas/User\"}}"
											},
											"examples": {}
										}
									},
									"name": "createUsersWithListInput"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/user/createWithList",
							"enabled": true
						},
						"/user/login": {
							"get": {
								"tags": [
									"user"
								],
								"summary": "Logs user into the system",
								"description": "",
								"operationId": "loginUser",
								"produces": [
									"application/xml",
									"application/json"
								],
								"parameters": [
									{
										"description": "The user name for login",
										"explode": true,
										"in": "query",
										"name": "username",
										"parameterSchema": {
											"type": "string"
										},
										"required": false,
										"style": "FORM",
										"type": "string"
									},
									{
										"description": "The password for login in clear text",
										"explode": true,
										"in": "query",
										"name": "password",
										"parameterSchema": {
											"type": "string"
										},
										"required": false,
										"style": "FORM",
										"type": "string"
									}
								],
								"responses": {
									"200": {
										"description": "successful operation",
										"headersV3": {
											"X-Expires-After": {
												"name": "X-Expires-After",
												"in": "header",
												"description": "date in UTC when token expires",
												"required": false,
												"type": "string",
												"format": "date-time",
												"style": "SIMPLE",
												"explode": false,
												"examples": {},
												"parameterSchema": {
													"type": "string",
													"format": "date-time"
												}
											},
											"X-Rate-Limit": {
												"name": "X-Rate-Limit",
												"in": "header",
												"description": "calls per hour allowed by the user",
												"required": false,
												"type": "integer",
												"format": "int32",
												"style": "SIMPLE",
												"explode": false,
												"examples": {},
												"parameterSchema": {
													"type": "integer",
													"format": "int32"
												}
											}
										},
										"content": {
											"application/json": {
												"schema": {
													"type": "gateway",
													"schema": "{\"type\":\"string\"}"
												},
												"examples": {}
											},
											"application/xml": {
												"schema": {
													"type": "gateway",
													"schema": "{\"type\":\"string\"}"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"400": {
										"description": "Invalid username/password supplied",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "loginUser"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/user/login",
							"enabled": true
						},
						"/user/logout": {
							"get": {
								"tags": [
									"user"
								],
								"summary": "Logs out current logged in user session",
								"description": "",
								"operationId": "logoutUser",
								"parameters": [],
								"responses": {
									"default": {
										"description": "successful operation",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "logoutUser"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/user/logout",
							"enabled": true
						},
						"/user/{username}": {
							"get": {
								"tags": [
									"user"
								],
								"summary": "Get user by user name",
								"description": "",
								"operationId": "getUserByName",
								"produces": [
									"application/xml",
									"application/json"
								],
								"parameters": [
									{
										"description": "The name that needs to be fetched. Use user1 for testing. ",
										"explode": false,
										"in": "path",
										"name": "username",
										"parameterSchema": {
											"type": "string"
										},
										"required": true,
										"style": "SIMPLE",
										"type": "string"
									}
								],
								"responses": {
									"200": {
										"description": "successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"$ref": "#/components/schemas/User"
												},
												"examples": {}
											},
											"application/xml": {
												"schema": {
													"$ref": "#/components/schemas/User"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"400": {
										"description": "Invalid username supplied",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"404": {
										"description": "User not found",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "getUserByName"
								}
							},
							"put": {
								"tags": [
									"user"
								],
								"summary": "Update user",
								"description": "This can only be done by the logged in user.",
								"operationId": "updateUser",
								"consumes": [
									"application/xml",
									"application/json",
									"application/x-www-form-urlencoded"
								],
								"parameters": [
									{
										"description": "name that need to be deleted",
										"explode": false,
										"in": "path",
										"name": "username",
										"parameterSchema": {
											"type": "string"
										},
										"required": true,
										"style": "SIMPLE",
										"type": "string"
									}
								],
								"responses": {
									"default": {
										"description": "successful operation",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {
										"application/json": {
											"schema": {
												"$ref": "#/components/schemas/User"
											},
											"examples": {}
										},
										"application/x-www-form-urlencoded": {
											"schema": {
												"$ref": "#/components/schemas/User"
											},
											"examples": {}
										},
										"application/xml": {
											"schema": {
												"$ref": "#/components/schemas/User"
											},
											"examples": {}
										}
									},
									"name": "updateUser"
								}
							},
							"delete": {
								"tags": [
									"user"
								],
								"summary": "Delete user",
								"description": "This can only be done by the logged in user.",
								"operationId": "deleteUser",
								"parameters": [
									{
										"description": "The name that needs to be deleted",
										"explode": false,
										"in": "path",
										"name": "username",
										"parameterSchema": {
											"type": "string"
										},
										"required": true,
										"style": "SIMPLE",
										"type": "string"
									}
								],
								"responses": {
									"400": {
										"description": "Invalid username supplied",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"404": {
										"description": "User not found",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "deleteUser"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/user/{username}",
							"enabled": true
						}
					},
					"securityDefinitions": {},
					"definitions": {},
					"parameters": {},
					"baseUriParameters": [],
					"externalDocs": [
						{
							"description": "Find out more about Swagger",
							"url": "http://swagger.io"
						}
					],
					"servers": [
						{
							"url": "/api/v3",
							"variables": {}
						}
					],
					"components": {
						"schemas": {
							"Address": {
								"type": "gateway",
								"schema": "{\"type\":\"object\",\"properties\":{\"city\":{\"type\":\"string\",\"example\":\"Palo Alto\"},\"state\":{\"type\":\"string\",\"example\":\"CA\"},\"street\":{\"type\":\"string\",\"example\":\"437 Lytton\"},\"zip\":{\"type\":\"string\",\"example\":\"94301\"}},\"xml\":{\"name\":\"address\"}}"
							},
							"ApiResponse": {
								"type": "gateway",
								"schema": "{\"type\":\"object\",\"properties\":{\"code\":{\"type\":\"integer\",\"format\":\"int32\"},\"message\":{\"type\":\"string\"},\"type\":{\"type\":\"string\"}},\"xml\":{\"name\":\"##default\"}}"
							},
							"Category": {
								"type": "gateway",
								"schema": "{\"type\":\"object\",\"properties\":{\"id\":{\"type\":\"integer\",\"format\":\"int64\",\"example\":\"1\"},\"name\":{\"type\":\"string\",\"example\":\"Dogs\"}},\"xml\":{\"name\":\"category\"}}"
							},
							"Customer": {
								"type": "gateway",
								"schema": "{\"type\":\"object\",\"properties\":{\"address\":{\"type\":\"array\",\"xml\":{\"name\":\"addresses\",\"wrapped\":true},\"items\":{\"type\":\"object\",\"$ref\":\"#/components/schemas/Address\"}},\"id\":{\"type\":\"integer\",\"format\":\"int64\",\"example\":\"100000\"},\"username\":{\"type\":\"string\",\"example\":\"fehguy\"}},\"xml\":{\"name\":\"customer\"}}"
							},
							"Order": {
								"type": "gateway",
								"schema": "{\"type\":\"object\",\"properties\":{\"complete\":{\"type\":\"boolean\"},\"id\":{\"type\":\"integer\",\"format\":\"int64\",\"example\":\"10\"},\"petId\":{\"type\":\"integer\",\"format\":\"int64\",\"example\":\"198772\"},\"quantity\":{\"type\":\"integer\",\"format\":\"int32\",\"example\":\"7\"},\"shipDate\":{\"type\":\"string\",\"format\":\"date-time\"},\"status\":{\"type\":\"string\",\"description\":\"Order Status\",\"example\":\"approved\",\"enum\":[\"placed\",\"approved\",\"delivered\"]}},\"xml\":{\"name\":\"order\"}}"
							},
							"Pet": {
								"type": "gateway",
								"schema": "{\"required\":[\"name\",\"photoUrls\"],\"type\":\"object\",\"properties\":{\"category\":{\"type\":\"object\",\"$ref\":\"#/components/schemas/Category\"},\"id\":{\"type\":\"integer\",\"format\":\"int64\",\"example\":\"10\"},\"name\":{\"type\":\"string\",\"example\":\"doggie\"},\"photoUrls\":{\"type\":\"array\",\"xml\":{\"wrapped\":true},\"items\":{\"type\":\"string\",\"xml\":{\"name\":\"photoUrl\"}}},\"status\":{\"type\":\"string\",\"description\":\"pet status in the store\",\"enum\":[\"available\",\"pending\",\"sold\"]},\"tags\":{\"type\":\"array\",\"xml\":{\"wrapped\":true},\"items\":{\"type\":\"object\",\"$ref\":\"#/components/schemas/Tag\"}}},\"xml\":{\"name\":\"pet\"}}"
							},
							"Tag": {
								"type": "gateway",
								"schema": "{\"type\":\"object\",\"properties\":{\"id\":{\"type\":\"integer\",\"format\":\"int64\"},\"name\":{\"type\":\"string\"}},\"xml\":{\"name\":\"tag\"}}"
							},
							"User": {
								"type": "gateway",
								"schema": "{\"type\":\"object\",\"properties\":{\"email\":{\"type\":\"string\",\"example\":\"john@email.com\"},\"firstName\":{\"type\":\"string\",\"example\":\"John\"},\"id\":{\"type\":\"integer\",\"format\":\"int64\",\"example\":\"10\"},\"lastName\":{\"type\":\"string\",\"example\":\"James\"},\"password\":{\"type\":\"string\",\"example\":\"12345\"},\"phone\":{\"type\":\"string\",\"example\":\"12345\"},\"userStatus\":{\"type\":\"integer\",\"description\":\"User Status\",\"format\":\"int32\",\"example\":\"1\"},\"username\":{\"type\":\"string\",\"example\":\"theUser\"}},\"xml\":{\"name\":\"user\"}}"
							}
						},
						"responses": {},
						"parameters": {},
						"examples": {},
						"requestBodies": {
							"Pet": {
								"content": {
									"application/json": {
										"schema": {
											"$ref": "#/components/schemas/Pet"
										},
										"examples": {}
									},
									"application/xml": {
										"schema": {
											"$ref": "#/components/schemas/Pet"
										},
										"examples": {}
									}
								},
								"name": "Pet"
							},
							"UserArray": {
								"content": {
									"application/json": {
										"schema": {
											"type": "gateway",
											"schema": "{\"type\":\"array\",\"items\":{\"type\":\"object\",\"$ref\":\"#/components/schemas/User\"}}"
										},
										"examples": {}
									}
								},
								"name": "UserArray"
							}
						},
						"headers": {},
						"links": {},
						"callbacks": {}
					},
					"type": "rest"
				},
				"nativeEndpoint": [
					{
						"passSecurityHeaders": true,
						"uri": "/api/v3",
						"connectionTimeoutDuration": 0,
						"alias": false
					}
				],
				"apiName": "petstore",
				"apiVersion": "1.0.17",
				"apiDescription": "This is a sample Pet Store Server based on the OpenAPI 3.0 specification.  You can find out more about\nSwagger at [http://swagger.io](http://swagger.io). In the third iteration of the pet store, we've switched to the design first approach!\nYou can now help us improve the API whether it's by making changes to the definition itself or to the code.\nThat way, with time, we can improve the API in general, and expose some of the new features in OAS3.\n\nSome useful links:\n- [The Pet Store repository](https://github.com/swagger-api/swagger-petstore)\n- [The source API definition for the Pet Store](https://github.com/swagger-api/swagger-petstore/blob/master/src/main/resources/openapi.yaml)",
				"maturityState": "Beta",
				"apiGroups": [
					"Finance Banking and Insurance",
					"Sales and Ordering"
				],
				"isActive": true,
				"type": "REST",
				"owner": "wecare@apiwheel.dev",
				"policies": [
					"4631db34-04a7-4cad-860c-6a182c019d4a"
				],
				"tracingEnabled": false,
				"scopes": [],
				"publishedPortals": [
					"69ca5c6b-8b0a-4bae-8021-4ccf0ffd136f"
				],
				"creationDate": "2023-03-18 06:29:48 GMT",
				"lastModified": "2023-03-29 13:43:56 GMT",
				"systemVersion": 1,
				"gatewayEndpoints": {},
				"deployments": [
					"APIGateway"
				],
				"microgatewayEndpoints": [],
				"appMeshEndpoints": [],
				"id": "2b598e47-3e0c-4b0f-8a72-da7fdf6c5ea2"
			},
			"responseStatus": "SUCCESS",
			"gatewayEndPoints": [
				"http://env688761.apigw-aw-us.webmethods.io/gateway/petstore/1.0.17"
			],
			"gatewayEndPointList": [
				{
					"endpointName": "DEFAULT_GATEWAY_ENDPOINT",
					"endpointDisplayName": "Default",
					"endpoint": "gateway/petstore/1.0.17",
					"endpointType": "DEFAULT",
					"endpointUrls": [
						"http://env688761.apigw-aw-us.webmethods.io/gateway/petstore/1.0.17"
					]
				}
			],
			"versions": [
				{
					"versionNumber": "1.0.17",
					"apiId": "2b598e47-3e0c-4b0f-8a72-da7fdf6c5ea2"
				}
			]
		}
	}`
	mc := &MockClient{}
	webMethodsClient, _ := NewClient(cfg, mc)
	mc.SendFunc = func(request coreapi.Request) (*coreapi.Response, error) {
		return &coreapi.Response{
			Code: 200,
			Body: []byte(response),
		}, nil
	}
	apiResponse, err := webMethodsClient.GetApiDetails("2b598e47-3e0c-4b0f-8a72-da7fdf6c5ea2")
	assert.Nil(t, err)
	assert.Equal(t, apiResponse.Api.MaturityState, "Beta")
	assert.Equal(t, apiResponse.GatewayEndPoints[0], "http://env688761.apigw-aw-us.webmethods.io/gateway/petstore/1.0.17")
	assert.Equal(t, apiResponse.Api.Owner, "wecare@apiwheel.dev")
	assert.Equal(t, apiResponse.Api.ApiVersion, "1.0.17")
	assert.Equal(t, apiResponse.Api.ApiDefinition.Info.Title, "Swagger Petstore - OpenAPI 3.0")
	assert.NotNil(t, apiResponse.Api.ApiDescription)
	assert.NotNil(t, apiResponse.Api.ApiGroups[0], "Finance Banking and Insurance")
	assert.NotNil(t, apiResponse.Api.ApiGroups[1], "Sales and Ordering")

}

func TestGetApiSpec(t *testing.T) {

	response := `{
		"openapi": "3.0.1",
		"info": {
			"title": "petstore",
			"description": "This is a sample Pet Store Server based on the OpenAPI 3.0 specification.  You can find out more about\nSwagger at [http://swagger.io](http://swagger.io). In the third iteration of the pet store, we've switched to the design first approach!\nYou can now help us improve the API whether it's by making changes to the definition itself or to the code.\nThat way, with time, we can improve the API in general, and expose some of the new features in OAS3.\n\nSome useful links:\n- [The Pet Store repository](https://github.com/swagger-api/swagger-petstore)\n- [The source API definition for the Pet Store](https://github.com/swagger-api/swagger-petstore/blob/master/src/main/resources/openapi.yaml)",
			"termsOfService": "http://swagger.io/terms/",
			"contact": {
				"email": "apiteam@swagger.io"
			},
			"license": {
				"name": "Apache 2.0",
				"url": "http://www.apache.org/licenses/LICENSE-2.0.html"
			},
			"version": "1.0.17"
		},
		"externalDocs": {
			"description": "Find out more about Swagger",
			"url": "http://swagger.io"
		},
		"servers": [
			{
				"url": "http://env688761.apigw-aw-us.webmethods.io/gateway/petstore/1.0.17"
			}
		],
		"security": [
			{
				"apiKey": []
			}
		],
		"tags": [
			{
				"name": "customrathna"
			},
			{
				"name": "custom2"
			},
			{
				"name": "pet",
				"description": "Everything about your Pets",
				"externalDocs": {
					"description": "Find out more",
					"url": "http://swagger.io"
				}
			},
			{
				"name": "store",
				"description": "Access to Petstore orders",
				"externalDocs": {
					"description": "Find out more about our store",
					"url": "http://swagger.io"
				}
			},
			{
				"name": "user",
				"description": "Operations about user"
			}
		],
		"paths": {
			"/pet": {
				"summary": "/pet",
				"put": {
					"tags": [
						"pet",
						"customrathna"
					],
					"summary": "Update an existing pet",
					"description": "Update an existing pet by Id",
					"operationId": "updatePet",
					"parameters": [],
					"requestBody": {
						"content": {
							"application/json": {
								"schema": {
									"$ref": "#/components/schemas/Pet"
								},
								"examples": {},
								"example": null
							},
							"application/x-www-form-urlencoded": {
								"schema": {
									"$ref": "#/components/schemas/Pet"
								},
								"examples": {},
								"example": null
							},
							"application/xml": {
								"schema": {
									"$ref": "#/components/schemas/Pet"
								},
								"examples": {},
								"example": null
							}
						}
					},
					"responses": {
						"200": {
							"description": "Successful operation",
							"headers": {},
							"content": {
								"application/json": {
									"schema": {
										"$ref": "#/components/schemas/Pet"
									},
									"examples": {},
									"example": null
								},
								"application/xml": {
									"schema": {
										"$ref": "#/components/schemas/Pet"
									},
									"examples": {},
									"example": null
								}
							},
							"links": {}
						},
						"400": {
							"description": "Invalid ID supplied",
							"headers": {},
							"content": {},
							"links": {}
						},
						"404": {
							"description": "Pet not found",
							"headers": {},
							"content": {},
							"links": {}
						},
						"405": {
							"description": "Validation exception",
							"headers": {},
							"content": {},
							"links": {}
						}
					}
				},
				"post": {
					"tags": [
						"pet",
						"custom2"
					],
					"summary": "Add a new pet to the store",
					"description": "Add a new pet to the store",
					"operationId": "addPet",
					"parameters": [],
					"requestBody": {
						"content": {
							"application/json": {
								"schema": {
									"$ref": "#/components/schemas/Pet"
								},
								"examples": {},
								"example": null
							},
							"application/x-www-form-urlencoded": {
								"schema": {
									"$ref": "#/components/schemas/Pet"
								},
								"examples": {},
								"example": null
							},
							"application/xml": {
								"schema": {
									"$ref": "#/components/schemas/Pet"
								},
								"examples": {},
								"example": null
							}
						}
					},
					"responses": {
						"200": {
							"description": "Successful operation",
							"headers": {},
							"content": {
								"application/json": {
									"schema": {
										"$ref": "#/components/schemas/Pet"
									},
									"examples": {},
									"example": null
								},
								"application/xml": {
									"schema": {
										"$ref": "#/components/schemas/Pet"
									},
									"examples": {},
									"example": null
								}
							},
							"links": {}
						},
						"405": {
							"description": "Invalid input",
							"headers": {},
							"content": {},
							"links": {}
						}
					}
				},
				"parameters": []
			},
			"/user/{username}": {
				"summary": "/user/{username}",
				"get": {
					"tags": [
						"user"
					],
					"summary": "Get user by user name",
					"description": "",
					"operationId": "getUserByName",
					"parameters": [
						{
							"name": "username",
							"in": "path",
							"description": "The name that needs to be fetched. Use user1 for testing. ",
							"required": true,
							"allowEmptyValue": false,
							"style": "simple",
							"explode": false,
							"schema": {
								"type": "string",
								"example": null
							}
						}
					],
					"responses": {
						"200": {
							"description": "successful operation",
							"headers": {},
							"content": {
								"application/json": {
									"schema": {
										"$ref": "#/components/schemas/User"
									},
									"examples": {},
									"example": null
								},
								"application/xml": {
									"schema": {
										"$ref": "#/components/schemas/User"
									},
									"examples": {},
									"example": null
								}
							},
							"links": {}
						},
						"400": {
							"description": "Invalid username supplied",
							"headers": {},
							"content": {},
							"links": {}
						},
						"404": {
							"description": "User not found",
							"headers": {},
							"content": {},
							"links": {}
						}
					}
				},
				"put": {
					"tags": [
						"user"
					],
					"summary": "Update user",
					"description": "This can only be done by the logged in user.",
					"operationId": "updateUser",
					"parameters": [
						{
							"name": "username",
							"in": "path",
							"description": "name that need to be deleted",
							"required": true,
							"allowEmptyValue": false,
							"style": "simple",
							"explode": false,
							"schema": {
								"type": "string",
								"example": null
							}
						}
					],
					"requestBody": {
						"content": {
							"application/json": {
								"schema": {
									"$ref": "#/components/schemas/User"
								},
								"examples": {},
								"example": null
							},
							"application/x-www-form-urlencoded": {
								"schema": {
									"$ref": "#/components/schemas/User"
								},
								"examples": {},
								"example": null
							},
							"application/xml": {
								"schema": {
									"$ref": "#/components/schemas/User"
								},
								"examples": {},
								"example": null
							}
						}
					},
					"responses": {
						"default": {
							"description": "successful operation",
							"headers": {},
							"content": {},
							"links": {}
						}
					}
				},
				"delete": {
					"tags": [
						"user"
					],
					"summary": "Delete user",
					"description": "This can only be done by the logged in user.",
					"operationId": "deleteUser",
					"parameters": [
						{
							"name": "username",
							"in": "path",
							"description": "The name that needs to be deleted",
							"required": true,
							"allowEmptyValue": false,
							"style": "simple",
							"explode": false,
							"schema": {
								"type": "string",
								"example": null
							}
						}
					],
					"responses": {
						"400": {
							"description": "Invalid username supplied",
							"headers": {},
							"content": {},
							"links": {}
						},
						"404": {
							"description": "User not found",
							"headers": {},
							"content": {},
							"links": {}
						}
					}
				},
				"parameters": []
			},
			"/pet/findByStatus": {
				"summary": "/pet/findByStatus",
				"get": {
					"tags": [
						"pet"
					],
					"summary": "Finds Pets by status",
					"description": "Multiple status values can be provided with comma separated strings",
					"operationId": "findPetsByStatus",
					"parameters": [
						{
							"name": "status",
							"in": "query",
							"description": "Status values that need to be considered for filter",
							"required": false,
							"allowEmptyValue": false,
							"style": "form",
							"explode": true,
							"schema": {
								"type": "string",
								"enum": [
									"available",
									"pending",
									"sold"
								],
								"default": "available",
								"example": null
							}
						}
					],
					"responses": {
						"200": {
							"description": "successful operation",
							"headers": {},
							"content": {
								"application/json": {
									"schema": {
										"type": "array",
										"items": {
											"$ref": "#/components/schemas/Pet"
										},
										"example": null
									},
									"examples": {},
									"example": null
								},
								"application/xml": {
									"schema": {
										"type": "array",
										"items": {
											"$ref": "#/components/schemas/Pet"
										},
										"example": null
									},
									"examples": {},
									"example": null
								}
							},
							"links": {}
						},
						"400": {
							"description": "Invalid status value",
							"headers": {},
							"content": {},
							"links": {}
						}
					}
				},
				"parameters": []
			},
			"/user/createWithList": {
				"summary": "/user/createWithList",
				"post": {
					"tags": [
						"user"
					],
					"summary": "Creates list of users with given input array",
					"description": "Creates list of users with given input array",
					"operationId": "createUsersWithListInput",
					"parameters": [],
					"requestBody": {
						"content": {
							"application/json": {
								"schema": {
									"type": "array",
									"items": {
										"$ref": "#/components/schemas/User"
									},
									"example": null
								},
								"examples": {},
								"example": null
							}
						}
					},
					"responses": {
						"200": {
							"description": "Successful operation",
							"headers": {},
							"content": {
								"application/json": {
									"schema": {
										"$ref": "#/components/schemas/User"
									},
									"examples": {},
									"example": null
								},
								"application/xml": {
									"schema": {
										"$ref": "#/components/schemas/User"
									},
									"examples": {},
									"example": null
								}
							},
							"links": {}
						},
						"default": {
							"description": "successful operation",
							"headers": {},
							"content": {},
							"links": {}
						}
					}
				},
				"parameters": []
			},
			"/pet/{petId}/uploadImage": {
				"summary": "/pet/{petId}/uploadImage",
				"post": {
					"tags": [
						"pet"
					],
					"summary": "uploads an image",
					"description": "",
					"operationId": "uploadFile",
					"parameters": [
						{
							"name": "petId",
							"in": "path",
							"description": "ID of pet to update",
							"required": true,
							"allowEmptyValue": false,
							"style": "simple",
							"explode": false,
							"schema": {
								"type": "integer",
								"format": "int64",
								"example": null
							}
						},
						{
							"name": "additionalMetadata",
							"in": "query",
							"description": "Additional Metadata",
							"required": false,
							"allowEmptyValue": false,
							"style": "form",
							"explode": true,
							"schema": {
								"type": "string",
								"example": null
							}
						}
					],
					"requestBody": {
						"content": {
							"application/octet-stream": {
								"schema": {
									"type": "string",
									"format": "binary",
									"example": null
								},
								"examples": {},
								"example": null
							}
						}
					},
					"responses": {
						"200": {
							"description": "successful operation",
							"headers": {},
							"content": {
								"application/json": {
									"schema": {
										"$ref": "#/components/schemas/ApiResponse"
									},
									"examples": {},
									"example": null
								}
							},
							"links": {}
						}
					}
				},
				"parameters": []
			},
			"/store/inventory": {
				"summary": "/store/inventory",
				"get": {
					"tags": [
						"store"
					],
					"summary": "Returns pet inventories by status",
					"description": "Returns a map of status codes to quantities",
					"operationId": "getInventory",
					"parameters": [],
					"responses": {
						"200": {
							"description": "successful operation",
							"headers": {},
							"content": {
								"application/json": {
									"schema": {
										"type": "object",
										"additionalProperties": {
											"type": "integer",
											"format": "int32",
											"example": null
										},
										"example": null
									},
									"examples": {},
									"example": null
								}
							},
							"links": {}
						}
					}
				},
				"parameters": []
			},
			"/user/login": {
				"summary": "/user/login",
				"get": {
					"tags": [
						"user"
					],
					"summary": "Logs user into the system",
					"description": "",
					"operationId": "loginUser",
					"parameters": [
						{
							"name": "username",
							"in": "query",
							"description": "The user name for login",
							"required": false,
							"allowEmptyValue": false,
							"style": "form",
							"explode": true,
							"schema": {
								"type": "string",
								"example": null
							}
						},
						{
							"name": "password",
							"in": "query",
							"description": "The password for login in clear text",
							"required": false,
							"allowEmptyValue": false,
							"style": "form",
							"explode": true,
							"schema": {
								"type": "string",
								"example": null
							}
						}
					],
					"responses": {
						"200": {
							"description": "successful operation",
							"headers": {
								"X-Rate-Limit": {
									"description": "calls per hour allowed by the user",
									"required": false,
									"style": "simple",
									"explode": false,
									"schema": {
										"type": "integer",
										"format": "int32",
										"example": null
									},
									"examples": {}
								},
								"X-Expires-After": {
									"description": "date in UTC when token expires",
									"required": false,
									"style": "simple",
									"explode": false,
									"schema": {
										"type": "string",
										"format": "date-time",
										"example": null
									},
									"examples": {}
								}
							},
							"content": {
								"application/json": {
									"schema": {
										"type": "string",
										"example": null
									},
									"examples": {},
									"example": null
								},
								"application/xml": {
									"schema": {
										"type": "string",
										"example": null
									},
									"examples": {},
									"example": null
								}
							},
							"links": {}
						},
						"400": {
							"description": "Invalid username/password supplied",
							"headers": {},
							"content": {},
							"links": {}
						}
					}
				},
				"parameters": []
			},
			"/user": {
				"summary": "/user",
				"post": {
					"tags": [
						"user"
					],
					"summary": "Create user",
					"description": "This can only be done by the logged in user.",
					"operationId": "createUser",
					"parameters": [],
					"requestBody": {
						"content": {
							"application/json": {
								"schema": {
									"$ref": "#/components/schemas/User"
								},
								"examples": {},
								"example": null
							},
							"application/x-www-form-urlencoded": {
								"schema": {
									"$ref": "#/components/schemas/User"
								},
								"examples": {},
								"example": null
							},
							"application/xml": {
								"schema": {
									"$ref": "#/components/schemas/User"
								},
								"examples": {},
								"example": null
							}
						}
					},
					"responses": {
						"default": {
							"description": "successful operation",
							"headers": {},
							"content": {
								"application/json": {
									"schema": {
										"$ref": "#/components/schemas/User"
									},
									"examples": {},
									"example": null
								},
								"application/xml": {
									"schema": {
										"$ref": "#/components/schemas/User"
									},
									"examples": {},
									"example": null
								}
							},
							"links": {}
						}
					}
				},
				"parameters": []
			},
			"/pet/findByTags": {
				"summary": "/pet/findByTags",
				"get": {
					"tags": [
						"pet"
					],
					"summary": "Finds Pets by tags",
					"description": "Multiple tags can be provided with comma separated strings. Use tag1, tag2, tag3 for testing.",
					"operationId": "findPetsByTags",
					"parameters": [
						{
							"name": "tags",
							"in": "query",
							"description": "Tags to filter by",
							"required": false,
							"allowEmptyValue": false,
							"style": "form",
							"explode": true,
							"schema": {
								"type": "array",
								"items": {
									"type": "string",
									"example": null
								},
								"example": null
							}
						}
					],
					"responses": {
						"200": {
							"description": "successful operation",
							"headers": {},
							"content": {
								"application/json": {
									"schema": {
										"type": "array",
										"items": {
											"$ref": "#/components/schemas/Pet"
										},
										"example": null
									},
									"examples": {},
									"example": null
								},
								"application/xml": {
									"schema": {
										"type": "array",
										"items": {
											"$ref": "#/components/schemas/Pet"
										},
										"example": null
									},
									"examples": {},
									"example": null
								}
							},
							"links": {}
						},
						"400": {
							"description": "Invalid tag value",
							"headers": {},
							"content": {},
							"links": {}
						}
					}
				},
				"parameters": []
			},
			"/store/order": {
				"summary": "/store/order",
				"post": {
					"tags": [
						"store"
					],
					"summary": "Place an order for a pet",
					"description": "Place a new order in the store",
					"operationId": "placeOrder",
					"parameters": [],
					"requestBody": {
						"content": {
							"application/json": {
								"schema": {
									"$ref": "#/components/schemas/Order"
								},
								"examples": {},
								"example": null
							},
							"application/x-www-form-urlencoded": {
								"schema": {
									"$ref": "#/components/schemas/Order"
								},
								"examples": {},
								"example": null
							},
							"application/xml": {
								"schema": {
									"$ref": "#/components/schemas/Order"
								},
								"examples": {},
								"example": null
							}
						}
					},
					"responses": {
						"200": {
							"description": "successful operation",
							"headers": {},
							"content": {
								"application/json": {
									"schema": {
										"$ref": "#/components/schemas/Order"
									},
									"examples": {},
									"example": null
								}
							},
							"links": {}
						},
						"405": {
							"description": "Invalid input",
							"headers": {},
							"content": {},
							"links": {}
						}
					}
				},
				"parameters": []
			},
			"/user/logout": {
				"summary": "/user/logout",
				"get": {
					"tags": [
						"user"
					],
					"summary": "Logs out current logged in user session",
					"description": "",
					"operationId": "logoutUser",
					"parameters": [],
					"responses": {
						"default": {
							"description": "successful operation",
							"headers": {},
							"content": {},
							"links": {}
						}
					}
				},
				"parameters": []
			},
			"/pet/{petId}": {
				"summary": "/pet/{petId}",
				"get": {
					"tags": [
						"pet"
					],
					"summary": "Find pet by ID",
					"description": "Returns a single pet",
					"operationId": "getPetById",
					"parameters": [
						{
							"name": "petId",
							"in": "path",
							"description": "ID of pet to return",
							"required": true,
							"allowEmptyValue": false,
							"style": "simple",
							"explode": false,
							"schema": {
								"type": "integer",
								"format": "int64",
								"example": null
							}
						}
					],
					"responses": {
						"200": {
							"description": "successful operation",
							"headers": {},
							"content": {
								"application/json": {
									"schema": {
										"$ref": "#/components/schemas/Pet"
									},
									"examples": {},
									"example": null
								},
								"application/xml": {
									"schema": {
										"$ref": "#/components/schemas/Pet"
									},
									"examples": {},
									"example": null
								}
							},
							"links": {}
						},
						"400": {
							"description": "Invalid ID supplied",
							"headers": {},
							"content": {},
							"links": {}
						},
						"404": {
							"description": "Pet not found",
							"headers": {},
							"content": {},
							"links": {}
						}
					}
				},
				"post": {
					"tags": [
						"pet"
					],
					"summary": "Updates a pet in the store with form data",
					"description": "",
					"operationId": "updatePetWithForm",
					"parameters": [
						{
							"name": "petId",
							"in": "path",
							"description": "ID of pet that needs to be updated",
							"required": true,
							"allowEmptyValue": false,
							"style": "simple",
							"explode": false,
							"schema": {
								"type": "integer",
								"format": "int64",
								"example": null
							}
						},
						{
							"name": "name",
							"in": "query",
							"description": "Name of pet that needs to be updated",
							"required": false,
							"allowEmptyValue": false,
							"style": "form",
							"explode": true,
							"schema": {
								"type": "string",
								"example": null
							}
						},
						{
							"name": "status",
							"in": "query",
							"description": "Status of pet that needs to be updated",
							"required": false,
							"allowEmptyValue": false,
							"style": "form",
							"explode": true,
							"schema": {
								"type": "string",
								"example": null
							}
						}
					],
					"requestBody": {
						"content": {}
					},
					"responses": {
						"405": {
							"description": "Invalid input",
							"headers": {},
							"content": {},
							"links": {}
						}
					}
				},
				"delete": {
					"tags": [
						"pet"
					],
					"summary": "Deletes a pet",
					"description": "",
					"operationId": "deletePet",
					"parameters": [
						{
							"name": "api_key",
							"in": "header",
							"description": "",
							"required": false,
							"allowEmptyValue": false,
							"style": "simple",
							"explode": false,
							"schema": {
								"type": "string",
								"example": null
							}
						},
						{
							"name": "petId",
							"in": "path",
							"description": "Pet id to delete",
							"required": true,
							"allowEmptyValue": false,
							"style": "simple",
							"explode": false,
							"schema": {
								"type": "integer",
								"format": "int64",
								"example": null
							}
						}
					],
					"responses": {
						"400": {
							"description": "Invalid pet value",
							"headers": {},
							"content": {},
							"links": {}
						}
					}
				},
				"parameters": []
			},
			"/store/order/{orderId}": {
				"summary": "/store/order/{orderId}",
				"get": {
					"tags": [
						"store"
					],
					"summary": "Find purchase order by ID",
					"description": "For valid response try integer IDs with value <= 5 or > 10. Other values will generate exceptions.",
					"operationId": "getOrderById",
					"parameters": [
						{
							"name": "orderId",
							"in": "path",
							"description": "ID of order that needs to be fetched",
							"required": true,
							"allowEmptyValue": false,
							"style": "simple",
							"explode": false,
							"schema": {
								"type": "integer",
								"format": "int64",
								"example": null
							}
						}
					],
					"responses": {
						"200": {
							"description": "successful operation",
							"headers": {},
							"content": {
								"application/json": {
									"schema": {
										"$ref": "#/components/schemas/Order"
									},
									"examples": {},
									"example": null
								},
								"application/xml": {
									"schema": {
										"$ref": "#/components/schemas/Order"
									},
									"examples": {},
									"example": null
								}
							},
							"links": {}
						},
						"400": {
							"description": "Invalid ID supplied",
							"headers": {},
							"content": {},
							"links": {}
						},
						"404": {
							"description": "Order not found",
							"headers": {},
							"content": {},
							"links": {}
						}
					}
				},
				"delete": {
					"tags": [
						"store"
					],
					"summary": "Delete purchase order by ID",
					"description": "For valid response try integer IDs with value < 1000. Anything above 1000 or nonintegers will generate API errors",
					"operationId": "deleteOrder",
					"parameters": [
						{
							"name": "orderId",
							"in": "path",
							"description": "ID of the order that needs to be deleted",
							"required": true,
							"allowEmptyValue": false,
							"style": "simple",
							"explode": false,
							"schema": {
								"type": "integer",
								"format": "int64",
								"example": null
							}
						}
					],
					"responses": {
						"400": {
							"description": "Invalid ID supplied",
							"headers": {},
							"content": {},
							"links": {}
						},
						"404": {
							"description": "Order not found",
							"headers": {},
							"content": {},
							"links": {}
						}
					}
				},
				"parameters": []
			}
		},
		"components": {
			"schemas": {
				"Order": {
					"type": "object",
					"properties": {
						"petId": {
							"type": "integer",
							"format": "int64",
							"example": 198772
						},
						"quantity": {
							"type": "integer",
							"format": "int32",
							"example": 7
						},
						"id": {
							"type": "integer",
							"format": "int64",
							"example": 10
						},
						"complete": {
							"type": "boolean",
							"example": null
						},
						"shipDate": {
							"type": "string",
							"format": "date-time",
							"example": null
						},
						"status": {
							"type": "string",
							"description": "Order Status",
							"example": "approved",
							"enum": [
								"placed",
								"approved",
								"delivered"
							]
						}
					},
					"xml": {
						"name": "order"
					},
					"example": null
				},
				"Category": {
					"type": "object",
					"properties": {
						"name": {
							"type": "string",
							"example": "Dogs"
						},
						"id": {
							"type": "integer",
							"format": "int64",
							"example": 1
						}
					},
					"xml": {
						"name": "category"
					},
					"example": null
				},
				"User": {
					"type": "object",
					"properties": {
						"firstName": {
							"type": "string",
							"example": "John"
						},
						"lastName": {
							"type": "string",
							"example": "James"
						},
						"password": {
							"type": "string",
							"example": "12345"
						},
						"userStatus": {
							"type": "integer",
							"description": "User Status",
							"format": "int32",
							"example": 1
						},
						"phone": {
							"type": "string",
							"example": "12345"
						},
						"id": {
							"type": "integer",
							"format": "int64",
							"example": 10
						},
						"email": {
							"type": "string",
							"example": "john@email.com"
						},
						"username": {
							"type": "string",
							"example": "theUser"
						}
					},
					"xml": {
						"name": "user"
					},
					"example": null
				},
				"Address": {
					"type": "object",
					"properties": {
						"zip": {
							"type": "string",
							"example": "94301"
						},
						"city": {
							"type": "string",
							"example": "Palo Alto"
						},
						"street": {
							"type": "string",
							"example": "437 Lytton"
						},
						"state": {
							"type": "string",
							"example": "CA"
						}
					},
					"xml": {
						"name": "address"
					},
					"example": null
				},
				"Customer": {
					"type": "object",
					"properties": {
						"address": {
							"type": "array",
							"xml": {
								"name": "addresses",
								"wrapped": true
							},
							"items": {
								"$ref": "#/components/schemas/Address"
							},
							"example": null
						},
						"id": {
							"type": "integer",
							"format": "int64",
							"example": 100000
						},
						"username": {
							"type": "string",
							"example": "fehguy"
						}
					},
					"xml": {
						"name": "customer"
					},
					"example": null
				},
				"Tag": {
					"type": "object",
					"properties": {
						"name": {
							"type": "string",
							"example": null
						},
						"id": {
							"type": "integer",
							"format": "int64",
							"example": null
						}
					},
					"xml": {
						"name": "tag"
					},
					"example": null
				},
				"ApiResponse": {
					"type": "object",
					"properties": {
						"code": {
							"type": "integer",
							"format": "int32",
							"example": null
						},
						"message": {
							"type": "string",
							"example": null
						},
						"type": {
							"type": "string",
							"example": null
						}
					},
					"xml": {
						"name": "##default"
					},
					"example": null
				},
				"Pet": {
					"required": [
						"name",
						"photoUrls"
					],
					"type": "object",
					"properties": {
						"photoUrls": {
							"type": "array",
							"xml": {
								"wrapped": true
							},
							"items": {
								"type": "string",
								"xml": {
									"name": "photoUrl"
								},
								"example": null
							},
							"example": null
						},
						"name": {
							"type": "string",
							"example": "doggie"
						},
						"id": {
							"type": "integer",
							"format": "int64",
							"example": 10
						},
						"category": {
							"$ref": "#/components/schemas/Category"
						},
						"status": {
							"type": "string",
							"description": "pet status in the store",
							"enum": [
								"available",
								"pending",
								"sold"
							],
							"example": null
						},
						"tags": {
							"type": "array",
							"xml": {
								"wrapped": true
							},
							"items": {
								"$ref": "#/components/schemas/Tag"
							},
							"example": null
						}
					},
					"xml": {
						"name": "pet"
					},
					"example": null
				}
			},
			"responses": {},
			"parameters": {},
			"examples": {},
			"requestBodies": {
				"UserArray": {
					"content": {
						"application/json": {
							"schema": {
								"type": "array",
								"items": {
									"$ref": "#/components/schemas/User"
								},
								"example": null
							},
							"examples": {},
							"example": null
						}
					}
				},
				"Pet": {
					"content": {
						"application/json": {
							"schema": {
								"$ref": "#/components/schemas/Pet"
							},
							"examples": {},
							"example": null
						},
						"application/xml": {
							"schema": {
								"$ref": "#/components/schemas/Pet"
							},
							"examples": {},
							"example": null
						}
					}
				}
			},
			"headers": {},
			"securitySchemes": {
				"apiKey": {
					"type": "apiKey",
					"name": "x-Gateway-APIKey",
					"in": "header"
				}
			},
			"links": {},
			"callbacks": {}
		}
	}`
	mc := &MockClient{}

	webMethodsClient, _ := NewClient(cfg, mc)

	mc.SendFunc = func(request coreapi.Request) (*coreapi.Response, error) {
		return &coreapi.Response{
			Code: 200,
			Body: []byte(response),
		}, nil
	}

	apiResponse, err := webMethodsClient.GetApiSpec("2b598e47-3e0c-4b0f-8a72-da7fdf6c5ea2")
	assert.Nil(t, err)
	assert.NotNil(t, apiResponse)

	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(apiResponse, &jsonMap)
	assert.Nil(t, err)
	_, isOpenAPI := jsonMap["openapi"]
	assert.True(t, isOpenAPI)
}

func TestGetWsdl(t *testing.T) {

	response := `<wsdl:definitions xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/" xmlns:ns="http://c.b.a" xmlns:wsaw="http://www.w3.org/2006/05/addressing/wsdl" xmlns:mime="http://schemas.xmlsoap.org/wsdl/mime/" xmlns:http="http://schemas.xmlsoap.org/wsdl/http/" xmlns:soap12="http://schemas.xmlsoap.org/wsdl/soap12/" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:ns1="http://org.apache.axis2/xsd" xmlns:soap="http://schemas.xmlsoap.org/wsdl/soap/" targetNamespace="http://c.b.a">
    <wsdl:documentation>Calculator</wsdl:documentation>
    <wsdl:types>
        <xs:schema attributeFormDefault="qualified" elementFormDefault="qualified" targetNamespace="http://c.b.a">
            <xs:element name="add">
                <xs:complexType>
                    <xs:sequence>
                        <xs:element minOccurs="0" name="n1" type="xs:int"/>
                        <xs:element minOccurs="0" name="n2" type="xs:int"/>
                    </xs:sequence>
                </xs:complexType>
            </xs:element>
            <xs:element name="addResponse">
                <xs:complexType>
                    <xs:sequence>
                        <xs:element minOccurs="0" name="return" type="xs:int"/>
                    </xs:sequence>
                </xs:complexType>
            </xs:element>
        </xs:schema>
    </wsdl:types>
    <wsdl:message name="addRequest">
        <wsdl:part name="parameters" element="ns:add">
    </wsdl:part>
    </wsdl:message>
    <wsdl:message name="addResponse">
        <wsdl:part name="parameters" element="ns:addResponse">
    </wsdl:part>
    </wsdl:message>
    <wsdl:portType name="CalculatorPortType">
        <wsdl:operation name="add">
            <wsdl:input message="ns:addRequest" wsaw:Action="urn:add">
    </wsdl:input>
            <wsdl:output message="ns:addResponse" wsaw:Action="urn:addResponse">
    </wsdl:output>
        </wsdl:operation>
    </wsdl:portType>
    <wsdl:binding name="CalculatorSoap11Binding" type="ns:CalculatorPortType">
        <soap:binding style="document" transport="http://schemas.xmlsoap.org/soap/http"/>
        <wsdl:operation name="add">
            <soap:operation soapAction="urn:add" style="document"/>
            <wsdl:input>
                <soap:body use="literal"/>
            </wsdl:input>
            <wsdl:output>
                <soap:body use="literal"/>
            </wsdl:output>
        </wsdl:operation>
    </wsdl:binding>
    <wsdl:binding name="CalculatorHttpBinding" type="ns:CalculatorPortType">
        <http:binding verb="POST"/>
        <wsdl:operation name="add">
            <http:operation location="add"/>
            <wsdl:input>
                <mime:content part="parameters" type="text/xml"/>
            </wsdl:input>
            <wsdl:output>
                <mime:content part="parameters" type="text/xml"/>
            </wsdl:output>
        </wsdl:operation>
    </wsdl:binding>
    <wsdl:service name="Calculator">
        <wsdl:port name="CalculatorHttpSoap11Endpoint" binding="ns:CalculatorSoap11Binding">
            <soap:address location="http://env688761.apigw-aw-us.webmethods.io/ws/Calculator.CalculatorHttpSoap11Endpoint/1"/>
        </wsdl:port>
        <wsdl:port name="CalculatorHttpSoap11Endpoint2" binding="ns:CalculatorSoap11Binding">
            <soap:address location="http://env688761.apigw-aw-us.webmethods.io/ws/Calculator/1"/>
        </wsdl:port>
        <wsdl:port name="CalculatorHttpEndpoint" binding="ns:CalculatorHttpBinding">
            <http:address location="http://env688761.apigw-aw-us.webmethods.io/ws/Calculator.CalculatorHttpEndpoint/1"/>
        </wsdl:port>
        <wsdl:port name="CalculatorHttpEndpoint2" binding="ns:CalculatorHttpBinding">
            <http:address location="http://env688761.apigw-aw-us.webmethods.io/ws/Calculator/1"/>
        </wsdl:port>
    </wsdl:service>
</wsdl:definitions>`
	mc := &MockClient{}
	webMethodsClient, _ := NewClient(cfg, mc)
	mc.SendFunc = func(request coreapi.Request) (*coreapi.Response, error) {
		return &coreapi.Response{
			Code: 200,
			Body: []byte(response),
		}, nil
	}
	apiResponse, err := webMethodsClient.GetWsdl("1178fcb2-faae-4fe4-94fa-f5efb0316285")
	assert.Nil(t, err)
	assert.NotNil(t, apiResponse)

}

func TestGetApplication(t *testing.T) {

	response := `{
		"applications": [
			{
				"name": "consumerapp",
				"description": "Testing subscription",
				"contactEmails": [],
				"identifiers": [],
				"siteURLs": [],
				"jsOrigins": [],
				"authStrategyIds": [],
				"version": null,
				"id": "1cc88555-b7df-4e5b-a9e3-1728cc0ecfe6",
				"created": "2023-03-28 02:25:04 GMT",
				"lastupdated": null,
				"owner": "wecare@apiwheel.dev",
				"ownerType": "user",
				"subscription": false,
				"consumingAPIs": [],
				"teams": [],
				"accessTokens": {
					"apiAccessKey_credentials": {
						"apiAccessKey": "e5c8641a-9b41-4b6b-8db2-4c75b24ed104",
						"expirationInterval": null,
						"expirationDate": null,
						"expired": false
					}
				}
			}
		]
	}`
	mc := &MockClient{}
	webMethodsClient, _ := NewClient(cfg, mc)
	mc.SendFunc = func(request coreapi.Request) (*coreapi.Response, error) {
		return &coreapi.Response{
			Code: 200,
			Body: []byte(response),
		}, nil
	}
	applicationResponse, err := webMethodsClient.GetApplication("1cc88555-b7df-4e5b-a9e3-1728cc0ecfe6")
	assert.Nil(t, err)
	assert.Equal(t, applicationResponse.Applications[0].Id, "1cc88555-b7df-4e5b-a9e3-1728cc0ecfe6")
	assert.Equal(t, applicationResponse.Applications[0].Name, "consumerapp")
	assert.Equal(t, applicationResponse.Applications[0].AccessTokens.ApiAccessKeyCredentials.ApiAccessKey, "e5c8641a-9b41-4b6b-8db2-4c75b24ed104")
}

func TestCreateApplication(t *testing.T) {

	request := `{
		"name": "consumerapp",
		"description": "Testing subscription",
		"siteURLs": [],
		"jsOrigins": [],
		"authStrategyIds": [],
		"subscription": true,
		"shell": false
	  }
	  `

	response := `{
		"name": "consumerapp",
		"description": "Testing subscription",
		"contactEmails": [],
		"identifiers": [],
		"siteURLs": [],
		"jsOrigins": [],
		"authStrategyIds": [],
		"version": null,
		"id": "1cc88555-b7df-4e5b-a9e3-1728cc0ecfe6",
		"created": null,
		"lastupdated": null,
		"owner": "wecare@apiwheel.dev",
		"ownerType": "user",
		"subscription": false,
		"consumingAPIs": [],
		"teams": [],
		"accessTokens": {}
	}`
	mc := &MockClient{}
	webMethodsClient, _ := NewClient(cfg, mc)
	mc.SendFunc = func(request coreapi.Request) (*coreapi.Response, error) {
		return &coreapi.Response{
			Code: 200,
			Body: []byte(response),
		}, nil
	}
	applicationRequest := &Application{}
	json.Unmarshal([]byte(request), applicationRequest)
	applicationResponse, err := webMethodsClient.CreateApplication(applicationRequest)
	assert.Nil(t, err)
	assert.Equal(t, applicationResponse.Id, "1cc88555-b7df-4e5b-a9e3-1728cc0ecfe6")
	assert.Equal(t, applicationResponse.Name, "consumerapp")

}

func TestSubscribeApplication(t *testing.T) {

	request := `{
		"apiIDs": [
			"2b598e47-3e0c-4b0f-8a72-da7fdf6c5ea2"
		]
	}`

	response := `{}`
	mc := &MockClient{}
	webMethodsClient, _ := NewClient(cfg, mc)
	mc.SendFunc = func(request coreapi.Request) (*coreapi.Response, error) {
		return &coreapi.Response{
			Code: 201,
			Body: []byte(response),
		}, nil
	}
	applicationApiSubscription := &ApplicationApiSubscription{}
	json.Unmarshal([]byte(request), applicationApiSubscription)
	err := webMethodsClient.SubscribeApplication("1cc88555-b7df-4e5b-a9e3-1728cc0ecfe6", applicationApiSubscription)
	assert.Nil(t, err)
}

func TestRotateApplicationApikey(t *testing.T) {

	response := `{}`
	mc := &MockClient{}
	webMethodsClient, _ := NewClient(cfg, mc)
	mc.SendFunc = func(request coreapi.Request) (*coreapi.Response, error) {
		return &coreapi.Response{
			Code: 201,
			Body: []byte(response),
		}, nil
	}
	err := webMethodsClient.RotateApplicationApikey("1cc88555-b7df-4e5b-a9e3-1728cc0ecfe6")
	assert.Nil(t, err)
}

func TestDeleteApplication(t *testing.T) {

	mc := &MockClient{}
	webMethodsClient, _ := NewClient(cfg, mc)
	mc.SendFunc = func(request coreapi.Request) (*coreapi.Response, error) {
		return &coreapi.Response{
			Code: 204,
		}, nil
	}
	err := webMethodsClient.DeleteApplication("1cc88555-b7df-4e5b-a9e3-1728cc0ecfe6")
	assert.Nil(t, err)
}

func TestDeleteApplicationAccessTokens(t *testing.T) {

	mc := &MockClient{}
	webMethodsClient, _ := NewClient(cfg, mc)
	mc.SendFunc = func(request coreapi.Request) (*coreapi.Response, error) {
		return &coreapi.Response{
			Code: 204,
		}, nil
	}
	err := webMethodsClient.DeleteApplicationAccessTokens("1cc88555-b7df-4e5b-a9e3-1728cc0ecfe6")
	assert.Nil(t, err)
}

func TestUnsubscribeApplication(t *testing.T) {

	mc := &MockClient{}
	webMethodsClient, _ := NewClient(cfg, mc)
	mc.SendFunc = func(request coreapi.Request) (*coreapi.Response, error) {
		return &coreapi.Response{
			Code: 204,
		}, nil
	}
	err := webMethodsClient.UnsubscribeApplication("1cc88555-b7df-4e5b-a9e3-1728cc0ecfe6", "1178fcb2-faae-4fe4-94fa-f5efb0316285")
	assert.Nil(t, err)
}

func TestIsAllowedTags(t *testing.T) {

	responseStr := `{
		"apiResponse": {
			"api": {
				"apiDefinition": {
					"info": {
						"description": "This is a sample Pet Store Server based on the OpenAPI 3.0 specification.  You can find out more about\nSwagger at [http://swagger.io](http://swagger.io). In the third iteration of the pet store, we've switched to the design first approach!\nYou can now help us improve the API whether it's by making changes to the definition itself or to the code.\nThat way, with time, we can improve the API in general, and expose some of the new features in OAS3.\n\nSome useful links:\n- [The Pet Store repository](https://github.com/swagger-api/swagger-petstore)\n- [The source API definition for the Pet Store](https://github.com/swagger-api/swagger-petstore/blob/master/src/main/resources/openapi.yaml)",
						"version": "1.0.17",
						"title": "Swagger Petstore - OpenAPI 3.0",
						"termsOfService": "http://swagger.io/terms/",
						"contact": {
							"email": "apiteam@swagger.io"
						},
						"license": {
							"name": "Apache 2.0",
							"url": "http://www.apache.org/licenses/LICENSE-2.0.html"
						}
					},
					"serviceRegistryDisplayName": "petstore_1.0.17",
					"tags": [
						{
							"name": "customrathna"
						},
						{
							"name": "custom2"
						},
						{
							"name": "pet",
							"description": "Everything about your Pets",
							"externalDocs": {
								"description": "Find out more",
								"url": "http://swagger.io"
							}
						},
						{
							"name": "store",
							"description": "Access to Petstore orders",
							"externalDocs": {
								"description": "Find out more about our store",
								"url": "http://swagger.io"
							}
						},
						{
							"name": "user",
							"description": "Operations about user"
						}
					],
					"schemes": [],
					"security": [],
					"paths": {
						"/pet": {
							"put": {
								"tags": [
									"pet",
									"customrathna"
								],
								"summary": "Update an existing pet",
								"description": "Update an existing pet by Id",
								"operationId": "updatePet",
								"consumes": [
									"application/xml",
									"application/json",
									"application/x-www-form-urlencoded"
								],
								"produces": [
									"application/xml",
									"application/json"
								],
								"parameters": [],
								"responses": {
									"200": {
										"description": "Successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"$ref": "#/components/schemas/Pet"
												},
												"examples": {}
											},
											"application/xml": {
												"schema": {
													"$ref": "#/components/schemas/Pet"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"400": {
										"description": "Invalid ID supplied",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"404": {
										"description": "Pet not found",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"405": {
										"description": "Validation exception",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"security": [
									{
										"requirements": {
											"petstore_auth": [
												"write:pets",
												"read:pets"
											]
										}
									}
								],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {
										"application/json": {
											"schema": {
												"$ref": "#/components/schemas/Pet"
											},
											"examples": {}
										},
										"application/x-www-form-urlencoded": {
											"schema": {
												"$ref": "#/components/schemas/Pet"
											},
											"examples": {}
										},
										"application/xml": {
											"schema": {
												"$ref": "#/components/schemas/Pet"
											},
											"examples": {}
										}
									},
									"name": "updatePet"
								}
							},
							"post": {
								"tags": [
									"pet",
									"custom2"
								],
								"summary": "Add a new pet to the store",
								"description": "Add a new pet to the store",
								"operationId": "addPet",
								"consumes": [
									"application/xml",
									"application/json",
									"application/x-www-form-urlencoded"
								],
								"produces": [
									"application/xml",
									"application/json"
								],
								"parameters": [],
								"responses": {
									"200": {
										"description": "Successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"$ref": "#/components/schemas/Pet"
												},
												"examples": {}
											},
											"application/xml": {
												"schema": {
													"$ref": "#/components/schemas/Pet"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"405": {
										"description": "Invalid input",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"security": [
									{
										"requirements": {
											"petstore_auth": [
												"write:pets",
												"read:pets"
											]
										}
									}
								],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {
										"application/json": {
											"schema": {
												"$ref": "#/components/schemas/Pet"
											},
											"examples": {}
										},
										"application/x-www-form-urlencoded": {
											"schema": {
												"$ref": "#/components/schemas/Pet"
											},
											"examples": {}
										},
										"application/xml": {
											"schema": {
												"$ref": "#/components/schemas/Pet"
											},
											"examples": {}
										}
									},
									"name": "addPet"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/pet",
							"enabled": true
						},
						"/pet/findByStatus": {
							"get": {
								"tags": [
									"pet"
								],
								"summary": "Finds Pets by status",
								"description": "Multiple status values can be provided with comma separated strings",
								"operationId": "findPetsByStatus",
								"produces": [
									"application/xml",
									"application/json"
								],
								"parameters": [
									{
										"default": "available",
										"description": "Status values that need to be considered for filter",
										"explode": true,
										"in": "query",
										"name": "status",
										"parameterSchema": {
											"default": "available",
											"enum": [
												"available",
												"pending",
												"sold"
											],
											"type": "string"
										},
										"required": false,
										"style": "FORM",
										"type": "string"
									}
								],
								"responses": {
									"200": {
										"description": "successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"type": "gateway",
													"schema": "{\"type\":\"array\",\"items\":{\"type\":\"object\",\"$ref\":\"#/components/schemas/Pet\"}}"
												},
												"examples": {}
											},
											"application/xml": {
												"schema": {
													"type": "gateway",
													"schema": "{\"type\":\"array\",\"items\":{\"type\":\"object\",\"$ref\":\"#/components/schemas/Pet\"}}"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"400": {
										"description": "Invalid status value",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"security": [
									{
										"requirements": {
											"petstore_auth": [
												"write:pets",
												"read:pets"
											]
										}
									}
								],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "findPetsByStatus"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/pet/findByStatus",
							"enabled": true
						},
						"/pet/findByTags": {
							"get": {
								"tags": [
									"pet"
								],
								"summary": "Finds Pets by tags",
								"description": "Multiple tags can be provided with comma separated strings. Use tag1, tag2, tag3 for testing.",
								"operationId": "findPetsByTags",
								"produces": [
									"application/xml",
									"application/json"
								],
								"parameters": [
									{
										"description": "Tags to filter by",
										"explode": true,
										"in": "query",
										"name": "tags",
										"parameterSchema": {
											"items": {
												"type": "string"
											},
											"type": "array"
										},
										"required": false,
										"style": "FORM"
									}
								],
								"responses": {
									"200": {
										"description": "successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"type": "gateway",
													"schema": "{\"type\":\"array\",\"items\":{\"type\":\"object\",\"$ref\":\"#/components/schemas/Pet\"}}"
												},
												"examples": {}
											},
											"application/xml": {
												"schema": {
													"type": "gateway",
													"schema": "{\"type\":\"array\",\"items\":{\"type\":\"object\",\"$ref\":\"#/components/schemas/Pet\"}}"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"400": {
										"description": "Invalid tag value",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"security": [
									{
										"requirements": {
											"petstore_auth": [
												"write:pets",
												"read:pets"
											]
										}
									}
								],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "findPetsByTags"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/pet/findByTags",
							"enabled": true
						},
						"/pet/{petId}": {
							"get": {
								"tags": [
									"pet"
								],
								"summary": "Find pet by ID",
								"description": "Returns a single pet",
								"operationId": "getPetById",
								"produces": [
									"application/xml",
									"application/json"
								],
								"parameters": [
									{
										"description": "ID of pet to return",
										"explode": false,
										"format": "int64",
										"in": "path",
										"name": "petId",
										"parameterSchema": {
											"format": "int64",
											"type": "integer"
										},
										"required": true,
										"style": "SIMPLE",
										"type": "integer"
									}
								],
								"responses": {
									"200": {
										"description": "successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"$ref": "#/components/schemas/Pet"
												},
												"examples": {}
											},
											"application/xml": {
												"schema": {
													"$ref": "#/components/schemas/Pet"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"400": {
										"description": "Invalid ID supplied",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"404": {
										"description": "Pet not found",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"security": [
									{
										"requirements": {
											"api_key": []
										}
									},
									{
										"requirements": {
											"petstore_auth": [
												"write:pets",
												"read:pets"
											]
										}
									}
								],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "getPetById"
								}
							},
							"post": {
								"tags": [
									"pet"
								],
								"summary": "Updates a pet in the store with form data",
								"description": "",
								"operationId": "updatePetWithForm",
								"parameters": [
									{
										"description": "ID of pet that needs to be updated",
										"explode": false,
										"format": "int64",
										"in": "path",
										"name": "petId",
										"parameterSchema": {
											"format": "int64",
											"type": "integer"
										},
										"required": true,
										"style": "SIMPLE",
										"type": "integer"
									},
									{
										"description": "Name of pet that needs to be updated",
										"explode": true,
										"in": "query",
										"name": "name",
										"parameterSchema": {
											"type": "string"
										},
										"required": false,
										"style": "FORM",
										"type": "string"
									},
									{
										"description": "Status of pet that needs to be updated",
										"explode": true,
										"in": "query",
										"name": "status",
										"parameterSchema": {
											"type": "string"
										},
										"required": false,
										"style": "FORM",
										"type": "string"
									}
								],
								"responses": {
									"405": {
										"description": "Invalid input",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"security": [
									{
										"requirements": {
											"petstore_auth": [
												"write:pets",
												"read:pets"
											]
										}
									}
								],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "updatePetWithForm"
								}
							},
							"delete": {
								"tags": [
									"pet"
								],
								"summary": "Deletes a pet",
								"description": "",
								"operationId": "deletePet",
								"parameters": [
									{
										"description": "",
										"explode": false,
										"in": "header",
										"name": "api_key",
										"parameterSchema": {
											"type": "string"
										},
										"required": false,
										"style": "SIMPLE",
										"type": "string"
									},
									{
										"description": "Pet id to delete",
										"explode": false,
										"format": "int64",
										"in": "path",
										"name": "petId",
										"parameterSchema": {
											"format": "int64",
											"type": "integer"
										},
										"required": true,
										"style": "SIMPLE",
										"type": "integer"
									}
								],
								"responses": {
									"400": {
										"description": "Invalid pet value",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"security": [
									{
										"requirements": {
											"petstore_auth": [
												"write:pets",
												"read:pets"
											]
										}
									}
								],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "deletePet"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/pet/{petId}",
							"enabled": true
						},
						"/pet/{petId}/uploadImage": {
							"post": {
								"tags": [
									"pet"
								],
								"summary": "uploads an image",
								"description": "",
								"operationId": "uploadFile",
								"consumes": [
									"application/octet-stream"
								],
								"produces": [
									"application/json"
								],
								"parameters": [
									{
										"description": "ID of pet to update",
										"explode": false,
										"format": "int64",
										"in": "path",
										"name": "petId",
										"parameterSchema": {
											"format": "int64",
											"type": "integer"
										},
										"required": true,
										"style": "SIMPLE",
										"type": "integer"
									},
									{
										"description": "Additional Metadata",
										"explode": true,
										"in": "query",
										"name": "additionalMetadata",
										"parameterSchema": {
											"type": "string"
										},
										"required": false,
										"style": "FORM",
										"type": "string"
									}
								],
								"responses": {
									"200": {
										"description": "successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"$ref": "#/components/schemas/ApiResponse"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"security": [
									{
										"requirements": {
											"petstore_auth": [
												"write:pets",
												"read:pets"
											]
										}
									}
								],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {
										"application/octet-stream": {
											"schema": {
												"type": "gateway",
												"schema": "{\"type\":\"string\",\"format\":\"binary\"}"
											},
											"examples": {}
										}
									},
									"name": "uploadFile"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/pet/{petId}/uploadImage",
							"enabled": true
						},
						"/store/inventory": {
							"get": {
								"tags": [
									"store"
								],
								"summary": "Returns pet inventories by status",
								"description": "Returns a map of status codes to quantities",
								"operationId": "getInventory",
								"produces": [
									"application/json"
								],
								"parameters": [],
								"responses": {
									"200": {
										"description": "successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"type": "gateway",
													"schema": "{\"type\":\"object\",\"additionalProperties\":{\"type\":\"integer\",\"format\":\"int32\"}}"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"security": [
									{
										"requirements": {
											"api_key": []
										}
									}
								],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "getInventory"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/store/inventory",
							"enabled": true
						},
						"/store/order": {
							"post": {
								"tags": [
									"store"
								],
								"summary": "Place an order for a pet",
								"description": "Place a new order in the store",
								"operationId": "placeOrder",
								"consumes": [
									"application/xml",
									"application/json",
									"application/x-www-form-urlencoded"
								],
								"produces": [
									"application/json"
								],
								"parameters": [],
								"responses": {
									"200": {
										"description": "successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"$ref": "#/components/schemas/Order"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"405": {
										"description": "Invalid input",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {
										"application/json": {
											"schema": {
												"$ref": "#/components/schemas/Order"
											},
											"examples": {}
										},
										"application/x-www-form-urlencoded": {
											"schema": {
												"$ref": "#/components/schemas/Order"
											},
											"examples": {}
										},
										"application/xml": {
											"schema": {
												"$ref": "#/components/schemas/Order"
											},
											"examples": {}
										}
									},
									"name": "placeOrder"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/store/order",
							"enabled": true
						},
						"/store/order/{orderId}": {
							"get": {
								"tags": [
									"store"
								],
								"summary": "Find purchase order by ID",
								"description": "For valid response try integer IDs with value <= 5 or > 10. Other values will generate exceptions.",
								"operationId": "getOrderById",
								"produces": [
									"application/xml",
									"application/json"
								],
								"parameters": [
									{
										"description": "ID of order that needs to be fetched",
										"explode": false,
										"format": "int64",
										"in": "path",
										"name": "orderId",
										"parameterSchema": {
											"format": "int64",
											"type": "integer"
										},
										"required": true,
										"style": "SIMPLE",
										"type": "integer"
									}
								],
								"responses": {
									"200": {
										"description": "successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"$ref": "#/components/schemas/Order"
												},
												"examples": {}
											},
											"application/xml": {
												"schema": {
													"$ref": "#/components/schemas/Order"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"400": {
										"description": "Invalid ID supplied",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"404": {
										"description": "Order not found",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "getOrderById"
								}
							},
							"delete": {
								"tags": [
									"store"
								],
								"summary": "Delete purchase order by ID",
								"description": "For valid response try integer IDs with value < 1000. Anything above 1000 or nonintegers will generate API errors",
								"operationId": "deleteOrder",
								"parameters": [
									{
										"description": "ID of the order that needs to be deleted",
										"explode": false,
										"format": "int64",
										"in": "path",
										"name": "orderId",
										"parameterSchema": {
											"format": "int64",
											"type": "integer"
										},
										"required": true,
										"style": "SIMPLE",
										"type": "integer"
									}
								],
								"responses": {
									"400": {
										"description": "Invalid ID supplied",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"404": {
										"description": "Order not found",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "deleteOrder"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/store/order/{orderId}",
							"enabled": true
						},
						"/user": {
							"post": {
								"tags": [
									"user"
								],
								"summary": "Create user",
								"description": "This can only be done by the logged in user.",
								"operationId": "createUser",
								"consumes": [
									"application/xml",
									"application/json",
									"application/x-www-form-urlencoded"
								],
								"produces": [
									"application/xml",
									"application/json"
								],
								"parameters": [],
								"responses": {
									"default": {
										"description": "successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"$ref": "#/components/schemas/User"
												},
												"examples": {}
											},
											"application/xml": {
												"schema": {
													"$ref": "#/components/schemas/User"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {
										"application/json": {
											"schema": {
												"$ref": "#/components/schemas/User"
											},
											"examples": {}
										},
										"application/x-www-form-urlencoded": {
											"schema": {
												"$ref": "#/components/schemas/User"
											},
											"examples": {}
										},
										"application/xml": {
											"schema": {
												"$ref": "#/components/schemas/User"
											},
											"examples": {}
										}
									},
									"name": "createUser"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/user",
							"enabled": true
						},
						"/user/createWithList": {
							"post": {
								"tags": [
									"user"
								],
								"summary": "Creates list of users with given input array",
								"description": "Creates list of users with given input array",
								"operationId": "createUsersWithListInput",
								"consumes": [
									"application/json"
								],
								"produces": [
									"application/xml",
									"application/json"
								],
								"parameters": [],
								"responses": {
									"200": {
										"description": "Successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"$ref": "#/components/schemas/User"
												},
												"examples": {}
											},
											"application/xml": {
												"schema": {
													"$ref": "#/components/schemas/User"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"default": {
										"description": "successful operation",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {
										"application/json": {
											"schema": {
												"type": "gateway",
												"schema": "{\"type\":\"array\",\"items\":{\"type\":\"object\",\"$ref\":\"#/components/schemas/User\"}}"
											},
											"examples": {}
										}
									},
									"name": "createUsersWithListInput"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/user/createWithList",
							"enabled": true
						},
						"/user/login": {
							"get": {
								"tags": [
									"user"
								],
								"summary": "Logs user into the system",
								"description": "",
								"operationId": "loginUser",
								"produces": [
									"application/xml",
									"application/json"
								],
								"parameters": [
									{
										"description": "The user name for login",
										"explode": true,
										"in": "query",
										"name": "username",
										"parameterSchema": {
											"type": "string"
										},
										"required": false,
										"style": "FORM",
										"type": "string"
									},
									{
										"description": "The password for login in clear text",
										"explode": true,
										"in": "query",
										"name": "password",
										"parameterSchema": {
											"type": "string"
										},
										"required": false,
										"style": "FORM",
										"type": "string"
									}
								],
								"responses": {
									"200": {
										"description": "successful operation",
										"headersV3": {
											"X-Expires-After": {
												"name": "X-Expires-After",
												"in": "header",
												"description": "date in UTC when token expires",
												"required": false,
												"type": "string",
												"format": "date-time",
												"style": "SIMPLE",
												"explode": false,
												"examples": {},
												"parameterSchema": {
													"type": "string",
													"format": "date-time"
												}
											},
											"X-Rate-Limit": {
												"name": "X-Rate-Limit",
												"in": "header",
												"description": "calls per hour allowed by the user",
												"required": false,
												"type": "integer",
												"format": "int32",
												"style": "SIMPLE",
												"explode": false,
												"examples": {},
												"parameterSchema": {
													"type": "integer",
													"format": "int32"
												}
											}
										},
										"content": {
											"application/json": {
												"schema": {
													"type": "gateway",
													"schema": "{\"type\":\"string\"}"
												},
												"examples": {}
											},
											"application/xml": {
												"schema": {
													"type": "gateway",
													"schema": "{\"type\":\"string\"}"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"400": {
										"description": "Invalid username/password supplied",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "loginUser"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/user/login",
							"enabled": true
						},
						"/user/logout": {
							"get": {
								"tags": [
									"user"
								],
								"summary": "Logs out current logged in user session",
								"description": "",
								"operationId": "logoutUser",
								"parameters": [],
								"responses": {
									"default": {
										"description": "successful operation",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "logoutUser"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/user/logout",
							"enabled": true
						},
						"/user/{username}": {
							"get": {
								"tags": [
									"user"
								],
								"summary": "Get user by user name",
								"description": "",
								"operationId": "getUserByName",
								"produces": [
									"application/xml",
									"application/json"
								],
								"parameters": [
									{
										"description": "The name that needs to be fetched. Use user1 for testing. ",
										"explode": false,
										"in": "path",
										"name": "username",
										"parameterSchema": {
											"type": "string"
										},
										"required": true,
										"style": "SIMPLE",
										"type": "string"
									}
								],
								"responses": {
									"200": {
										"description": "successful operation",
										"headersV3": {},
										"content": {
											"application/json": {
												"schema": {
													"$ref": "#/components/schemas/User"
												},
												"examples": {}
											},
											"application/xml": {
												"schema": {
													"$ref": "#/components/schemas/User"
												},
												"examples": {}
											}
										},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"400": {
										"description": "Invalid username supplied",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"404": {
										"description": "User not found",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "getUserByName"
								}
							},
							"put": {
								"tags": [
									"user"
								],
								"summary": "Update user",
								"description": "This can only be done by the logged in user.",
								"operationId": "updateUser",
								"consumes": [
									"application/xml",
									"application/json",
									"application/x-www-form-urlencoded"
								],
								"parameters": [
									{
										"description": "name that need to be deleted",
										"explode": false,
										"in": "path",
										"name": "username",
										"parameterSchema": {
											"type": "string"
										},
										"required": true,
										"style": "SIMPLE",
										"type": "string"
									}
								],
								"responses": {
									"default": {
										"description": "successful operation",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {
										"application/json": {
											"schema": {
												"$ref": "#/components/schemas/User"
											},
											"examples": {}
										},
										"application/x-www-form-urlencoded": {
											"schema": {
												"$ref": "#/components/schemas/User"
											},
											"examples": {}
										},
										"application/xml": {
											"schema": {
												"$ref": "#/components/schemas/User"
											},
											"examples": {}
										}
									},
									"name": "updateUser"
								}
							},
							"delete": {
								"tags": [
									"user"
								],
								"summary": "Delete user",
								"description": "This can only be done by the logged in user.",
								"operationId": "deleteUser",
								"parameters": [
									{
										"description": "The name that needs to be deleted",
										"explode": false,
										"in": "path",
										"name": "username",
										"parameterSchema": {
											"type": "string"
										},
										"required": true,
										"style": "SIMPLE",
										"type": "string"
									}
								],
								"responses": {
									"400": {
										"description": "Invalid username supplied",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									},
									"404": {
										"description": "User not found",
										"headersV3": {},
										"content": {},
										"links": {},
										"schema": {},
										"examples": {},
										"headers": {}
									}
								},
								"mockedResponses": {},
								"mockedConditionsBasedCustomResponsesList": [],
								"enabled": true,
								"scopes": [],
								"requestBody": {
									"content": {},
									"name": "deleteUser"
								}
							},
							"parameters": [],
							"scopes": [],
							"displayName": "/user/{username}",
							"enabled": true
						}
					},
					"securityDefinitions": {},
					"definitions": {},
					"parameters": {},
					"baseUriParameters": [],
					"externalDocs": [
						{
							"description": "Find out more about Swagger",
							"url": "http://swagger.io"
						}
					],
					"servers": [
						{
							"url": "/api/v3",
							"variables": {}
						}
					],
					"components": {
						"schemas": {
							"Address": {
								"type": "gateway",
								"schema": "{\"type\":\"object\",\"properties\":{\"city\":{\"type\":\"string\",\"example\":\"Palo Alto\"},\"state\":{\"type\":\"string\",\"example\":\"CA\"},\"street\":{\"type\":\"string\",\"example\":\"437 Lytton\"},\"zip\":{\"type\":\"string\",\"example\":\"94301\"}},\"xml\":{\"name\":\"address\"}}"
							},
							"ApiResponse": {
								"type": "gateway",
								"schema": "{\"type\":\"object\",\"properties\":{\"code\":{\"type\":\"integer\",\"format\":\"int32\"},\"message\":{\"type\":\"string\"},\"type\":{\"type\":\"string\"}},\"xml\":{\"name\":\"##default\"}}"
							},
							"Category": {
								"type": "gateway",
								"schema": "{\"type\":\"object\",\"properties\":{\"id\":{\"type\":\"integer\",\"format\":\"int64\",\"example\":\"1\"},\"name\":{\"type\":\"string\",\"example\":\"Dogs\"}},\"xml\":{\"name\":\"category\"}}"
							},
							"Customer": {
								"type": "gateway",
								"schema": "{\"type\":\"object\",\"properties\":{\"address\":{\"type\":\"array\",\"xml\":{\"name\":\"addresses\",\"wrapped\":true},\"items\":{\"type\":\"object\",\"$ref\":\"#/components/schemas/Address\"}},\"id\":{\"type\":\"integer\",\"format\":\"int64\",\"example\":\"100000\"},\"username\":{\"type\":\"string\",\"example\":\"fehguy\"}},\"xml\":{\"name\":\"customer\"}}"
							},
							"Order": {
								"type": "gateway",
								"schema": "{\"type\":\"object\",\"properties\":{\"complete\":{\"type\":\"boolean\"},\"id\":{\"type\":\"integer\",\"format\":\"int64\",\"example\":\"10\"},\"petId\":{\"type\":\"integer\",\"format\":\"int64\",\"example\":\"198772\"},\"quantity\":{\"type\":\"integer\",\"format\":\"int32\",\"example\":\"7\"},\"shipDate\":{\"type\":\"string\",\"format\":\"date-time\"},\"status\":{\"type\":\"string\",\"description\":\"Order Status\",\"example\":\"approved\",\"enum\":[\"placed\",\"approved\",\"delivered\"]}},\"xml\":{\"name\":\"order\"}}"
							},
							"Pet": {
								"type": "gateway",
								"schema": "{\"required\":[\"name\",\"photoUrls\"],\"type\":\"object\",\"properties\":{\"category\":{\"type\":\"object\",\"$ref\":\"#/components/schemas/Category\"},\"id\":{\"type\":\"integer\",\"format\":\"int64\",\"example\":\"10\"},\"name\":{\"type\":\"string\",\"example\":\"doggie\"},\"photoUrls\":{\"type\":\"array\",\"xml\":{\"wrapped\":true},\"items\":{\"type\":\"string\",\"xml\":{\"name\":\"photoUrl\"}}},\"status\":{\"type\":\"string\",\"description\":\"pet status in the store\",\"enum\":[\"available\",\"pending\",\"sold\"]},\"tags\":{\"type\":\"array\",\"xml\":{\"wrapped\":true},\"items\":{\"type\":\"object\",\"$ref\":\"#/components/schemas/Tag\"}}},\"xml\":{\"name\":\"pet\"}}"
							},
							"Tag": {
								"type": "gateway",
								"schema": "{\"type\":\"object\",\"properties\":{\"id\":{\"type\":\"integer\",\"format\":\"int64\"},\"name\":{\"type\":\"string\"}},\"xml\":{\"name\":\"tag\"}}"
							},
							"User": {
								"type": "gateway",
								"schema": "{\"type\":\"object\",\"properties\":{\"email\":{\"type\":\"string\",\"example\":\"john@email.com\"},\"firstName\":{\"type\":\"string\",\"example\":\"John\"},\"id\":{\"type\":\"integer\",\"format\":\"int64\",\"example\":\"10\"},\"lastName\":{\"type\":\"string\",\"example\":\"James\"},\"password\":{\"type\":\"string\",\"example\":\"12345\"},\"phone\":{\"type\":\"string\",\"example\":\"12345\"},\"userStatus\":{\"type\":\"integer\",\"description\":\"User Status\",\"format\":\"int32\",\"example\":\"1\"},\"username\":{\"type\":\"string\",\"example\":\"theUser\"}},\"xml\":{\"name\":\"user\"}}"
							}
						},
						"responses": {},
						"parameters": {},
						"examples": {},
						"requestBodies": {
							"Pet": {
								"content": {
									"application/json": {
										"schema": {
											"$ref": "#/components/schemas/Pet"
										},
										"examples": {}
									},
									"application/xml": {
										"schema": {
											"$ref": "#/components/schemas/Pet"
										},
										"examples": {}
									}
								},
								"name": "Pet"
							},
							"UserArray": {
								"content": {
									"application/json": {
										"schema": {
											"type": "gateway",
											"schema": "{\"type\":\"array\",\"items\":{\"type\":\"object\",\"$ref\":\"#/components/schemas/User\"}}"
										},
										"examples": {}
									}
								},
								"name": "UserArray"
							}
						},
						"headers": {},
						"links": {},
						"callbacks": {}
					},
					"type": "rest"
				},
				"nativeEndpoint": [
					{
						"passSecurityHeaders": true,
						"uri": "/api/v3",
						"connectionTimeoutDuration": 0,
						"alias": false
					}
				],
				"apiName": "petstore",
				"apiVersion": "1.0.17",
				"apiDescription": "This is a sample Pet Store Server based on the OpenAPI 3.0 specification.  You can find out more about\nSwagger at [http://swagger.io](http://swagger.io). In the third iteration of the pet store, we've switched to the design first approach!\nYou can now help us improve the API whether it's by making changes to the definition itself or to the code.\nThat way, with time, we can improve the API in general, and expose some of the new features in OAS3.\n\nSome useful links:\n- [The Pet Store repository](https://github.com/swagger-api/swagger-petstore)\n- [The source API definition for the Pet Store](https://github.com/swagger-api/swagger-petstore/blob/master/src/main/resources/openapi.yaml)",
				"maturityState": "Beta",
				"apiGroups": [
					"Finance Banking and Insurance",
					"Sales and Ordering"
				],
				"isActive": true,
				"type": "REST",
				"owner": "wecare@apiwheel.dev",
				"policies": [
					"4631db34-04a7-4cad-860c-6a182c019d4a"
				],
				"tracingEnabled": false,
				"scopes": [],
				"publishedPortals": [
					"69ca5c6b-8b0a-4bae-8021-4ccf0ffd136f"
				],
				"creationDate": "2023-03-18 06:29:48 GMT",
				"lastModified": "2023-03-29 13:43:56 GMT",
				"systemVersion": 1,
				"gatewayEndpoints": {},
				"deployments": [
					"APIGateway"
				],
				"microgatewayEndpoints": [],
				"appMeshEndpoints": [],
				"id": "2b598e47-3e0c-4b0f-8a72-da7fdf6c5ea2"
			},
			"responseStatus": "SUCCESS",
			"gatewayEndPoints": [
				"http://env688761.apigw-aw-us.webmethods.io/gateway/petstore/1.0.17"
			],
			"gatewayEndPointList": [
				{
					"endpointName": "DEFAULT_GATEWAY_ENDPOINT",
					"endpointDisplayName": "Default",
					"endpoint": "gateway/petstore/1.0.17",
					"endpointType": "DEFAULT",
					"endpointUrls": [
						"http://env688761.apigw-aw-us.webmethods.io/gateway/petstore/1.0.17"
					]
				}
			],
			"versions": [
				{
					"versionNumber": "1.0.17",
					"apiId": "2b598e47-3e0c-4b0f-8a72-da7fdf6c5ea2"
				}
			]
		}
	}`

	apiResponse := &GetApiDetails{}
	json.Unmarshal([]byte(responseStr), apiResponse)
	mc := &MockClient{}
	cfg.Filter = `tag.custom2.Exists() && tag.customrathna.Exists()`

	// cfg.Filter = `tag.Contains(custom)` //invalid config
	//cfg.Filter = `tag.custom2 == "custom2"` // Not supported as tag does not have any value.
	webMethodsClient, err := NewClient(cfg, mc)
	assert.Nil(t, err)
	logrus.SetLevel(logrus.DebugLevel)
	result := webMethodsClient.IsAllowedTags(apiResponse.ApiResponse.Api.ApiDefinition.Tags)
	assert.True(t, result)

	cfg.Filter = `tag.custom2.Exists() || tag.customrathna2.Exists()`
	webMethodsClient, err = NewClient(cfg, mc)
	assert.Nil(t, err)
	result = webMethodsClient.IsAllowedTags(apiResponse.ApiResponse.Api.ApiDefinition.Tags)
	assert.True(t, result)

	cfg.Filter = `tag.temp.Exists()`
	webMethodsClient, err = NewClient(cfg, mc)
	assert.Nil(t, err)
	result = webMethodsClient.IsAllowedTags(apiResponse.ApiResponse.Api.ApiDefinition.Tags)
	assert.False(t, result)

	cfg.Filter = `tag.pet.Exists()`
	webMethodsClient, err = NewClient(cfg, mc)
	assert.Nil(t, err)
	result = webMethodsClient.IsAllowedTags(apiResponse.ApiResponse.Api.ApiDefinition.Tags)
	assert.True(t, result)

}
