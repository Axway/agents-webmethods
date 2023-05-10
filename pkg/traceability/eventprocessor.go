package traceability

import (
	"fmt"

	"github.com/Axway/agent-sdk/pkg/transaction"

	"github.com/elastic/beats/v7/libbeat/beat"
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
	// var gatewayTrafficLogEntry GwTrafficLogEntry
	// err := json.Unmarshal(rawEvent, &gatewayTrafficLogEntry)
	// if err != nil {
	// 	log.Error(err.Error())
	// 	return nil
	// }
	fmt.Println(string(rawEvent))
	// data := strings.Split(string(rawEvent), "|")
	// fmt.Println(len(data))
	// Map the log entry to log event structure expected by AmplifyCentral Observer
	// summaryEvent, logEvents, err := ep.eventMapper.processMapping(gatewayTrafficLogEntry)
	// if err != nil {
	// 	log.Error(err.Error())
	// 	return nil
	// }

	// events, err := ep.eventGenerator.CreateEvents(*summaryEvent, logEvents, time.Now(), nil, nil, nil)
	// if err != nil {
	// 	log.Error(err.Error())
	// 	return nil
	// }
	// return events
	return nil
}
