package discovery

import (
	"crypto/sha256"
	"fmt"

	"github.com/Axway/agent-sdk/pkg/apic/provisioning"
	"github.com/Axway/agent-sdk/pkg/cache"
	"github.com/Axway/agents-webmethods/pkg/common"
	"github.com/Axway/agents-webmethods/pkg/subscription"
	"github.com/Axway/agents-webmethods/pkg/webmethods"

	"github.com/sirupsen/logrus"

	"github.com/Axway/agent-sdk/pkg/apic"
	"github.com/Axway/agents-webmethods/pkg/config"
)

const (
	marketplace = "marketplace"
)

// ServiceHandler converts a webmethods APIM to an array of ServiceDetails
type ServiceHandler interface {
	ToServiceDetail(api *webmethods.AmplifyAPI) *ServiceDetail
	OnConfigChange(cfg *config.WebMethodConfig)
}

type serviceHandler struct {
	client webmethods.Client
	cache  cache.Cache
	mode   string
}

func (s *serviceHandler) OnConfigChange(cfg *config.WebMethodConfig) {
}

// ToServiceDetails gathers the ServiceDetail for a Webmethods APIM.
func (s *serviceHandler) ToServiceDetail(api *webmethods.AmplifyAPI) *ServiceDetail {
	logger := logrus.WithFields(logrus.Fields{
		"name": api.Name,
		"id":   api.ID,
	})

	serviceDetail, err := s.getServiceDetail(api)
	if err != nil {
		logger.Errorf("error getting the service details: %s", err.Error())
	}
	return serviceDetail
}

// getServiceDetail gets the ServiceDetail for the API asset.
func (s *serviceHandler) getServiceDetail(api *webmethods.AmplifyAPI) (*ServiceDetail, error) {
	logger := logrus.WithFields(logrus.Fields{
		"name": api.Name,
		"id":   api.ID,
	})

	isAlreadyPublished, checksum := isPublished(api, s.cache)
	// If true, then the api is published and there were no changes detected
	if isAlreadyPublished {
		logger.Debug("api is already published")
		return nil, nil
	}

	err := s.cache.Set(checksum, *api)
	if err != nil {
		logger.Errorf("failed to save api to cache: %s", err)
	}

	var ardName string
	crds := make([]string, 1)

	specType := getSpecType(api.ApiType)

	if specType == "" {
		return nil, fmt.Errorf("unknown spec type")
	}

	logger.Infof("Processing Spec type :%s", specType)

	if specType == "oas3" {
		specParser := apic.NewSpecResourceParser(api.ApiSpec, apic.Oas3)
		specParser.Parse()
		specProcessor := specParser.GetSpecProcessor()
		if processor, ok := specProcessor.(apic.OasSpecProcessor); ok {
			processor.ParseAuthInfo()
			authPolicies := processor.GetAuthPolicies()
			oauthScopes := processor.GetOAuthScopes()
			oauthScopesList := make([]string, 0, len(oauthScopes))
			for k := range oauthScopes {
				oauthScopesList = append(oauthScopesList, k)
			}
			for _, value := range authPolicies {
				logger.Infof("API Authentication Type %s", value)
				if value == apic.Apikey {
					ardName = provisioning.APIKeyARD
					crds[0] = provisioning.APIKeyARD
					break
				}
				if value == apic.Oauth {
					ardName = subscription.OAuth2AuthType
					crds[0] = subscription.OAuth2AuthType
					break
				}
			}
		}
	} else {
		logger.Info("Ignoring authentication")
	}

	return &ServiceDetail{
		AccessRequestDefinition: ardName,
		CRDs:                    crds,
		APIName:                 api.Name,
		APISpec:                 api.ApiSpec,
		//AuthPolicy:              api.AuthPolicy,
		Description: api.Description,
		// Use the Asset ID for the externalAPIID so that apis linked to the asset are created as a revision
		ID:                api.ID,
		ResourceType:      specType,
		ServiceAttributes: map[string]string{},
		AgentDetails: map[string]string{
			common.AttrAPIID:    api.ID,
			common.AttrChecksum: checksum,
		},
		Title:   api.Name,
		Version: api.Version,
		Status:  apic.PublishedStatus,
	}, nil
}

// getSpecType determines the correct resource type for the asset.
func getSpecType(apiType string) string {

	if apiType == "REST" {
		return apic.Oas3
	} else if apiType == "SOAP" {
		return apic.Wsdl
	}
	return ""
}

// makeChecksum generates a makeChecksum for the api for change detection
func makeChecksum(val interface{}) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%v", val)))
	return fmt.Sprintf("%x", sum)
}

// isPublished checks if an api is published with the latest changes. Returns true if it is, and false if it is not.
func isPublished(api *webmethods.AmplifyAPI, c cache.Cache) (bool, string) {
	// Change detection (asset + policies)
	checksum := makeChecksum(api)
	item, err := c.Get(checksum)
	if err != nil || item == nil {
		return false, checksum
	}
	return true, checksum
}
