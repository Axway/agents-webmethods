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

type Apis struct {
	WebmethodsApi []WebmethodsApi `json:"api"`
}

type ListApiResponse struct {
	WebmethodsApi  WebmethodsApi `json:"api"`
	ResponseStatus string        `json:"responseStatus"`
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
}

type Api struct {
	ApiDescription string
	ApiVersion     string
	Title          string
	ApiDefinition  ApiDefinition
	MaturityState  string
	ApiGroups      []string
	Owner          string
}

type ApiDefinition struct {
	Info Info
	Tags []Tag `json:"tags"`
}
type Info struct {
	Description string
	Version     string
	Title       string
}

type Tag struct {
	Name string
}

// application

type AccessTokens struct {
	ApiAccessKeyCredentials ApiAccessKeyCredentials `json:"apiAccessKey_credentials"`
	Oauth2Token             Oauth2Token             `json:"oauth2Token"`
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

type ApiAccessKeyCredentials struct {
	ApiAccessKey       string `json:"apiAccessKey"`
	ExpirationInterval string `json:"expirationInterval"`
	ExpirationDate     string `json:"expirationDate"`
}

type ApplicationResponse struct {
	Applications []Application `json:"applications"`
}

type SearchApplicationResponse struct {
	SearchApplication []SearchApplication `json:"application"`
}

type SearchApplication struct {
	ApplicationID string `json:"applicationID"`
	Name          string `json:"name"`
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

type ApplicationApiSubscription struct {
	ApiIDs []string `json:"apiIDs"`
}

type Search struct {
	Types          []string `json:"types"`
	Scope          []Scope  `json:"scope"`
	ResponseFields []string `json:"responseFields"`
	Condition      string   `json:"condition,omitempty"`
	SortByField    string   `json:"sortByField,omitempty"`
}

type Scope struct {
	AttributeName string `json:"attributeName"`
	Keyword       string `json:"keyword"`
}

type Strategy struct {
	Id                 string             `json:"id,omitempty"`
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	AuthServerAlias    string             `json:"authServerAlias"`
	Audience           string             `json:"audience"`
	Type               string             `json:"type"`
	DcrConfig          DcrConfig          `json:"dcrConfig"`
	ClientRegistration ClientRegistration `json:"clientRegistration"`
}

type DcrConfig struct {
	AllowedGrantTypes  []string `json:"allowedGrantTypes"`
	Scopes             []string `json:"scopes"`
	RedirectUris       []string `json:"redirectUris"`
	AuthServer         string   `json:"authServer"`
	ApplicationType    string   `json:"applicationType"`
	ClientType         string   `json:"clientType"`
	ExpirationInterval string   `json:"expirationInterval"`
	RefreshCount       string   `json:"refreshCount"`
	PkceType           string   `json:"pkceType"`
}

type StrategyResponse struct {
	Strategy Strategy `json:"strategy"`
}

type ClientRegistration struct {
	Shell                    bool          `json:"shell"`
	ClientId                 string        `json:"clientId"`
	Name                     string        `json:"name"`
	TokenLifetime            int           `json:"tokenLifetime"`
	TokenRefreshLimit        int           `json:"tokenRefreshLimit"`
	ClientSecret             string        `json:"clientSecret"`
	Enabled                  bool          `json:"enabled"`
	RedirectUris             []string      `json:"redirectUris"`
	ClScopes                 []interface{} `json:"clScopes"`
	AuthCodeAllowed          bool          `json:"authCodeAllowed"`
	ImplicitAllowed          bool          `json:"implicitAllowed"`
	ClientCredentialsAllowed bool          `json:"clientCredentialsAllowed"`
	ResourceOwnerAllowed     bool          `json:"resourceOwnerAllowed"`
	PkceType                 string        `json:"pkceType"`
}

type OauthServers struct {
	Alias []struct {
		ID                  string       `json:"id"`
		Name                string       `json:"name"`
		Description         string       `json:"description,omitempty"`
		Type                string       `json:"type"`
		Scopes              []OauthScope `json:"scopes"`
		SupportedGrantTypes []string     `json:"supportedGrantTypes"`
	} `json:"alias"`
}

type OauthScope struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
