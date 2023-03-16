package discovery

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"git.ecd.axway.org/apigov/agents-webmethods/pkg/common"
	"git.ecd.axway.org/apigov/agents-webmethods/pkg/webmethods"
	"github.com/Axway/agent-sdk/pkg/apic/provisioning"
	"github.com/Axway/agent-sdk/pkg/cache"

	"github.com/sirupsen/logrus"

	"git.ecd.axway.org/apigov/agents-webmethods/pkg/config"
	"github.com/Axway/agent-sdk/pkg/apic"
)

const (
	marketplace = "marketplace"
	catalog     = "unified-catalog"
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

	ard := provisioning.APIKeyARD
	crds := []string{provisioning.APIKey}

	specType, err := getSpecType(api.ApiSpec)
	if err != nil {
		return nil, err
	}

	if specType == "" {
		return nil, fmt.Errorf("unknown spec type")
	}

	return &ServiceDetail{
		AccessRequestDefinition: ard,
		CRDs:                    crds,
		APIName:                 api.Name,
		APISpec:                 api.ApiSpec,
		AuthPolicy:              "api-key",
		Description:             api.Description,
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
func getSpecType(specContent []byte) (string, error) {
	if specContent != nil {
		jsonMap := make(map[string]interface{})
		err := json.Unmarshal(specContent, &jsonMap)
		if err != nil {
			return "", err
		}
		if _, isSwagger := jsonMap["swagger"]; isSwagger {
			return apic.Oas2, nil
		} else if _, isOpenAPI := jsonMap["openapi"]; isOpenAPI {
			return apic.Oas3, nil
		}
	}
	return "", nil
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
