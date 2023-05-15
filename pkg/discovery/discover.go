package discovery

import (
	"time"

	"github.com/Axway/agent-sdk/pkg/apic"

	"github.com/Axway/agent-sdk/pkg/cache"

	"github.com/Axway/agents-webmethods/pkg/config"
	"github.com/Axway/agents-webmethods/pkg/webmethods"

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
	maturityState     string
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
	apis, err := d.client.SearchAPIs()
	log.Infof("%v", apis)
	if err != nil {
		log.Errorf("Unable to list Apis :%v", err)
		return
	}

	for _, api := range apis.WebmethodsApi {
		go func(api webmethods.WebmethodsApi) {
			apiResponse, err := d.client.GetApiDetails(api.Id)
			if err != nil {
				log.Errorf("Unable to get API Details : %v", err)
				return
			}

			if !d.client.IsAllowedTags(apiResponse.Api.ApiDefinition.Tags) {
				log.Infof("API not matched with filtered tags : %v, hence ignoring for discovery", err)
				return
			}

			if apiResponse.Api.MaturityState == d.maturityState {
				var specification []byte
				if api.ApiType == "REST" {
					specification, err = d.client.GetApiSpec(api.Id)
				} else if api.ApiType == "SOAP" {
					specification, err = d.client.GetWsdl(apiResponse.GatewayEndPoints[0])
				}
				if err != nil {
					log.Errorf("Unable to read API specification : %v", err)
					return
				}
				amplifyApi := webmethods.AmplifyAPI{
					ID:          api.Id,
					Name:        api.ApiName,
					Description: api.ApiDescription,
					Version:     api.ApiVersion,
					// append endpoint url
					Url:           apiResponse.GatewayEndPoints[0],
					Documentation: []byte(apiResponse.Api.ApiDefinition.Info.Description),
					ApiSpec:       specification,
					ApiType:       api.ApiType,
				}
				svcDetail := d.serviceHandler.ToServiceDetail(&amplifyApi)
				if svcDetail != nil {
					d.apiChan <- svcDetail
				}
			} else {
				log.Infof("Ignoring API %s with MaturityState %s", api.ApiName, apiResponse.Api.MaturityState)
			}

		}(api)
	}
}
