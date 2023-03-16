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
}

type WebmethodsApi struct {
	ApiName          string
	ApiVersion       string
	ApiDescription   string
	IsActive         bool
	ApiType          string `type`
	TracingEnabled   bool
	PublishedPortals []string
	SystemVersion    int
	Id               string
}

type ApiResponse struct {
	Api              Api
	GatewayEndPoints []string
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
