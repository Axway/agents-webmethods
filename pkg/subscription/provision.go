package subscription

import (
	"errors"
	"fmt"

	"git.ecd.axway.org/apigov/agents-webmethods/pkg/common"
	"git.ecd.axway.org/apigov/agents-webmethods/pkg/webmethods"
	prov "github.com/Axway/agent-sdk/pkg/apic/provisioning"
	"github.com/Axway/agent-sdk/pkg/util"
	"github.com/Axway/agent-sdk/pkg/util/log"
	"github.com/sirupsen/logrus"
)

const (
	// CorsField -
	CorsField = "cors"
	// RedirectURLsField -
	RedirectURLsField = "redirectURLs"
	OauthServerField  = "oauthServer"

	OAuth2AuthType = "oauth2"

	ApplicationTypeField = "applicationType"
	// ClientTypeField -
	ClientTypeField = "clientType"
	AudienceField   = "audience"
)

type provisioner struct {
	client webmethods.Client
	log    logrus.FieldLogger
}

// NewProvisioner creates a type to implement the SDK Provisioning methods for handling subscriptions
func NewProvisioner(client webmethods.Client, log logrus.FieldLogger) prov.Provisioning {
	return &provisioner{
		client: client,
		log:    log.WithField("component", "mp-provisioner"),
	}
}

// AccessRequestDeprovision deletes a contract
func (p provisioner) AccessRequestDeprovision(req prov.AccessRequest) prov.RequestStatus {
	p.log.Info("deprovisioning access request")
	rs := prov.NewRequestStatusBuilder()
	instDetails := req.GetInstanceDetails()

	apiID := util.ToString(instDetails[common.AttrAPIID])
	if apiID == "" {
		return p.failed(rs, notFound(common.AttrAPIID))
	}

	// process access request delete
	// webmethodsApplicationId := req.GetApplicationDetailsValue(common.AttrAppID)
	// if webmethodsApplicationId == "" {
	// 	return p.failed(rs, notFound(common.AttrAppID))
	// }

	// err := p.client.UnsubscribeApplication(webmethodsApplicationId, apiID)
	// if err != nil {
	// 	return p.failed(rs, notFound("Error removing API from Webmethods Application"))
	// }

	p.log.
		WithField("api", apiID).
		WithField("app", req.GetApplicationName()).
		Info("removed access")
	return rs.Success()
}

// AccessRequestProvision adds an API to an app
func (p provisioner) AccessRequestProvision(req prov.AccessRequest) (prov.RequestStatus, prov.AccessData) {
	p.log.Info("provisioning access request")
	rs := prov.NewRequestStatusBuilder()
	instDetails := req.GetInstanceDetails()

	apiID := util.ToString(instDetails[common.AttrAPIID])
	if apiID == "" {
		return p.failed(rs, notFound(common.AttrAPIID)), nil
	}

	webmethodsApplicationId := req.GetApplicationDetailsValue(common.AttrAppID)
	if webmethodsApplicationId == "" {
		return p.failed(rs, notFound(common.AttrAppID)), nil
	}

	webmethodsApplication, err := p.client.GetApplication(webmethodsApplicationId)
	if err != nil || len(webmethodsApplication.Applications) == 0 {
		return p.failed(rs, notFound("Webmethods Application not found")), nil
	}

	apiIds := []string{apiID}
	applicationApiSubscription := webmethods.ApplicationApiSubscription{
		ApiIDs: apiIds,
	}

	err = p.client.SubscribeApplication(webmethodsApplicationId, &applicationApiSubscription)
	if err != nil {
		return p.failed(rs, notFound("Error assocating API to Webmethods Application")), nil
	}
	// process access request create
	rs.AddProperty(common.AttrAppID, webmethodsApplicationId)
	p.log.
		WithField("api", apiID).
		WithField("app", req.GetApplicationName()).
		Info("granted access")
	return rs.Success(), nil
}

// ApplicationRequestDeprovision deletes an app
func (p provisioner) ApplicationRequestDeprovision(req prov.ApplicationRequest) prov.RequestStatus {
	p.log.Info("deprovisioning application")
	rs := prov.NewRequestStatusBuilder()

	appID := req.GetApplicationDetailsValue(common.AppID)
	webmethodsApplicationId := req.GetApplicationDetailsValue(common.AttrAppID)
	if webmethodsApplicationId == "" {
		return p.failed(rs, notFound(common.AttrAppID))
	}
	err := p.client.DeleteApplication(webmethodsApplicationId)
	if err != nil {
		return p.failed(rs, notFound("Error Deleting Webmethods application"))
	}
	log.Infof("Application with Id %s deleted successfully on webmethods", webmethodsApplicationId)
	p.log.
		WithField("appName", req.GetManagedApplicationName()).
		WithField("appID", appID).
		Info("removed application")
	return rs.Success()
}

