package traceability

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Axway/agent-sdk/pkg/agent"
	v1 "github.com/Axway/agent-sdk/pkg/apic/apiserver/models/api/v1"
	"github.com/Axway/agent-sdk/pkg/transaction"
	transutil "github.com/Axway/agent-sdk/pkg/transaction/util"
	"github.com/Axway/agent-sdk/pkg/util"
	"github.com/Axway/agent-sdk/pkg/util/log"
	"github.com/Axway/agents-webmethods/pkg/config"
	"github.com/sirupsen/logrus"

	"github.com/elastic/beats/v7/libbeat/beat"
)

const Client = "Client"
const WebmethodsProxy = "Webmethods.APIProxy"

type Processor interface {
	ProcessRaw(webmethodsEvent WebmethodsEvent) []beat.Event
}

type cacheManager interface {
	GetManagedApplicationCacheKeys() []string
	GetManagedApplication(id string) *v1.ResourceInstance
}

// EventProcessor - represents the processor for received event for Amplify Central
// The event processing can be done either when the beat input receives the log entry or before the beat transport
// publishes the event to transport.
// When processing the received log entry on input, the log entry is mapped to structure expected for Amplify Central Observer
// and then beat.Event is published to beat output that produces the event over the configured transport.
// When processing the log entry on output, the log entry is published to output as beat.Event. The output transport invokes
// the Process(events []publisher.Event) method which is set as output event processor. The Process() method processes the received
// log entry and performs the mapping to structure expected for Amplify Central Observer. The method returns the converted Events to
// transport publisher which then produces the events over the transport.
type ApiEventProcessor struct {
	cfg            *config.AgentConfig
	eventGenerator transaction.EventGenerator
	cacheManager   cacheManager
	appIDToManApp  map[string]string
}

func NewApiEventProcessor(
	gateway *config.AgentConfig,
	eventGenerator transaction.EventGenerator,
) *ApiEventProcessor {
	ep := &ApiEventProcessor{
		cfg:            gateway,
		eventGenerator: eventGenerator,
		cacheManager:   agent.GetCacheManager(),
		appIDToManApp:  make(map[string]string),
	}
	return ep
}

// ProcessRaw - process the received log entry and returns the event to be published to Amplifyingestion service
func (aep *ApiEventProcessor) ProcessRaw(webmethodsEvent WebmethodsEvent) []beat.Event {

	logrus.Infof("%+v\n", webmethodsEvent)
	summaryEvent, logEvents, err := aep.processMapping(webmethodsEvent)
	if err != nil {
		logrus.Error(err.Error())
		return nil
	}
	// Generates the beat.Event with attributes by Amplify ingestion service
	events, err := aep.eventGenerator.CreateEvents(*summaryEvent, logEvents, time.Now(), nil, nil, nil)
	if err != nil {
		logrus.Error(err.Error())
		return nil
	}
	return events
	//return nil

}

func (aep *ApiEventProcessor) processMapping(webmethodsEvent WebmethodsEvent) (*transaction.LogEvent, []transaction.LogEvent, error) {
	centralCfg := agent.GetCentralConfig()

	eventTimestamp := time.UnixMilli(webmethodsEvent.CreationDate)
	eventTime := eventTimestamp.UTC().UnixNano() / int64(time.Millisecond)
	//eventTime := time.Now().UTC().Format(gatewayTrafficLogEntry.EventTimestamp)
	txID := webmethodsEvent.SessionId
	txEventID := webmethodsEvent.CorrelationID
	leg0ID := FormatLeg0(txEventID)
	leg1ID := FormatLeg1(txEventID)

	transInboundLogEventLeg, err := aep.createTransactionEvent(eventTime, txID, webmethodsEvent, leg0ID, "", "Inbound")
	if err != nil {
		return nil, nil, err

	}
	transOutboundLogEventLeg, err := aep.createOutboundTransactionEvent(txID, webmethodsEvent, leg1ID, leg0ID, "Outbound")
	if err != nil {
		return nil, nil, err
	}

	transSummaryLogEvent, err := aep.createSummaryEvent(eventTime, txID, webmethodsEvent, centralCfg.GetTeamID())
	if err != nil {
		return nil, nil, err
	}

	logrus.Infof("outbound leg %v", transOutboundLogEventLeg)
	if transOutboundLogEventLeg == nil {
		return transSummaryLogEvent, []transaction.LogEvent{
			*transInboundLogEventLeg,
		}, nil
	}

	return transSummaryLogEvent, []transaction.LogEvent{
		*transOutboundLogEventLeg,
		*transInboundLogEventLeg,
	}, nil
}

