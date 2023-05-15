package traceability

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/beats/v7/libbeat/publisher"

	"github.com/Axway/agent-sdk/pkg/agent"
	v1 "github.com/Axway/agent-sdk/pkg/apic/apiserver/models/api/v1"
	"github.com/Axway/agent-sdk/pkg/transaction"
	"github.com/Axway/agent-sdk/pkg/util"
	"github.com/Axway/agent-sdk/pkg/util/log"

	transutil "github.com/Axway/agent-sdk/pkg/transaction/util"
	sdkerrors "github.com/Axway/agent-sdk/pkg/util/errors"
	agenterrors "github.com/Axway/agents-webmethods/pkg/errors"
)

const (
	condorKey  = "condor"
	retriesKey = "retries"
)

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
type EventProcessor struct {
	cfg            *AgentConfigTraceability
	eventGenerator transaction.EventGenerator
	tenantID       string
	cacheManager   cacheManager
	appIDToManApp  map[string]string
}

// func NewEventProcessor(
// 	gateway *AgentConfigTraceability,
// 	eventGenerator transaction.EventGenerator,
// 	mapper *EventMapper,
// ) *EventProcessor {
// 	ep := &EventProcessor{
// 		cfg:            gateway,
// 		eventGenerator: eventGenerator,
// 		eventMapper:    mapper,
// 		cacheManager:   agent.GetCacheManager(),
// 		appIDToManApp:  make(map[string]string),
// 		//tenantID:       agentConfig.Central.GetTenantID(),
// 	}
// 	return ep
// }

func NewEventProcessor(
	gateway *AgentConfigTraceability,
	maxRetries int,
) *EventProcessor {
	ep := &EventProcessor{
		cfg:            gateway,
		eventGenerator: transaction.NewEventGenerator(),
		cacheManager:   agent.GetCacheManager(),
		appIDToManApp:  make(map[string]string),
		//tenantID:       agentConfig.Central.GetTenantID(),
	}
	log.Trace("Event Processor Created")
	return ep
}

// Process - Process the log file, waiting for events
func (p *EventProcessor) Process(events []publisher.Event) []publisher.Event {
	newEvents := make([]publisher.Event, 0)
	for _, event := range events {
		newEvents, _ = p.ProcessEvent(newEvents, event)
	}
	for _, newEvent := range newEvents {
		str := fmt.Sprintf("%v", newEvent)
		log.Trace("New event message to process - " + str)
	}
	return newEvents
}

func (p *EventProcessor) ProcessEvent(newEvents []publisher.Event, event publisher.Event) ([]publisher.Event, error) {
	eventMsgFieldVal, err := event.Content.Fields.GetValue("message")
	if err != nil {
		return newEvents, sdkerrors.Wrap(agenterrors.ErrEventNoMsg, err.Error()).FormatError(event)
	}
	eventMsg, ok := eventMsgFieldVal.(string)
	if !ok {
		return newEvents, nil
	}
	data := strings.Split(string(eventMsg), "|")
	if len(data) == 39 {
		totalTime, _ := strconv.Atoi(data[17])
		nativeTime, _ := strconv.Atoi(data[18])
		gatewayTrafficLogEntry := GwTrafficLogEntry{
			EventType:             data[0],
			RootContext:           data[1],
			ParentContext:         data[2],
			CurrentContext:        data[3],
			Uuid:                  data[4],
			ServerId:              data[5],
			EventTimestamp:        data[6],
			CurrentTimestamp:      data[7],
			SessionId:             data[8],
			ApiName:               data[9],
			ApiVersion:            data[10],
			TargetName:            data[11],
			ApplicationName:       data[12],
			ApplicationIp:         data[13],
			ApplicationId:         data[14],
			Request:               data[15],
			Response:              data[16],
			TotalTime:             totalTime,
			NativeTime:            nativeTime,
			RequestStatus:         data[19],
			OperationName:         data[20],
			NatvieEndpoint:        data[21],
			PartnerId:             data[22],
			ApiId:                 data[23],
			ServiceName:           data[24],
			RequestHeaders:        data[25],
			QueryParam:            data[26],
			ResponseHeaders:       data[27],
			CorrelationId:         data[28],
			ErrorOrigin:           data[29],
			Custom:                data[30],
			NativeRequestHeaders:  data[31],
			NativeRequestPayload:  data[32],
			NativeResponseHeaders: data[33],
			NativeResponsePayload: data[34],
			NativeHttpMethod:      data[35],
			NativeUrl:             data[36],
			ExternalCalls:         data[37],
			SourceGatewayNode:     data[38],
		}
		summaryEvent, logEvents, err := p.processMapping(gatewayTrafficLogEntry)
		if err != nil {
			log.Error(err.Error())
			return newEvents, nil
		}

		//events, err := p.eventGenerator.CreateEvents(*summaryEvent, logEvents, time.Now(), nil, nil, nil)
		events, err := p.createCondorEvents(event, *summaryEvent, logEvents)
		if err != nil {
			return newEvents, nil
		}
		newEvents = append(newEvents, events...)

	} else {
		log.Errorf("Invalid record %s", eventMsg)
	}
	return newEvents, nil

}

func (p *EventProcessor) getApplicationByName(appId string, appName string) (string, string) {
	// Add the V7 Application ID, with prefix, and Name to the event

	// find the manged application in the cache
	manAppName := p.getManagedApplicationNameByID(appId)
	if manAppName != "" {
		return appId, manAppName
	}
	return appId, appName

}