// ApplicationRequestProvision creates an app
func (p provisioner) ApplicationRequestProvision(req prov.ApplicationRequest) prov.RequestStatus {
	p.log.Info("provisioning application")
	rs := prov.NewRequestStatusBuilder()

	appName := req.GetManagedApplicationName()
	if appName == "" {
		return p.failed(rs, notFound("managed application name"))
	}

	// process application create
	var application webmethods.Application
	application.Name = appName
	application.Version = "1.0"
	application.Description = "Amplify " + appName
	createdApplication, err := p.client.CreateApplication(&application)
	if err != nil {
		return p.failed(rs, notFound("Error creating application"))
	}
	rs.AddProperty(common.AttrAppID, createdApplication.Id)
	p.log.
		WithField("appName", req.GetManagedApplicationName()).
		Info("created application")

	return rs.Success()
}

// CredentialDeprovision returns success since credentials are removed with the app
func (p provisioner) CredentialDeprovision(req prov.CredentialRequest) prov.RequestStatus {
	msg := "credentials will be removed when the subscription is deleted"
	p.log.Info(msg)
	rs := prov.NewRequestStatusBuilder()
	log.Infof("Credential Type %s", req.GetCredentialType())
	// process credential delete
	webmethodsApplicationId := req.GetCredentialDetailsValue(common.AttrAppID)
	log.Infof("webmethodsApplicationId : %s", webmethodsApplicationId)

	if webmethodsApplicationId == "" {
		return p.failed(rs, notFound(common.AttrAppID))
	}

	switch req.GetCredentialType() {
	case prov.APIKeyCRD:
		err := p.client.DeleteApplicationAccessTokens(webmethodsApplicationId)
		if err != nil {
			return p.failed(rs, notFound("Unable to clear application credentials from Webmethods"))
		}
	case OAuth2AuthType:
		log.Info("Removing oauth credential")
	}

	return rs.Success()
}

// CredentialProvision retrieves the credentials from an app
func (p provisioner) CredentialProvision(req prov.CredentialRequest) (prov.RequestStatus, prov.Credential) {
	p.log.Info("provisioning credentials")
	rs := prov.NewRequestStatusBuilder()

	appName := req.GetApplicationName()
	if appName == "" {
		return p.failed(rs, notFound("appName")), nil
	}

	webmethodsApplicationId := req.GetApplicationDetailsValue(common.AttrAppID)
	log.Infof("webmethodsApplicationId : %s", webmethodsApplicationId)
	if webmethodsApplicationId == "" {
		return p.failed(rs, notFound(common.AttrAppID)), nil
	}

	log.Infof("Credential Type %s", req.GetCredentialType())
	applicationsResponse, err := p.client.GetApplication(webmethodsApplicationId)
	if err != nil {
		return p.failed(rs, notFound("Unable to get application from Webmethods")), nil
	}
	var credential prov.Credential
	provData := getCredProvData(req.GetCredentialData())

	switch req.GetCredentialType() {
	case prov.APIKeyCRD:
		credential = prov.NewCredentialBuilder().SetAPIKey(applicationsResponse.Applications[0].AccessTokens.ApiAccessKeyCredentials.ApiAccessKey)
	case OAuth2AuthType:
		credential, err = createOrGetOauthCredential(applicationsResponse.Applications[0], provData, p)
		if err != nil {
			return p.failed(rs, notFound(err.Error())), nil
		}
	}
	rs.AddProperty(common.AttrAppID, webmethodsApplicationId)
	p.log.Info("created credentials")
	return rs.Success(), credential
}

func (p provisioner) CredentialUpdate(req prov.CredentialRequest) (prov.RequestStatus, prov.Credential) {
	p.log.Info("updating credential for app %s", req.GetApplicationName())
	rs := prov.NewRequestStatusBuilder()
	appName := req.GetApplicationName()
	if appName == "" {
		return p.failed(rs, notFound("appName")), nil
	}
	webmethodsApplicationId := req.GetApplicationDetailsValue(common.AttrAppID)
	if webmethodsApplicationId == "" {
		return p.failed(rs, notFound(common.AttrAppID)), nil
	}
	err := p.client.RotateApplicationApikey(webmethodsApplicationId)
	if err != nil {
		return p.failed(rs, notFound("Unable to Rotate Webmethods Application APIkey")), nil
	}
	applicationsResponse, err := p.client.GetApplication(webmethodsApplicationId)
	if err != nil {
		return p.failed(rs, notFound("Unable to get application from Webmethods")), nil
	}
	var credential prov.Credential

	switch req.GetCredentialType() {
	case prov.APIKeyCRD:
		credential = prov.NewCredentialBuilder().SetAPIKey(applicationsResponse.Applications[0].AccessTokens.ApiAccessKeyCredentials.ApiAccessKey)
	case OAuth2AuthType:
		strategyId := applicationsResponse.Applications[0].AuthStrategyIds[0]
		log.Infof("Using existing Oauth Strategy named %s with id %s", applicationsResponse.Applications[0].Name, strategyId)
		strategyResponse, err := p.client.GetStrategy(strategyId)
		if err != nil {
			return p.failed(rs, notFound("Unable to get strategy from Webmethods")), nil
		}
		credential = prov.NewCredentialBuilder().SetOAuthIDAndSecret(strategyResponse.Strategy.ClientRegistration.ClientId, strategyResponse.Strategy.ClientRegistration.ClientSecret)
	}
	p.log.Infof("updated credentials for app %s", req.GetApplicationName())
	return rs.Success(), credential
}

