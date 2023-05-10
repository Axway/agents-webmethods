package traceability

import (
	corecmd "github.com/Axway/agent-sdk/pkg/cmd"
	"github.com/Axway/agent-sdk/pkg/cmd/service"
	corecfg "github.com/Axway/agent-sdk/pkg/config"
	"github.com/Axway/agents-webmethods/pkg/config"
	"github.com/Axway/agents-webmethods/pkg/traceability"

	libcmd "github.com/elastic/beats/v7/libbeat/cmd"
	"github.com/elastic/beats/v7/libbeat/cmd/instance"
)

// RootCmd - Agent root command
var RootCmd corecmd.AgentRootCmd
var beatCmd *libcmd.BeatsRootCmd

func init() {
	name := "webmethods_traceability_agent"
	settings := instance.Settings{
		Name:            name,
		HasDashboards:   true,
		ConfigOverrides: corecfg.LogConfigOverrides(),
	}

	// Initialize the beat command
	beatCmd = libcmd.GenRootCmdWithSettings(traceability.NewBeater, settings)
	cmd := beatCmd.Command
	// Wrap the beat command with the agent command processor with callbacks to initialize the agent config and command execution.
	// The first parameter identifies the name of the yaml file that agent will look for to load the config
	RootCmd = corecmd.NewCmd(
		&cmd,
		name, // Name of the agent and yaml config file
		"Webmethods Traceability Agent",
		initConfig,
		run,
		corecfg.TraceabilityAgent,
	)
	config.AddConfigProperties(RootCmd.GetProperties())
	RootCmd.AddCommand(service.GenServiceCmd("pathConfig"))
}

// Callback that agent will call to process the execution
func run() error {
	return beatCmd.Execute()
}

// Callback that agent will call to initialize the config. CentralConfig is parsed by Agent SDK
// and passed to the callback allowing the agent code to access the central config
func initConfig(centralConfig corecfg.CentralConfig) (interface{}, error) {
	agentConfig := &traceability.AgentConfigTraceability{
		CentralConfig:              centralConfig,
		WebMethodConfigTracability: traceability.NewWebmothodsConfig(RootCmd.GetProperties(), centralConfig.GetAgentType()),
	}
	traceability.SetConfig(agentConfig)
	return agentConfig, nil
}
