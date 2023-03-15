package discovery

import (
	"time"

	"github.com/Axway/agent-sdk/pkg/apic"

	"github.com/Axway/agent-sdk/pkg/cache"

	"git.ecd.axway.org/apigov/agents-webmethods/pkg/boomi"
	"git.ecd.axway.org/apigov/agents-webmethods/pkg/config"
	"git.ecd.axway.org/apigov/agents-webmethods/pkg/webmethods"

	"github.com/sirupsen/logrus"

	"github.com/Axway/agent-sdk/pkg/util/log"
)

// discovery implements the Repeater interface. Polls boomi for APIs.
type discovery struct {
	apiChan           chan *ServiceDetail
	cache             cache.Cache
	centralClient     apic.Client
	client            webmethods.ListAPIClient
	discoveryPageSize int
	pollInterval      time.Duration
	stopDiscovery     chan bool
	serviceHandler    ServiceHandler
}

func (d *discovery) Stop() {
	d.stopDiscovery <- true
}

func (d *discovery) OnConfigChange(cfg *config.BoomiConfig) {
	d.pollInterval = cfg.PollInterval
	d.serviceHandler.OnConfigChange(cfg)
}

// Loop Discovery event loop.
func (d *discovery) Loop() {
	go func() {
		// Instant fist "tick"
		d.discoverAPIs()
		logrus.Info("Starting poller for Boomi APIs")
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
		go func(api boomi.API) {
			svcDetail := d.serviceHandler.ToServiceDetail(&api)
			if svcDetail != nil {
				d.apiChan <- svcDetail
			}
		}(api)
	}
}
