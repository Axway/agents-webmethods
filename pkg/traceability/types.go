package traceability

// Headers - Type for request/response headers
type Headers map[string]string

// GwTransaction - Type for gateway transaction detail
type GwTransaction struct {
	ID              string  `json:"id"`
	SourceHost      string  `json:"srcHost"`
	SourcePort      int     `json:"srcPort"`
	DesHost         string  `json:"destHost"`
	DestPort        int     `json:"destPort"`
	URI             string  `json:"uri"`
	Method          string  `json:"method"`
	StatusCode      int     `json:"statusCode"`
	RequestHeaders  Headers `json:"requestHeaders"`
	ResponseHeaders Headers `json:"responseHeaders"`
	RequestBytes    int     `json:"requestByte"`
	ResponseBytes   int     `json:"responseByte"`
}

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
	TraceID               string        `json:"traceId"`
	APIName               string        `json:"apiName"`
	InboundTransaction    GwTransaction `json:"inbound"`
	OutboundTransaction   GwTransaction `json:"outbound"`
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
