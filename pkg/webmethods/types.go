package webmethods

// API -
type API struct {
	ApiSpec       []byte
	ID            string
	Name          string
	Description   string
	Version       string
	Url           string
	Documentation []byte
}

type WebmthodsApi struct {
	apiName          string
	apiVersion       string
	isActive         bool
	apiType          string `type`
	tracingEnabled   bool
	publishedPortals []string
	systemVersion    int
	id               string
}