func (aep *ApiEventProcessor) getTransactionEventStatus(code int) transaction.TxEventStatus {
	if code >= 400 {
		return transaction.TxEventStatusFail
	}
	return transaction.TxEventStatusPass
}

func (aep *ApiEventProcessor) getTransactionSummaryStatus(statusCode int) transaction.TxSummaryStatus {
	transSummaryStatus := transaction.TxSummaryStatusUnknown
	if statusCode >= http.StatusOK && statusCode < http.StatusBadRequest {
		transSummaryStatus = transaction.TxSummaryStatusSuccess
	} else if statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError {
		transSummaryStatus = transaction.TxSummaryStatusFailure
	} else if statusCode >= http.StatusInternalServerError && statusCode < http.StatusNetworkAuthenticationRequired {
		transSummaryStatus = transaction.TxSummaryStatusException
	}
	return transSummaryStatus
}

func (aep *ApiEventProcessor) createTransactionEvent(eventTime int64, txID string, webmethodsEvent WebmethodsEvent, eventID, parentId, direction string) (*transaction.LogEvent, error) {

	req := map[string]string{
		"User-AgentName": "test",
		"Request-ID":     eventID,
	}
	res := map[string]string{

		"Response-Time": "20",
	}

	httpStatus, _ := strconv.Atoi(webmethodsEvent.ResponseCode)
	host := webmethodsEvent.ServerID
	port := 443
	if strings.Index(host, ":") != -1 {
		uris := strings.Split(host, ":")
		host = uris[0]
		port, _ = strconv.Atoi(uris[1])
	}
	httpProtocolDetails, err := transaction.NewHTTPProtocolBuilder().
		SetByteLength(100, 100).
		SetRemoteAddress("test", "localhost", 443).
		SetURI(webmethodsEvent.OperationName).
		SetMethod(webmethodsEvent.HTTPMethod).
		SetHeaders(buildHeaders(req), buildHeaders(res)).
		SetStatus(httpStatus, http.StatusText(httpStatus)).
		SetHost(host).
		SetLocalAddress(host, port).
		Build()
	if err != nil {
		return nil, err
	}

	return transaction.NewTransactionEventBuilder().
		SetTimestamp(eventTime).
		SetTransactionID(txID).
		SetID(eventID).
		SetParentID(parentId).
		SetSource(WebmethodsProxy).
		SetDestination(webmethodsEvent.SourceGatewayNode).
		SetDirection(direction).
		SetStatus(aep.getTransactionEventStatus(httpStatus)).
		SetProtocolDetail(httpProtocolDetails).
		Build()
}

