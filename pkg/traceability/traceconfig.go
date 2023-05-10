package traceability

import (
	"fmt"
	"os"

	"github.com/Axway/agent-sdk/pkg/cmd/properties"
	corecfg "github.com/Axway/agent-sdk/pkg/config"
)

var agentConfig *AgentConfigTraceability

const (
	pathLogFile = "webmethods.logFile"
)

// SetConfig sets the global AgentConfig reference.
func SetConfig(newConfig *AgentConfigTraceability) {
	agentConfig = newConfig
}

// GetConfig gets the AgentConfig
func GetConfig() *AgentConfigTraceability {
	return agentConfig
}

// AgentConfig - represents the config for agent
type AgentConfigTraceability struct {
	CentralConfig              corecfg.CentralConfig       `config:"central"`
	WebMethodConfigTracability *WebMethodConfigTracability `config:"webmethod"`
}

// WebMethodConfig - represents the config for the Webmethods APIM
type WebMethodConfigTracability struct {
	corecfg.IConfigValidator
	AgentType corecfg.AgentType
	LogFile   string            `config:"logFile"`
	TLS       corecfg.TLSConfig `config:"ssl"`
}

// ValidateCfg - Validates the gateway config
func (c *WebMethodConfigTracability) ValidateCfg() (err error) {

	if c.LogFile == "" {
		return fmt.Errorf("invalid Webmethods APIM configuration: logFile is not configured")
	}

	if c.AgentType == corecfg.TraceabilityAgent && c.LogFile != "" {
		if _, err := os.Stat(c.LogFile); os.IsNotExist(err) {
			return fmt.Errorf("invalid  Webmethods APIM log path: path does not exist: %s", c.LogFile)
		}
	}
	return
}

// AddConfigProperties - Adds the command properties needed for Webmethods agent
func AddConfigProperties(props properties.Properties) {
	props.AddStringProperty(pathLogFile, "./logs/traffic.log", "Sample log file with traffic event from gateway")
}

// NewWebmothodsConfig - parse the props and create an Webmethods Configuration structure
func NewWebmothodsConfig(props properties.Properties, agentType corecfg.AgentType) *WebMethodConfigTracability {
	return &WebMethodConfigTracability{
		AgentType: agentType,
		LogFile:   props.StringPropertyValue(pathLogFile),
	}
}
