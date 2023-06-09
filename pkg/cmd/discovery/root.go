package discovery

import (
	"errors"

	"github.com/Axway/agent-sdk/pkg/agent"
	"github.com/Axway/agent-sdk/pkg/apic/provisioning"
	"github.com/Axway/agent-sdk/pkg/util/log"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"

	coreagent "github.com/Axway/agent-sdk/pkg/agent"
	coreapi "github.com/Axway/agent-sdk/pkg/api"
	corecmd "github.com/Axway/agent-sdk/pkg/cmd"
	"github.com/Axway/agents-webmethods/pkg/subscription"
	subs "github.com/Axway/agents-webmethods/pkg/subscription"
	"github.com/Axway/agents-webmethods/pkg/webmethods"

	"github.com/Axway/agent-sdk/pkg/cmd/service"
	corecfg "github.com/Axway/agent-sdk/pkg/config"

	"github.com/Axway/agents-webmethods/pkg/config"
	"github.com/Axway/agents-webmethods/pkg/discovery"
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
		"webmethods_discovery_agent", // Name of the yaml file
		"Webmethods Discovery Agent", // Agent description
		initConfig,                   // Callback for initializing the agent config
		run,                          // Callback for executing the agent
		corecfg.DiscoveryAgent,       // Agent Type (Discovery or Traceability)
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
		CentralConfig:   centralConfig,
		WebMethodConfig: config.NewWebmothodsConfig(RootCmd.GetProperties(), centralConfig.GetAgentType()),
	}

	config.SetConfig(conf)

	logger := logrus.WithFields(logrus.Fields{
		"component": "agent",
	})
	client := coreapi.NewClient(conf.WebMethodConfig.TLS, conf.WebMethodConfig.ProxyURL)
	gatewayClient, err := webmethods.NewClient(conf.WebMethodConfig, client)
	if err != nil {
		return nil, err
	}

	oauthServersResponse, err := gatewayClient.ListOauth2Servers()

	if err != nil {
		return nil, err
	}
	servers := []string{}

	for _, server := range oauthServersResponse.Alias {
		servers = append(servers, server.Name)
	}

	// filter oauth authz server based on config
	//if conf.WebMethodConfig.Oauth2AuthzServerAlias != "" {
	if slices.Contains(servers, conf.WebMethodConfig.Oauth2AuthzServerAlias) {
		servers = []string{conf.WebMethodConfig.Oauth2AuthzServerAlias}
	} else {
		return nil, errors.New("Invalid Oauth2 Authorization Server alias name")
	}

	scopes := []string{}

	for _, server := range oauthServersResponse.Alias {
		var serverPlusScope string
		for _, scope := range server.Scopes {
			if server.Name == conf.WebMethodConfig.Oauth2AuthzServerAlias {
				serverPlusScope = scope.Name

			} else {
				// Handle multiple IDP todo - create multipel credential definition if needed.
				serverPlusScope = server.Name + "-" + scope.Name
			}
			scopes = append(scopes, serverPlusScope)
		}
	}

	log.Infof("Available scopes from IDP %v", scopes)

	corsProp := getCorsSchemaPropertyBuilder()
	agent.RegisterProvisioner(subs.NewProvisioner(gatewayClient, logger))
	agent.NewAPIKeyAccessRequestBuilder().Register()
	agent.NewAPIKeyCredentialRequestBuilder(coreagent.WithCRDRequestSchemaProperty(corsProp)).IsRenewable().Register()

	oAuthRedirects := getAuthRedirectSchemaPropertyBuilder()
	oAuthServers := provisioning.NewSchemaPropertyBuilder().
		SetName(subscription.OauthServerField).SetRequired().SetLabel("Oauth Server").
		IsString().SetEnumValues(servers)

	oAuthType := provisioning.NewSchemaPropertyBuilder().
		SetName(subscription.ApplicationTypeField).SetRequired().SetLabel("Application Type").
		IsString().SetEnumValues([]string{"Confidential", "Public"}).SetDefaultValue("Confidential")

	// audience := provisioning.NewSchemaPropertyBuilder().
	// 	SetName(subscription.AudienceField).SetLabel("Audience").IsString().SetAsTextArea()

	oAuthApiScope := provisioning.NewSchemaPropertyBuilder().
		SetName(subscription.OauthScopes).SetLabel("Scopes").IsArray().AddItem(
		provisioning.NewSchemaPropertyBuilder().SetName("scope").IsString().SetEnumValues(scopes))

	agent.NewAccessRequestBuilder().SetName(subscription.OAuth2AuthType).Register()

	agent.NewOAuthCredentialRequestBuilder(
		coreagent.WithCRDOAuthSecret(),
		coreagent.WithCRDRequestSchemaProperty(oAuthServers),
		coreagent.WithCRDRequestSchemaProperty(oAuthType),
		//	coreagent.WithCRDRequestSchemaProperty(audience),
		coreagent.WithCRDRequestSchemaProperty(oAuthApiScope),
		coreagent.WithCRDRequestSchemaProperty(oAuthRedirects),
		coreagent.WithCRDRequestSchemaProperty(corsProp)).SetName(subscription.OAuth2AuthType).IsRenewable().Register()

	discoveryAgent = discovery.NewAgent(conf, gatewayClient)
	return conf, nil
}

func getCorsSchemaPropertyBuilder() provisioning.PropertyBuilder {
	return provisioning.NewSchemaPropertyBuilder().
		SetName(subscription.CorsField).
		SetLabel("Javascript Origins").
		IsArray().
		AddItem(
			provisioning.NewSchemaPropertyBuilder().
				SetName("Origins").
				IsString())
}

func getAuthRedirectSchemaPropertyBuilder() provisioning.PropertyBuilder {
	return provisioning.NewSchemaPropertyBuilder().
		SetName(subscription.RedirectURLsField).
		SetLabel("Redirect URLs").
		IsArray().
		AddItem(
			provisioning.NewSchemaPropertyBuilder().
				SetName("URL").
				IsString())
}
