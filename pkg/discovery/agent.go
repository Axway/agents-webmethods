package discovery

import (
	"os"
	"os/signal"
	"syscall"

	"git.ecd.axway.org/apigov/agents-webmethods/pkg/config"
	"git.ecd.axway.org/apigov/agents-webmethods/pkg/webmethods"
	coreAgent "github.com/Axway/agent-sdk/pkg/agent"
	"github.com/Axway/agent-sdk/pkg/cache"
)

type Repeater interface {
	Loop()
	OnConfigChange(cfg *config.WebMethodConfig)
	Stop()
}

// Agent -
type Agent struct {
	client    webmethods.Client
	stopAgent chan bool
	discovery Repeater
	publisher Repeater
}

// NewAgent creates a new agent
func NewAgent(cfg *config.AgentConfig, client webmethods.Client) (agent *Agent) {
	buffer := 5
	apiChan := make(chan *ServiceDetail, buffer)

	pub := &publisher{
		apiChan:     apiChan,
		stopPublish: make(chan bool),
		publishAPI:  coreAgent.PublishAPI,
	}

	c := cache.New()

	svcHandler := &serviceHandler{
		client: client,
		cache:  c,
	}

	svcHandler.mode = marketplace

	disc := &discovery{
		apiChan:           apiChan,
		cache:             c,
		client:            client,
		centralClient:     coreAgent.GetCentralClient(),
		discoveryPageSize: 50,
		pollInterval:      cfg.WebMethodConfig.PollInterval,
		stopDiscovery:     make(chan bool),
		serviceHandler:    svcHandler,
	}

	return newAgent(client, disc, pub)
}

func newAgent(
	client webmethods.Client,
	discovery Repeater,
	publisher Repeater,
) *Agent {
	return &Agent{
		client:    client,
		discovery: discovery,
		publisher: publisher,
		stopAgent: make(chan bool),
	}
}

// onConfigChange apply configuration changes
func (a *Agent) onConfigChange() {
	cfg := config.GetConfig()

	// Stop Discovery & Publish
	a.discovery.Stop()
	a.publisher.Stop()

	a.client.OnConfigChange(cfg.WebMethodConfig)
	a.discovery.OnConfigChange(cfg.WebMethodConfig)

	// Restart Discovery & Publish
	go a.discovery.Loop()
	go a.publisher.Loop()
}

// Run the agent loop
func (a *Agent) Run() {
	coreAgent.OnConfigChange(a.onConfigChange)

	go a.discovery.Loop()
	go a.publisher.Loop()

	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, os.Interrupt)

	<-gracefulStop
	a.Stop()
}

// Stop stops the discovery agent.
func (a *Agent) Stop() {
	a.discovery.Stop()
	a.publisher.Stop()
	close(a.stopAgent)
}
