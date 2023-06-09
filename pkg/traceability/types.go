package traceability

// # Transaction event logs
// 1 EVENT_TYPE
// 2 ROOT_CONTEXT
// 3 PARENT_CONTEXT
// 4 CURRENT_CONTEXT
// 5 UUID
// 6 SERVER_ID
// 7 EVENT_TIMESTAMP
// 8 CURRENT_TIMESTAMP
// 9 SESSION_ID
// 10 API_NAME
// 11 API_VERSION
// 12 TARGET_NAME
// 13 APPLICATION_NAME
// 14 APPLICATION_IP
// 15 APPLICATION_ID
// 16 RESPONSE
// 17 REQUEST
// 18 TOTAL_TIME
// 19 NATIVE_TIME
// 20 REQUEST_STATUS
// 21 OPERATION_NAME
// 22 NATIVE_ENDPOINT
// 23 PARTNER_ID
// 24 API_ID
// 25 SERVICE_NAME
// 26 REQ_HEADERS
// 27 QUERY_PARAM
// 28 RES_HEADERS
// 29 CORRELATIONID
// 30 ERROR_ORIGIN
// 31 CUSTOM
// 32 NATIVE_REQUEST_HEADERS
// 33 NATIVE_REQUEST_PAYLOAD
// 34 NATIVE_RESPONSE_HEADERS
// 35 NATIVE_RESPONSE_PAYLOAD
// 36 NATIVE_HTTP_METHOD GET
// 37 NATIVE_URL
// 38 EXTERNAL_CALLS
// 39 SOURCE_GATEWAY_NODE

// GwTrafficLogEntry - Represents the structure of log entry the agent will receive
type GwTrafficLogEntry struct {
	EventType             string
	RootContext           string
	ParentContext         string
	CurrentContext        string
	Uuid                  string
	ServerId              string
	EventTimestamp        string
	CurrentTimestamp      string
	SessionId             string
	ApiName               string
	ApiVersion            string
	TargetName            string
	ApplicationName       string
	ApplicationIp         string
	ApplicationId         string
	Request               string
	Response              string
	TotalTime             int
	NativeTime            int
	RequestStatus         string
	OperationName         string
	NatvieEndpoint        string
	PartnerId             string
	ApiId                 string
	ServiceName           string
	RequestHeaders        string
	QueryParam            string
	ResponseHeaders       string
	CorrelationId         string
	ErrorOrigin           string
	Custom                string
	NativeRequestHeaders  string
	NativeRequestPayload  string
	NativeResponseHeaders string
	NativeResponsePayload string
	NativeHttpMethod      string
	NativeUrl             string
	ExternalCalls         string
	SourceGatewayNode     string
}

type WebmethodsEvent struct {
	EventType       string      `json:"eventType"`
	SourceGateway   string      `json:"sourceGateway"`
	CreationDate    int64       `json:"creationDate"`
	ApiName         string      `json:"apiName"`
	ApiVersion      string      `json:"apiVersion"`
	ApiId           string      `json:"apiId"`
	TotalTime       int         `json:"totalTime"`
	SessionId       string      `json:"sessionId"`
	GatewayTime     int         `json:"gatewayTime"`
	ApplicationName string      `json:"applicationName"`
	ApplicationIp   string      `json:"applicationIp"`
	ApplicationId   string      `json:"applicationId"`
	Status          string      `json:"status"`
	ReqPayload      string      `json:"reqPayload"`
	ResPayload      string      `json:"resPayload"`
	TotalDataSize   int         `json:"totalDataSize"`
	ResponseCode    string      `json:"responseCode"`
	OperationName   string      `json:"operationName"`
	HTTPMethod      string      `json:"httpMethod"`
	RequestHeaders  HttpHeaders `json:"requestHeaders"`
	ResponseHeaders HttpHeaders `json:"responseHeaders"`
	QueryParameters struct {
	} `json:"queryParameters"`
	CorrelationID string `json:"correlationID"`
	CustomFields  struct {
	} `json:"customFields"`
	ErrorOrigin           string        `json:"errorOrigin"`
	NativeRequestHeaders  HttpHeaders   `json:"nativeRequestHeaders"`
	NativeReqPayload      string        `json:"nativeReqPayload"`
	NativeResponseHeaders HttpHeaders   `json:"nativeResponseHeaders"`
	NativeResPayload      string        `json:"nativeResPayload"`
	NativeHTTPMethod      string        `json:"nativeHttpMethod"`
	NativeURL             string        `json:"nativeURL"`
	ServerID              string        `json:"serverID"`
	ExternalCalls         []BackendCall `json:"externalCalls"`
	SourceGatewayNode     string        `json:"sourceGatewayNode"`
	CallbackRequest       bool          `json:"callbackRequest"`
}

type HttpHeaders struct {
}

type BackendCall struct {
	ExternalCallType string `json:"externalCallType"`
	ExternalURL      string `json:"externalURL"`
	CallStartTime    int64  `json:"callStartTime"`
	CallEndTime      int64  `json:"callEndTime"`
	CallDuration     int    `json:"callDuration"`
	ResponseCode     string `json:"responseCode"`
}
