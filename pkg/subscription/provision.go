package subscription

import (
	"fmt"

	"git.ecd.axway.org/apigov/agents-webmethods/pkg/common"
	"git.ecd.axway.org/apigov/agents-webmethods/pkg/webmethods"
	prov "github.com/Axway/agent-sdk/pkg/apic/provisioning"
	"github.com/Axway/agent-sdk/pkg/util"
	"github.com/sirupsen/logrus"
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

	webmethodsApplicationId := req.GetApplicationDetailsValue("webmethodsApplicationId")
	if webmethodsApplicationId == "" {
		return p.failed(rs, notFound(common.AttrAppID)), nil
	}

	webmethodsApplication, err := p.client.GetApplication(webmethodsApplicationId)
	if err != nil {
		return p.failed(rs, notFound("Webmethods Application not found")), nil
	}
	apiIds := []string{apiID}
	webmethodsApplication.NewApisForAssociation = apiIds

	updatedWebmethodsApplication, err := p.client.UpdateApplication(webmethodsApplication)
	if err != nil {
		return p.failed(rs, notFound("Error assocating API to Webmethods Application")), nil
	}

	var matchAPI bool

	for _, value := range updatedWebmethodsApplication.ConsumingAPIs {
		if value == apiID {
			p.log.Info("Successfully updated application %s with apikey  %s", updatedWebmethodsApplication.Name, apiID)
			matchAPI = true
			break
		}
	}

	if !matchAPI {
		return p.failed(rs, notFound("Error assocating API to Webmethods Application")), nil
	}

	// process access request create

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

	// process application delete

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
	application.Description = appName

	createdApplication, err := p.client.CreateApplication(&application)
	if err != nil {
		return p.failed(rs, notFound("Error creating application"))
	}
	rs.AddProperty("webmethodsApplicationId", createdApplication.Id)
	p.log.
		WithField("appName", req.GetManagedApplicationName()).
		Info("created application")

	return rs.Success()
}

// CredentialDeprovision returns success since credentials are removed with the app
func (p provisioner) CredentialDeprovision(_ prov.CredentialRequest) prov.RequestStatus {
	msg := "credentials will be removed when the subscription is deleted"
	p.log.Info(msg)

	// process credential delete

	return prov.NewRequestStatusBuilder().
		SetMessage("credentials will be removed when the application is deleted").
		Success()
}

// CredentialProvision retrieves the credentials from an app
func (p provisioner) CredentialProvision(req prov.CredentialRequest) (prov.RequestStatus, prov.Credential) {
	p.log.Info("provisioning credentials")
	rs := prov.NewRequestStatusBuilder()

	appName := req.GetApplicationName()
	if appName == "" {
		return p.failed(rs, notFound("appName")), nil
	}

	appID := req.GetApplicationDetailsValue(common.AppID)
	if appName == "" {
		return p.failed(rs, notFound("appID")), nil
	}
	application, err := p.client.GetApplication(appID)
	if err != nil {
		return p.failed(rs, notFound("Unable to get application from Webmethods")), nil
	}
	cr := prov.NewCredentialBuilder().SetAPIKey(application.AccessTokens.ApiAccessKey.APIAccessKey)
	p.log.Info("created credentials")

	return rs.Success(), cr
}

func (p provisioner) CredentialUpdate(req prov.CredentialRequest) (prov.RequestStatus, prov.Credential) {
	p.log.Info("updating credential for app %s", req.GetApplicationName())
	rs := prov.NewRequestStatusBuilder()

	appName := req.GetApplicationName()
	if appName == "" {
		return p.failed(rs, notFound("appName")), nil
	}

	cr := prov.NewCredentialBuilder().SetAPIKey(appName)

	p.log.Infof("updated credentials for app %s", req.GetApplicationName())

	return rs.Success(), cr
}

func (p provisioner) failed(rs prov.RequestStatusBuilder, err error) prov.RequestStatus {
	rs.SetMessage(err.Error())
	p.log.Error(err)
	return rs.Failed()
}

func notFound(msg string) error {
	return fmt.Errorf("%s not found", msg)
}
