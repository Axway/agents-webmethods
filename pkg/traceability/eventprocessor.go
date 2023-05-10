package traceability

import (
	"strconv"
	"strings"
	"time"

	"github.com/Axway/agent-sdk/pkg/transaction"

	"github.com/elastic/beats/v7/libbeat/beat"

	"github.com/Axway/agent-sdk/pkg/util/log"
)

type Processor interface {
	ProcessRaw(rawEvent []byte) []beat.Event
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
	eventMapper    *EventMapper
}

func NewEventProcessor(
	gateway *AgentConfigTraceability,
	eventGenerator transaction.EventGenerator,
	mapper *EventMapper,
) *EventProcessor {
	ep := &EventProcessor{
		cfg:            gateway,
		eventGenerator: eventGenerator,
		eventMapper:    mapper,
	}
	return ep
}

// ProcessRaw - process the received log entry and returns the event to be published to Amplifyingestion service
func (ep *EventProcessor) ProcessRaw(rawEvent []byte) []beat.Event {
	var gatewayTrafficLogEntry GwTrafficLogEntry
	data := strings.Split(string(rawEvent), "|")
	if len(data) == 39 && data[0] == "#AGW_EVENT_TXN" {
		totalTime, _ := strconv.Atoi(data[17])
		nativeTime, _ := strconv.Atoi(data[18])
		gatewayTrafficLogEntry = GwTrafficLogEntry{
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
	} else {
		log.Errorf("Invalid record %s", string(rawEvent))
		return nil
	}

	//Map the log entry to log event structure expected by AmplifyCentral Observer
	summaryEvent, logEvents, err := ep.eventMapper.processMapping(gatewayTrafficLogEntry)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	events, err := ep.eventGenerator.CreateEvents(*summaryEvent, logEvents, time.Now(), nil, nil, nil)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return events
}
