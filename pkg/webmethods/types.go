package webmethods

// API -
type AmplifyAPI struct {
	ApiSpec       []byte
	ID            string
	Name          string
	Description   string
	Version       string
	Url           string
	Documentation []byte
	AuthPolicy    string
	ApiType       string
}

type ListApi struct {
	ListApiResponse []ListApiResponse `json:"apiResponse"`
}

type ListApiResponse struct {
	WebmethodsApi WebmethodsApi `json:"api"`
}

type WebmethodsApi struct {
	ApiName          string
	ApiVersion       string
	ApiDescription   string
	IsActive         bool
	ApiType          string `json:"type"`
	TracingEnabled   bool
	PublishedPortals []string
	SystemVersion    int
	Id               string
}

type GetApiDetails struct {
	ApiResponse ApiResponse `json:"apiResponse"`
}

type ApiResponse struct {
	Api              Api      `json:"api"`
	GatewayEndPoints []string `json:"gatewayEndPoints"`
	MaturityState    string
	Owner            string
}

type Api struct {
	ApiDefinition ApiDefinition
}
type ApiDefinition struct {
	Info Info
}
type Info struct {
	Description string
	Version     string
	Title       string
}

// application

type AccessTokens struct {
	ApiAccessKey ApiAccessKey `json:"apiAccessKey"`
	Oauth2Token  Oauth2Token  `json:"oauth2Token"`
}

type Oauth2Token struct {
	Type               string   `json:"type"`
	ClientID           string   `json:"clientId"`
	ClientSecret       string   `json:"clientSecret"`
	ClientName         string   `json:"clientName"`
	Scopes             []string `json:"scopes"`
	ExpirationInterval string   `json:"expirationInterval"`
	RefreshCount       string   `json:"refreshCount"`
	RedirectUris       []string `json:"redirectUris"`
}

type ApiAccessKey struct {
	APIAccessKey       string   `json:"apiAccessKey"`
	ExpirationInterval string   `json:"expirationInterval"`
	ExpirationDate     string   `json:"expirationDate"`
	ApiGroups          []string `json:"apiGroups"`
}

type Application struct {
	Id            string `json:"id"`
	ApplicationID string `json:"applicationID"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Owner         string `json:"owner"`
	Identifiers   []struct {
		ID    string   `json:"id"`
		Key   string   `json:"key"`
		Name  string   `json:"name"`
		Value []string `json:"value"`
	} `json:"identifiers"`
	ContactEmails         []string     `json:"contactEmails"`
	IconbyteArray         string       `json:"iconbyteArray"`
	AccessTokens          AccessTokens `json:"accessTokens"`
	CreationDate          string       `json:"creationDate"`
	LastModified          string       `json:"lastModified"`
	LastUpdated           string       `json:"lastUpdated"`
	SiteURLs              []string     `json:"siteURLs"`
	JsOrigins             []string     `json:"jsOrigins"`
	Version               string       `json:"version"`
	IsSuspended           string       `json:"isSuspended"`
	AuthStrategyIds       []string     `json:"authStrategyIds"`
	Subscription          bool         `json:"subscription"`
	ConsumingAPIs         []string     `json:"consumingAPIs"`
	NewApisForAssociation []string     `json:"newApisForAssociation"`
}

type SavedSettings struct {
	ExtendedKeys ExtendedKeys `json:"extendedKeys"`
}

type ExtendedKeys struct {
	ApiMaturityStatePossibleValues string `json:"apiMaturityStatePossibleValues"`
	ApiGroupingPossibleValues      string `json:"apiGroupingPossibleValues"`
}
