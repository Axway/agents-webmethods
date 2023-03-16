package discovery

import (
	"fmt"
	"time"

	"github.com/Axway/agent-sdk/pkg/apic"

	"github.com/Axway/agent-sdk/pkg/cache"

	"git.ecd.axway.org/apigov/agents-webmethods/pkg/config"
	"git.ecd.axway.org/apigov/agents-webmethods/pkg/webmethods"

	"github.com/sirupsen/logrus"

	"github.com/Axway/agent-sdk/pkg/util/log"
)

// discovery implements the Repeater interface. Polls Webmethos  APIM.
type discovery struct {
	apiChan           chan *ServiceDetail
	cache             cache.Cache
	centralClient     apic.Client
	client            webmethods.Client
	discoveryPageSize int
	pollInterval      time.Duration
	stopDiscovery     chan bool
	serviceHandler    ServiceHandler
}

func (d *discovery) Stop() {
	d.stopDiscovery <- true
}

func (d *discovery) OnConfigChange(cfg *config.WebMethodConfig) {
	d.pollInterval = cfg.PollInterval
	d.serviceHandler.OnConfigChange(cfg)
}

// Loop Discovery event loop.
func (d *discovery) Loop() {
	go func() {
		// Instant fist "tick"
		d.discoverAPIs()
		logrus.Info("Starting poller for Webmethods APIM")
		ticker := time.NewTicker(d.pollInterval)
		for {
			select {
			case <-ticker.C:
				d.discoverAPIs()
				break
			case <-d.stopDiscovery:
				log.Debug("stopping discovery loop")
				ticker.Stop()
				break
			}
		}
	}()
}

// discoverAPIs Finds APIs from exchange
func (d *discovery) discoverAPIs() {
	apis, err := d.client.ListAPIs()
	if err != nil {
		log.Error(err)
	}

	for _, api := range apis {
		go func(api webmethods.WebmethodsApi) {
			apiResponse, err := d.client.GetApiDetails(api.Id)
			if err != nil {
				panic(fmt.Sprintf("Unable to decompress : %s", err))
			}

			var specification []byte
			if api.ApiType == "REST" {
				specification, err = d.client.GetApiSpec(api.Id)
			} else if api.ApiType == "SOAP" {
				specification, err = d.client.GetWsdl(apiResponse.GatewayEndPoints[0])
			}
			if err != nil {
				panic(fmt.Sprintf("Unable to decompress : %s", err))
			}
			authPolicy := handleAuthPolicy()

			amplifyApi := webmethods.AmplifyAPI{
				ID:          api.Id,
				Name:        api.ApiName,
				Description: api.ApiDescription,
				Version:     api.ApiVersion,
				// append endpoint url
				Url:           apiResponse.GatewayEndPoints[0],
				Documentation: []byte(apiResponse.Api.ApiDefinition.Info.Description),
				ApiSpec:       specification,
				AuthPolicy:    authPolicy,
			}
			svcDetail := d.serviceHandler.ToServiceDetail(&amplifyApi)
			if svcDetail != nil {
				d.apiChan <- svcDetail
			}
		}(api)
	}
}

func handleAuthPolicy() string {
	authPolicy := "api-key" // Default auth i.e api key authentication
	return authPolicy
}
