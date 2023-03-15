package discovery

import (
	"github.com/Axway/agent-sdk/pkg/agent"
	"github.com/sirupsen/logrus"

	"git.ecd.axway.org/apigov/agents-webmethods/pkg/boomi"
	subs "git.ecd.axway.org/apigov/agents-webmethods/pkg/subscription"
	corecmd "github.com/Axway/agent-sdk/pkg/cmd"
	"github.com/Axway/agent-sdk/pkg/cmd/service"
	corecfg "github.com/Axway/agent-sdk/pkg/config"

	"git.ecd.axway.org/apigov/agents-webmethods/pkg/config"
	"git.ecd.axway.org/apigov/agents-webmethods/pkg/discovery"
)

// RootCmd - Agent root command
var (
	RootCmd        corecmd.AgentRootCmd
	discoveryAgent *discovery.Agent
)

func init() {
	// Create new root command with callbacks to initialize the agent config and command execution.
	// The first parameter identifies the name of the yaml file that agent will look for to load the config
	RootCmd = corecmd.NewRootCmd(
		"boomi_discovery_agent", // Name of the yaml file
		"Boomi Discovery Agent", // Agent description
		initConfig,              // Callback for initializing the agent config
		run,                     // Callback for executing the agent
		corecfg.DiscoveryAgent,  // Agent Type (Discovery or Traceability)
	)
	config.AddConfigProperties(RootCmd.GetProperties())

	RootCmd.AddCommand(service.GenServiceCmd("pathConfig"))
}

// run Callback that agent will call to process the execution
func run() error {
	discoveryAgent.Run()
	return nil
}

// initConfig Callback that agent will call to initialize the config. CentralConfig is parsed by Agent SDK
// and passed to the callback allowing the agent code to access the central config
func initConfig(centralConfig corecfg.CentralConfig) (interface{}, error) {
	conf := &config.AgentConfig{
		CentralConfig: centralConfig,
		BoomiConfig:   config.NewBoomiConfig(RootCmd.GetProperties()),
	}

	config.SetConfig(conf)

	logger := logrus.WithFields(logrus.Fields{
		"component": "agent",
	})

	gatewayClient := boomi.NewClient(conf.BoomiConfig)
	if centralConfig.IsMarketplaceSubsEnabled() {
		agent.RegisterProvisioner(subs.NewProvisioner(gatewayClient, logger))
		agent.NewAPIKeyAccessRequestBuilder().Register()
		agent.NewOAuthCredentialRequestBuilder(agent.WithCRDOAuthSecret()).IsRenewable().Register()
	}

	discoveryAgent = discovery.NewAgent(conf, gatewayClient)

	return conf, nil
}
