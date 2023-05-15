package traceability

import (
	"os"
	"os/signal"
	"syscall"

	coreagent "github.com/Axway/agent-sdk/pkg/agent"
	"github.com/Axway/agent-sdk/pkg/transaction"
	agenterrors "github.com/Axway/agent-sdk/pkg/util/errors"
	hc "github.com/Axway/agent-sdk/pkg/util/healthcheck"
	"github.com/Axway/agents-webmethods/pkg/config"
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
)

// Agent - Webmethods Beater configuration. Implements the beat.Beater interface.
type Agent struct {
	client         beat.Client
	doneCh         chan struct{}
	eventChannel   chan string
	eventProcessor Processor
	webmethods     Emitter
}

// NewBeater creates an instance of webmethods_traceability_agent.
func NewBeater(_ *beat.Beat, _ *common.Config) (beat.Beater, error) {
	eventChannel := make(chan string)
	agentConfig := config.GetConfig()

	generator := transaction.NewEventGenerator()
	mapper := &EventMapper{}
	processor := NewEventProcessor(agentConfig, generator, mapper)
	emitter := NewWebmethodsEventEmitter(agentConfig.WebMethodConfig.LogFile, eventChannel)

	return newAgent(processor, emitter, eventChannel)
}

func newAgent(
	processor Processor,
	emitter Emitter,
	eventChannel chan string,
) (*Agent, error) {
	a := &Agent{
		doneCh:         make(chan struct{}),
		eventChannel:   eventChannel,
		eventProcessor: processor,
		webmethods:     emitter,
	}

	// Validate that all necessary services are up and running. If not, return error
	if hc.RunChecks() != hc.OK {
		return nil, agenterrors.ErrInitServicesNotReady
	}

	return a, nil
}

// Run starts the webmethods traceability agent.
func (a *Agent) Run(b *beat.Beat) error {
	coreagent.OnConfigChange(a.onConfigChange)

	var err error
	a.client, err = b.Publisher.Connect()
	if err != nil {
		coreagent.UpdateStatus(coreagent.AgentFailed, err.Error())
		return err
	}

	go a.webmethods.Start()

	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, os.Interrupt)

	for {
		select {
		case <-a.doneCh:
			return a.client.Close()
		case <-gracefulStop:
			return a.client.Close()
		case event := <-a.eventChannel:
			eventsToPublish := a.eventProcessor.ProcessRaw([]byte(event))
			a.client.PublishAll(eventsToPublish)
		}
	}
}

// onConfigChange apply configuration changes
func (a *Agent) onConfigChange() {
}

// Stop stops the agent.
func (a *Agent) Stop() {
	a.doneCh <- struct{}{}
}
