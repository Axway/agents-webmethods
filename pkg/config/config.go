package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/Axway/agent-sdk/pkg/cmd/properties"
	corecfg "github.com/Axway/agent-sdk/pkg/config"
)

var config *AgentConfig

const (
	pathPollInterval   = "webmethods.pollInterval"
	pathLogFile        = "webmethods.logFile"
	pathProcessOnInput = "webmethods.processOnInput"

	pathWebmethodsApimUrl = "webmethods.url"
	pathAuthUsername      = "webmethods.auth.username"
	pathAuthPassword      = "webmethods.auth.password"

	pathSSLNextProtos         = "webmethods.ssl.nextProtos"
	pathSSLInsecureSkipVerify = "webmethods.ssl.insecureSkipVerify"
	pathSSLCipherSuites       = "webmethods.ssl.cipherSuites"
	pathSSLMinVersion         = "webmethods.ssl.minVersion"
	pathSSLMaxVersion         = "webmethods.ssl.maxVersion"
	pathProxyURL              = "webmethods.proxyUrl"
	pathCachePath             = "webmethods.cachePath"
)

// SetConfig sets the global AgentConfig reference.
func SetConfig(newConfig *AgentConfig) {
	config = newConfig
}

// GetConfig gets the AgentConfig
func GetConfig() *AgentConfig {
	return config
}

// AgentConfig - represents the config for agent
type AgentConfig struct {
	CentralConfig   corecfg.CentralConfig `config:"central"`
	WebMethodConfig *WebMethodConfig      `config:"webmethod"`
}

// WebMethodConfig - represents the config for the Boomi gateway
type WebMethodConfig struct {
	corecfg.IConfigValidator
	AgentType         corecfg.AgentType
	PollInterval      time.Duration     `config:"pollInterval"`
	LogFile           string            `config:"logFile"`
	ProcessOnInput    bool              `config:"processOnInput"`
	CachePath         string            `config:"cachePath"`
	WebmethodsApimUrl string            `config:"url"`
	Environment       string            `config:"environment"`
	Username          string            `config:"auth.username"`
	Password          string            `config:"auth.password"`
	ProxyURL          string            `config:"proxyUrl"`
	TLS               corecfg.TLSConfig `config:"ssl"`
}

// ValidateCfg - Validates the gateway config
func (c *WebMethodConfig) ValidateCfg() (err error) {
	if c.WebmethodsApimUrl == "" {
		return fmt.Errorf("invalid boomi configuration: atomSphereUrl is not configured")
	} else {
		_, err := url.ParseRequestURI(c.WebmethodsApimUrl)
		if err != nil {
			return fmt.Errorf("invalid boomi Platform URL: %s", c.WebmethodsApimUrl)
		}
	}

	if c.Username == "" {
		return fmt.Errorf("invalid boomi configuration: username is not configured")
	}

	if c.Password == "" {
		return fmt.Errorf("invalid boomi configuration: password is not configured")
	}

	if c.PollInterval == 0 {
		return errors.New("invalid boomi configuration: pollInterval is invalid")
	}

	if _, err := os.Stat(c.CachePath); os.IsNotExist(err) {
		return fmt.Errorf("invalid boomi cache path: path does not exist: %s", c.CachePath)
	}
	c.CachePath = filepath.Clean(c.CachePath)

	if c.AgentType == corecfg.TraceabilityAgent && c.LogFile != "" {
		if _, err := os.Stat(c.LogFile); os.IsNotExist(err) {
			return fmt.Errorf("invalid boomi log path: path does not exist: %s", c.LogFile)
		}
	}
	return
}

// AddConfigProperties - Adds the command properties needed for Boomi agent
func AddConfigProperties(props properties.Properties) {
	props.AddDurationProperty(pathPollInterval, 30*time.Second, "Poll interval for read spec discovery/traffic log")
	props.AddStringProperty(pathLogFile, "./logs/traffic.log", "Sample log file with traffic event from gateway")
	props.AddBoolProperty(pathProcessOnInput, true, "Flag to process received event on input or by output before publishing the event by transport")
	props.AddStringProperty(pathWebmethodsApimUrl, "https://api.webmethod.com", "Webmethods APIM URL.")
	props.AddStringProperty(pathAuthUsername, "", "Webmethods APIM username.")
	props.AddStringProperty(pathAuthPassword, "", "Webmethods APIM password.")
	props.AddStringProperty(pathCachePath, "/tmp", "Webmethods Cache Path")
	// ssl properties and command flags
	props.AddStringSliceProperty(pathSSLNextProtos, []string{}, "List of supported application level protocols, comma separated.")
	props.AddBoolProperty(pathSSLInsecureSkipVerify, false, "Controls whether a client verifies the server's certificate chain and host name.")
	props.AddStringSliceProperty(pathSSLCipherSuites, corecfg.TLSDefaultCipherSuitesStringSlice(), "List of supported cipher suites, comma separated.")
	props.AddStringProperty(pathSSLMinVersion, corecfg.TLSDefaultMinVersionString(), "Minimum acceptable SSL/TLS protocol version.")
	props.AddStringProperty(pathSSLMaxVersion, "0", "Maximum acceptable SSL/TLS protocol version.")
}

// NewBoomiConfig - parse the props and create an Boomi Configuration structure
func NewBoomiConfig(props properties.Properties, agentType corecfg.AgentType) *WebMethodConfig {
	return &WebMethodConfig{
		AgentType:         agentType,
		PollInterval:      props.DurationPropertyValue(pathPollInterval),
		LogFile:           props.StringPropertyValue(pathLogFile),
		ProcessOnInput:    props.BoolPropertyValue(pathProcessOnInput),
		WebmethodsApimUrl: props.StringPropertyValue(pathWebmethodsApimUrl),
		CachePath:         props.StringPropertyValue(pathCachePath),
		Password:          props.StringPropertyValue(pathAuthPassword),
		ProxyURL:          props.StringPropertyValue(pathProxyURL),
		Username:          props.StringPropertyValue(pathAuthUsername),
		TLS: &corecfg.TLSConfiguration{
			NextProtos:         props.StringSlicePropertyValue(pathSSLNextProtos),
			InsecureSkipVerify: props.BoolPropertyValue(pathSSLInsecureSkipVerify),
			CipherSuites:       corecfg.NewCipherArray(props.StringSlicePropertyValue(pathSSLCipherSuites)),
			MinVersion:         corecfg.TLSVersionAsValue(props.StringPropertyValue(pathSSLMinVersion)),
			MaxVersion:         corecfg.TLSVersionAsValue(props.StringPropertyValue(pathSSLMaxVersion)),
		},
	}
}
