.PHONY: clean dep-check dep update-sdk dep-version dep-sdk test test-sonar sonar build docker

WORKSPACE ?= $$(pwd)
GO_PKG_LIST := $(shell go list ./...)
export GOFLAGS := -mod=mod
export GOPRIVATE := git.ecd.axway.org

all: clean dep test build docker
	@echo "Done"

clean:
	@rm -rf ./bin/
	@mkdir -p ./bin
	@echo "Clean complete"

dep-check:
	@go mod verify

dep:
	@echo "Resolving go package dependencies"
	@go mod tidy
	@echo "Package dependencies completed"

update-sdk:
	@echo "Updating SDK dependencies"
	@export GOFLAGS="" && go mod edit -require "github.com/Axway/agent-sdk@${version}"

dep-branch:
	@make sdk=`git branch --show-current` dep-version

dep-version:
	@export version=$(sdk) && make update-sdk && make dep

dep-sdk: 
	@make sdk=main dep-version

test: dep
	@go vet ${GO_PKG_LIST}
	@go test -race -short -coverprofile=${WORKSPACE}/gocoverage.out -count=1 ${GO_PKG_LIST}

test-sonar: dep
	@go vet ${GO_PKG_LIST}
	@go test -short -coverpkg=./... -coverprofile=${WORKSPACE}/gocoverage.out -count=1 ${GO_PKG_LIST} -json > ${WORKSPACE}/goreport.json

sonar: test-sonar
	./sonar.sh $(sonarHost)

${WORKSPACE}/bin/webmethods_agents:
	go version
	@export time=`date +%Y%m%d%H%M%S` && \
	export version=`cat version` && \
	export commit_id=`cat commit_id` && \
	export sdk_version=`go list -m github.com/Axway/agent-sdk | awk '{print $$2}' | awk -F'-' '{print substr($$1, 2)}'` && \
	go build -tags static_all \
		-ldflags="-X 'github.com/Axway/agent-sdk/pkg/cmd.BuildTime=$${time}' \
				-X 'github.com/Axway/agent-sdk/pkg/cmd.BuildVersion=$${version}' \
				-X 'github.com/Axway/agent-sdk/pkg/cmd.BuildCommitSha=$${commit_id}' \
				-X 'github.com/Axway/agent-sdk/pkg/cmd.SDKBuildVersion=$${sdk_version}' \
				-X 'github.com/Axway/agent-sdk/pkg/cmd.BuildAgentName=WebMethodsAgents' \
				-X 'github.com/Axway/agent-sdk/pkg/cmd.BuildAgentDescription=Amplify webMethods Agents'" \
		-a -o ${WORKSPACE}/bin/webmethods_agent ${WORKSPACE}/webmethods_agent.go

build: dep ${WORKSPACE}/bin/azure_discovery_agent
	@echo "Build complete"

docker: dep
	docker build -t azure_discovery_agent:latest -f ${WORKSPACE}/build/docker/Dockerfile .
	@echo "Docker build complete"