func (p provisioner) failed(rs prov.RequestStatusBuilder, err error) prov.RequestStatus {
	log.Info("handle failed event")
	rs.SetMessage(err.Error())
	p.log.Error(err)
	return rs.Failed()
}

func notFound(msg string) error {
	return fmt.Errorf("%s not found", msg)
}

func getCredProvData(credData map[string]interface{}) credentialMetaData {
	// defaults
	credMetaData := credentialMetaData{
		cors:         []string{"*"},
		redirectURLs: []string{},
		appType:      "Confidential",
	}

	// get cors from credential request
	if data, ok := credData[CorsField]; ok && data != nil {
		credMetaData.cors = []string{}
		for _, c := range data.([]interface{}) {
			credMetaData.cors = append(credMetaData.cors, c.(string))
		}
	}
	// get redirectURLs
	if data, ok := credData[RedirectURLsField]; ok && data != nil {
		credMetaData.redirectURLs = []string{}
		for _, u := range data.([]interface{}) {
			credMetaData.redirectURLs = append(credMetaData.redirectURLs, u.(string))
		}
	}
	// Oauth Server  field
	if data, ok := credData[OauthServerField]; ok && data != nil {
		credMetaData.oauthServerName = data.(string)
	}
	// credential type field
	if data, ok := credData[ApplicationTypeField]; ok && data != nil {
		credMetaData.appType = data.(string)
	}

	// Audience type field
	if data, ok := credData[AudienceField]; ok && data != nil {
		credMetaData.audience = data.(string)
	}

	return credMetaData
}

type credentialMetaData struct {
	cors            []string
	redirectURLs    []string
	oauthServerName string
	appType         string
	audience        string
}

func createOrGetOauthCredential(application webmethods.Application, provData credentialMetaData, p provisioner) (prov.Credential, error) {
	var strategyResponse *webmethods.StrategyResponse
	var err error
	if len(application.AuthStrategyIds) == 0 {
		log.Infof("Creating new Oauth Strategy named %s", application.Name)
		dcrconfig := webmethods.DcrConfig{
			AllowedGrantTypes: []string{"authorization_code",
				"password",
				"client_credentials",
				"refresh_token",
				"implicit"},
			RedirectUris:       provData.redirectURLs,
			AuthServer:         provData.oauthServerName,
			ApplicationType:    "web",
			ClientType:         provData.appType,
			ExpirationInterval: "3600",
			RefreshCount:       "100",
		}
		strategy := &webmethods.Strategy{
			Name:            application.Name,
			Description:     application.Name,
			AuthServerAlias: provData.oauthServerName,
			Audience:        provData.audience,
			Type:            "OAUTH2",
			DcrConfig:       dcrconfig,
		}

		strategyResponse, err = p.client.CreateOauth2Strategy(strategy)
		if err != nil {
			return nil, errors.New("Unable to get application from Webmethods")
		}

		application.AuthStrategyIds = []string{strategyResponse.Strategy.Id}
		applicationsResponse, err := p.client.UpdateApplication(&application)
		if err != nil {
			return nil, errors.New("Unable to get update  Webmethods applicaiton")
		}
		if applicationsResponse == nil {
			return nil, errors.New("Unable to get update  Webmethods applicaiton")
		}
	} else {
		strategyId := application.AuthStrategyIds[0]
		log.Infof("Using existing Oauth Strategy named %s with id %s", application.Name, strategyId)
		strategyResponse, err = p.client.GetStrategy(strategyId)
		if err != nil {
			return nil, errors.New("Unable to get strategy from Webmethods")
		}
	}
	credential := prov.NewCredentialBuilder().SetOAuthIDAndSecret(strategyResponse.Strategy.ClientRegistration.ClientId, strategyResponse.Strategy.ClientRegistration.ClientSecret)
	return credential, nil
}
