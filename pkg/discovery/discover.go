package discovery

import (
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
	apis, err := d.client.ListAPIs()
	if err != nil {
		log.Errorf("Unable to list Apis :%v", err)
		return
	}

	for _, api := range apis {
		go func(api webmethods.ListApiResponse) {
			if api.ResponseStatus == "SUCCESS" {
				apiResponse, err := d.client.GetApiDetails(api.WebmethodsApi.Id)
				if err != nil {
					log.Errorf("Unable to get API Details : %v", err)
					return
				}

				if !d.client.IsAllowedTags(apiResponse.Api.ApiDefinition.Tags) {
					log.Infof("API matched with filtered tags : %v, hence ignoring for discovery", err)
					return
				}

				if apiResponse.Api.MaturityState == d.maturityState {
					var specification []byte
					if api.WebmethodsApi.ApiType == "REST" {
						specification, err = d.client.GetApiSpec(api.WebmethodsApi.Id)
					} else if api.WebmethodsApi.ApiType == "SOAP" {
						specification, err = d.client.GetWsdl(apiResponse.GatewayEndPoints[0])
					}
					if err != nil {
						log.Error("Unable to read API specification : ", err)
						return
					}
					amplifyApi := webmethods.AmplifyAPI{
						ID:          api.WebmethodsApi.Id,
						Name:        api.WebmethodsApi.ApiName,
						Description: api.WebmethodsApi.ApiDescription,
						Version:     api.WebmethodsApi.ApiVersion,
						// append endpoint url
						Url:           apiResponse.GatewayEndPoints[0],
						Documentation: []byte(apiResponse.Api.ApiDefinition.Info.Description),
						ApiSpec:       specification,
						//	AuthPolicy:    authPolicy,
						ApiType: api.WebmethodsApi.ApiType,
					}
					svcDetail := d.serviceHandler.ToServiceDetail(&amplifyApi)
					if svcDetail != nil {
						d.apiChan <- svcDetail
					}
				} else {
					log.Infof("Ignoring API %s with MaturityState %s", api.WebmethodsApi.ApiName, apiResponse.Api.MaturityState)
				}
			}
		}(api)
	}
}