func (p *EventProcessor) getManagedApplicationNameByID(appID string) string {
	if name, ok := p.appIDToManApp[appID]; ok {
		return name
	}
	for _, key := range p.cacheManager.GetManagedApplicationCacheKeys() {
		ri := p.cacheManager.GetManagedApplication(key)
		val, _ := util.GetAgentDetailsValue(ri, "webmethodsApplicationId")
		if val == appID {
			p.appIDToManApp[appID] = ri.Name
			return ri.Name
		}
	}
	return ""
}

func (p *EventProcessor) createCondorEvents(originalLogEvent publisher.Event, summaryEvent transaction.LogEvent, detailEvents []transaction.LogEvent) ([]publisher.Event, error) {
	// Create the beat event then wrap that in the publisher event for Condor

	// Add a Retry count to Meta
	if originalLogEvent.Content.Meta == nil {
		originalLogEvent.Content.Meta = make(map[string]interface{})
	}
	originalLogEvent.Content.Meta[retriesKey] = 0
	originalLogEvent.Content.Meta[condorKey] = true

	beatEvents, err := p.eventGenerator.CreateEvents(summaryEvent, detailEvents, originalLogEvent.Content.Timestamp, originalLogEvent.Content.Meta, originalLogEvent.Content.Fields, originalLogEvent.Content.Private)

	if err != nil {
		return nil, sdkerrors.Wrap(agenterrors.ErrCreateCondorEvent, err.Error())
	}

	events := make([]publisher.Event, 0)
	for _, beatEvent := range beatEvents {
		events = append(events, publisher.Event{
			Content: beatEvent,
			Flags:   originalLogEvent.Flags,
		})
	}
	return events, nil
}

func (ep *EventProcessor) processMapping(gatewayTrafficLogEntry GwTrafficLogEntry) (*transaction.LogEvent, []transaction.LogEvent, error) {
	centralCfg := agent.GetCentralConfig()

	eventTimestamp, _ := time.Parse(time.RFC3339, gatewayTrafficLogEntry.EventTimestamp)
	eventTime := eventTimestamp.UTC().UnixNano() / int64(time.Millisecond)
	//eventTime := time.Now().UTC().Format(gatewayTrafficLogEntry.EventTimestamp)
	txID := gatewayTrafficLogEntry.Uuid
	txEventID := gatewayTrafficLogEntry.CorrelationId
	transInboundLogEventLeg, err := ep.createTransactionEvent(eventTime, txID, gatewayTrafficLogEntry, txEventID, "Inbound")
	if err != nil {
		return nil, nil, err
	}

	transSummaryLogEvent, err := ep.createSummaryEvent(eventTime, txID, gatewayTrafficLogEntry, centralCfg.GetTeamID())
	if err != nil {
		return nil, nil, err
	}

	return transSummaryLogEvent, []transaction.LogEvent{
		*transInboundLogEventLeg,
	}, nil
}

func (ep *EventProcessor) getTransactionEventStatus(code string) transaction.TxEventStatus {
	if code != "SUCCESS" {
		return transaction.TxEventStatusFail
	}
	return transaction.TxEventStatusPass
}

func (ep *EventProcessor) getHttpStatusCode(code string) int {
	if code == "SUCCESS" {
		return 200
	}
	return 500
}

func (ep *EventProcessor) getTransactionSummaryStatus(statusCode int) transaction.TxSummaryStatus {
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

func (ep *EventProcessor) createTransactionEvent(eventTime int64, txID string, txDetails GwTrafficLogEntry, eventID, direction string) (*transaction.LogEvent, error) {

	httpStatus := ep.getHttpStatusCode(txDetails.RequestStatus)

	host := txDetails.ServerId
	port := 443
	if strings.Index(host, ":") != -1 {
		uris := strings.Split(host, ":")
		host = uris[0]
		port, _ = strconv.Atoi(uris[1])

	}

	httpProtocolDetails, err := transaction.NewHTTPProtocolBuilder().
		SetURI(txDetails.OperationName).
		SetMethod(txDetails.NativeHttpMethod).
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
		SetSource(txDetails.ServerId).
		SetDirection(direction).
		SetStatus(ep.getTransactionEventStatus(txDetails.RequestStatus)).
		SetProtocolDetail(httpProtocolDetails).
		Build()
}

func (ep *EventProcessor) createSummaryEvent(eventTime int64, txID string, gatewayTrafficLogEntry GwTrafficLogEntry, teamID string) (*transaction.LogEvent, error) {
	statusCode := ep.getHttpStatusCode(gatewayTrafficLogEntry.RequestStatus)
	method := gatewayTrafficLogEntry.NativeHttpMethod
	uri := gatewayTrafficLogEntry.OperationName
	host := gatewayTrafficLogEntry.ApplicationIp

	builder := transaction.NewTransactionSummaryBuilder().
		SetTimestamp(eventTime).
		SetTransactionID(txID).
		SetStatus(ep.getTransactionSummaryStatus(statusCode), strconv.Itoa(statusCode)).
		SetTeam(teamID).
		SetDuration(gatewayTrafficLogEntry.TotalTime).
		SetEntryPoint("http", method, uri, host).
		// If the API is published to Central as unified catalog item/API service, se the Proxy details with the API definition
		// The Proxy.Name represents the name of the API
		// The Proxy.ID should be of format "remoteApiId_<ID Of the API on remote gateway>". Use transaction.FormatProxyID(<ID Of the API on remote gateway>) to get the formatted value.
		SetProxy(transutil.FormatProxyID(gatewayTrafficLogEntry.ApiId), gatewayTrafficLogEntry.ApiName, 0)

	if gatewayTrafficLogEntry.ApplicationName != "Unknown" && gatewayTrafficLogEntry.ApplicationId != "Unknown" {
		builder.SetApplication(gatewayTrafficLogEntry.ApplicationId, gatewayTrafficLogEntry.ApplicationName)
	}
	return builder.Build()
}
