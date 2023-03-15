module git.ecd.axway.org/apigov/agents-webmethods

go 1.18

replace (
	github.com/Shopify/sarama => github.com/elastic/sarama v1.19.1-0.20210823122811-11c3ef800752
	github.com/docker/docker => github.com/docker/engine v17.12.0-ce-rc1.0.20190717161051-705d9623b7c1+incompatible
	github.com/dop251/goja => github.com/andrewkroh/goja v0.0.0-20190128172624-dd2ac4456e20
	github.com/dop251/goja_nodejs => github.com/dop251/goja_nodejs v0.0.0-20171011081505-adff31b136e6
	github.com/getkin/kin-openapi => github.com/getkin/kin-openapi v0.67.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.21.1
	k8s.io/client-go => k8s.io/client-go v0.21.1
)

require (
	github.com/Axway/agent-sdk v1.1.47
	github.com/elastic/beats/v7 v7.17.7
	github.com/hpcloud/tail v1.0.0
	github.com/sirupsen/logrus v1.9.0
)