func (aep *ApiEventProcessor) createOutboundTransactionEvent(txID string, webmethodsEvent WebmethodsEvent, eventID, parentId, direction string) (*transaction.LogEvent, error) {

	req := map[string]string{
		"User-AgentName": "test",
		"Request-ID":     eventID,
	}
	res := map[string]string{

		"Response-Time": "20",
	}
	backendCall := webmethodsEvent.ExternalCalls
	log.Infof("BAckend call - %v", backendCall)
	if len(backendCall) > 0 {
		logrus.Info("Processing outbound calls")
		httpStatus, _ := strconv.Atoi(backendCall[0].ResponseCode)
		httpUrl, err := url.Parse(backendCall[0].ExternalURL)
		if err != nil {
			logrus.Error("Unable to parse URL", err)
			return nil, nil
		}
		var host, port string
		scheme := httpUrl.Scheme
		if strings.Index(httpUrl.Host, ":") != -1 {
			host, port, _ = net.SplitHostPort(httpUrl.Host)
		} else {
			host = httpUrl.Host
			if scheme == "https" {
				port = "443"
			} else {
				port = "80"
			}
		}
		portInt, _ := strconv.Atoi(port)
		httpProtocolDetails, err := transaction.NewHTTPProtocolBuilder().
			SetByteLength(100, 100).
			SetRemoteAddress("test", "localhost", 443).
			SetURI(webmethodsEvent.OperationName).
			SetMethod(webmethodsEvent.HTTPMethod).
			SetHeaders(buildHeaders(req), buildHeaders(res)).
			SetStatus(httpStatus, http.StatusText(httpStatus)).
			SetHost(host).
			SetLocalAddress(host, portInt).
			Build()
		if err != nil {
			return nil, err
		}
		eventTimestamp := time.UnixMilli(backendCall[0].CallStartTime)
		eventTime := eventTimestamp.UTC().UnixNano() / int64(time.Millisecond)
		logrus.Info("Outbound object created")
		return transaction.NewTransactionEventBuilder().
			SetTimestamp(eventTime).
			SetTransactionID(txID).
			SetID(eventID).
			SetParentID(parentId).
			SetSource(Client).
			SetDestination(WebmethodsProxy).
			SetDirection(direction).
			SetDuration(backendCall[0].CallDuration).
			SetStatus(aep.getTransactionEventStatus(httpStatus)).
			SetProtocolDetail(httpProtocolDetails).
			Build()
	}
	return nil, nil
}

func (aep *ApiEventProcessor) createSummaryEvent(eventTime int64, txID string, webmethodsEvent WebmethodsEvent, teamID string) (*transaction.LogEvent, error) {
	statusCode, _ := strconv.Atoi(webmethodsEvent.ResponseCode)
	method := webmethodsEvent.HTTPMethod
	uri := webmethodsEvent.OperationName
	host := webmethodsEvent.ApplicationIp

	builder := transaction.NewTransactionSummaryBuilder().
		SetTimestamp(eventTime).
		SetTransactionID(txID).
		SetStatus(aep.getTransactionSummaryStatus(statusCode), strconv.Itoa(statusCode)).
		SetTeam(teamID).
		SetDuration(webmethodsEvent.TotalTime).
		SetEntryPoint("http", method, uri, host).
		// If the API is published to Central as unified catalog item/API service, se the Proxy details with the API definition
		// The Proxy.Name represents the name of the API
		// The Proxy.ID should be of format "remoteApiId_<ID Of the API on remote gateway>". Use transaction.FormatProxyID(<ID Of the API on remote gateway>) to get the formatted value.
		SetProxy(transutil.FormatProxyID(webmethodsEvent.ApiId), webmethodsEvent.ApiName, 0)

	if webmethodsEvent.ApplicationName != "Unknown" && webmethodsEvent.ApplicationId != "Unknown" {
		builder.SetApplication(transutil.FormatApplicationID(webmethodsEvent.ApplicationId), webmethodsEvent.ApplicationName)
	}
	return builder.Build()

}

func (aep *ApiEventProcessor) getApplicationByName(appId string, appName string) (string, string) {
	// find the manged application in the cache
	manAppName := aep.getManagedApplicationNameByID(appId)
	if manAppName != "" {
		return appId, manAppName
	}
	return appId, appName

}

func (aep *ApiEventProcessor) getManagedApplicationNameByID(appID string) string {
	if name, ok := aep.appIDToManApp[appID]; ok {
		return name
	}
	for _, key := range aep.cacheManager.GetManagedApplicationCacheKeys() {
		ri := aep.cacheManager.GetManagedApplication(key)
		val, _ := util.GetAgentDetailsValue(ri, "webmethodsApplicationId")
		if val == appID {
			aep.appIDToManApp[appID] = ri.Name
			return ri.Name
		}
	}
	return ""
}

func FormatLeg0(id string) string {
	return fmt.Sprintf("%s-leg0", id)
}

func FormatLeg1(id string) string {
	return fmt.Sprintf("%s-leg1", id)
}

func buildHeaders(headers map[string]string) string {
	jsonHeader, err := json.Marshal(headers)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonHeader)
}
