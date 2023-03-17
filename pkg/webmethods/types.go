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
