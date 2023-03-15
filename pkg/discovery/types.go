package discovery

// ServiceDetail is the information for the ex
type ServiceDetail struct {
	AccessRequestDefinition string
	APIName                 string
	APISpec                 []byte
	APIUpdateSeverity       string
	AuthPolicy              string
	CRDs                    []string
	Description             string
	Documentation           []byte
	ID                      string
	Image                   string
	ImageContentType        string
	ResourceType            string
	ServiceAttributes       map[string]string
	AgentDetails            map[string]string
	Stage                   string
	State                   string
	Status                  string
	SubscriptionName        string
	Tags                    []string
	Title                   string
	URL                     string
	Version                 string
}
