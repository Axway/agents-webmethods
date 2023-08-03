package traceability

import (
	"fmt"
	hc "github.com/Axway/agent-sdk/pkg/util/healthcheck"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"

	coreagent "github.com/Axway/agent-sdk/pkg/agent"
	coreapi "github.com/Axway/agent-sdk/pkg/api"
	"github.com/Axway/agent-sdk/pkg/transaction"
	"github.com/Axway/agent-sdk/pkg/util/errors"
	agenterrors "github.com/Axway/agent-sdk/pkg/util/errors"
	"github.com/Axway/agents-webmethods/pkg/config"
	localerrors "github.com/Axway/agents-webmethods/pkg/errors"
	"github.com/Axway/agents-webmethods/pkg/webmethods"
)

// Agent - mulesoft Beater configuration. Implements the beat.Beater interface.
type Agent struct {
	client         beat.Client
	doneCh         chan struct{}
	eventChannel   chan WebmethodsEvent
	eventProcessor Processor
	webmethods     Emitter
}

func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {

	agentCfg := config.GetConfig()
	if agentCfg == nil {
		return nil, localerrors.ErrConfigFile
	}
	generator := transaction.NewEventGenerator()
	httpClient := coreapi.NewClient(agentCfg.WebMethodConfig.TLS, agentCfg.WebMethodConfig.ProxyURL)
	client, err := webmethods.NewClient(agentCfg.WebMethodConfig, httpClient)
	if err != nil {
		return nil, err
	}
	webmethodsAPIGatewayTimezone := agentCfg.WebMethodConfig.Timezone
	timezoneLocation, err := time.LoadLocation("UTC")
	if webmethodsAPIGatewayTimezone != "" {
		timezoneLocation, err = time.LoadLocation(webmethodsAPIGatewayTimezone)
		if err != nil {
			return nil, errors.Newf(4001, "invalid timestamp %s")
		}
	}
	processor := NewApiEventProcessor(agentCfg, generator)
	eventChannel := make(chan WebmethodsEvent)
	emitter := NewWebmethodsEventEmitter(agentCfg.WebMethodConfig.CachePath, agentCfg.WebMethodConfig.PollInterval, eventChannel, client, *timezoneLocation)
	emitterJob, err := NewMuleEventEmitterJob(emitter, agentCfg.WebMethodConfig.PollInterval, client)
	if err != nil {
		return nil, err
	}
	return newAgent(processor, emitterJob, eventChannel)
}

func newAgent(
	processor Processor,
	emitter Emitter,
	eventChannel chan WebmethodsEvent,
) (*Agent, error) {
	a := &Agent{
		doneCh:         make(chan struct{}),
		eventChannel:   eventChannel,
		eventProcessor: processor,
		webmethods:     emitter,
	}

	//Validate that all necessary services are up and running. If not, return error
	if hc.RunChecks() != hc.OK {
		return nil, agenterrors.ErrInitServicesNotReady
	}

	return a, nil
}

// Run starts the Mulesoft traceability agent.
func (a *Agent) Run(b *beat.Beat) error {
	coreagent.OnConfigChange(a.onConfigChange)

	var err error
	a.client, err = b.Publisher.Connect()
	if err != nil {
		coreagent.UpdateStatus(coreagent.AgentFailed, err.Error())
		return err
	}

	go func() {
		err := a.webmethods.Start()
		if err != nil {
			fmt.Printf("Unable to start :%s", err)
		}
	}()

	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, os.Interrupt)

	for {
		select {
		case <-a.doneCh:
			return a.client.Close()
		case <-gracefulStop:
			return a.client.Close()
		case event := <-a.eventChannel:
			eventsToPublish := a.eventProcessor.ProcessRaw(event)
			a.client.PublishAll(eventsToPublish)
		}
	}
}

// onConfigChange apply configuration changes
func (a *Agent) onConfigChange() {
	cfg := config.GetConfig()
	a.webmethods.OnConfigChange(cfg)
}

// Stop stops the agent.
func (a *Agent) Stop() {
	a.doneCh <- struct{}{}
}
