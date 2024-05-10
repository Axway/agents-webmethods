.PHONY: clean dep-check dep update-sdk dep-version dep-sdk test test-sonar sonar build docker

WORKSPACE ?= $$(pwd)
GO_PKG_LIST := $(shell go list ./...)
SDK_VERSION := $(shell go list -m github.com/Axway/agent-sdk | cut -d ' ' -f 2 | cut -c 2-)
VERSION := $(shell git tag -l --sort="version:refname" | grep -Eo "[0-9]{1,}\.[0-9]{1,}\.[0-9]{1,}" | tail -1)
COMMIT_ID := $(shell git rev-parse --short HEAD)

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

sdk-version:
	@echo $(SDK_VERSION)

run-discovery:
	@go run ./cmd/discovery/main.go

run-trace:
	@go run ./cmd/traceability/main.go

build-discovery:
	@echo "building discovery agent with sdk version $(SDK_VERSION)"
	export CGO_ENABLED=0
	export TIME=`date +%Y%m%d%H%M%S`
	@go build \
		-ldflags="-X 'github.com/Axway/agent-sdk/pkg/cmd.BuildTime=${TIME}' \
			-X 'github.com/Axway/agent-sdk/pkg/cmd.BuildVersion=${VERSION}' \
			-X 'github.com/Axway/agent-sdk/pkg/cmd.BuildCommitSha=${COMMIT_ID}' \
			-X 'github.com/Axway/agent-sdk/pkg/cmd.SDKBuildVersion=${SDK_VERSION}' \
			-X 'github.com/Axway/agent-sdk/pkg/cmd.BuildAgentName=WebmethodsDiscoveryAgent'" \
		-o bin/discovery ./cmd/discovery/main.go
	@echo "discovery agent binary placed at bin/discovery"

build-trace:
	@echo "building traceability agent with sdk version $(SDK_VERSION)"
	export CGO_ENABLED=0
	export TIME=`date +%Y%m%d%H%M%S`
	@go build \
		-ldflags="-X 'github.com/Axway/agent-sdk/pkg/cmd.BuildTime=${TIME}' \
			-X 'github.com/Axway/agent-sdk/pkg/cmd.BuildVersion=${VERSION}' \
			-X 'github.com/Axway/agent-sdk/pkg/cmd.BuildCommitSha=${COMMIT_ID}' \
			-X 'github.com/Axway/agent-sdk/pkg/cmd.SDKBuildVersion=${SDK_VERSION}' \
			-X 'github.com/Axway/agent-sdk/pkg/cmd.BuildAgentName=WebmethodsTraceabilityAgent'" \
		-o bin/traceability ./cmd/traceability/main.go
	@echo "traceability agent binary placed at bin/traceability"

build-trace-docker:
	@go build -o /app/traceability ./cmd/traceability/main.go

docker-build-discovery:
	@docker build -t webmethods_discovery_agent:latest -f ${WORKSPACE}/build/discovery.Dockerfile .
	@echo "Docker build complete"

docker-build-traceability:
	@docker build -t webmethods_traceability_agent:latest -f ${WORKSPACE}/build/traceability.Dockerfile .
	@echo "Docker build complete"
