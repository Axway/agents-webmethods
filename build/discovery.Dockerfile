# Build image
# golang:1.21.6-alpine3.19 linux/amd64
FROM docker.io/golang@sha256:2523a6f68a0f515fe251aad40b18545155135ca6a5b2e61da8254df9153e3648 AS builder

ARG VERSION
ARG COMMIT_ID

ENV APP_HOME /go/src/github.com/Axway/agents-webmethods
ENV APP_USER axway
ENV AGENT=${APP_HOME}/cmd/discovery


RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME

# Copy necessary files
COPY . .

RUN export time=`date +%Y%m%d%H%M%S` && \
  export commit_id=${COMMIT_ID} && \
  export version=${VERSION} && \
  export sdk_version=`go list -m github.com/Axway/agent-sdk | awk '{print $2}' | awk -F'-' '{print substr($1, 2)}'` && \
  export GOOS=linux && \
  export CGO_ENABLED=0 && \
  export GOARCH=amd64 && \
  go build -tags static_all \
  -ldflags="-X 'github.com/Axway/agent-sdk/pkg/cmd.BuildTime=${time}' \
  -X 'github.com/Axway/agent-sdk/pkg/cmd.BuildVersion=${version}' \
  -X 'github.com/Axway/agent-sdk/pkg/cmd.BuildCommitSha=${commit_id}' \
  -X 'github.com/Axway/agent-sdk/pkg/cmd.SDKBuildVersion=${sdk_version}' \
  -X 'github.com/Axway/agent-sdk/pkg/cmd.BuildAgentName=webMethodsDiscoveryAgent'" \
  -a -o ${APP_HOME}/bin/webmethods_discovery_agent ${AGENT}/main.go

# Create non-root user
RUN addgroup -g 2500 $APP_USER && adduser -u 2500 -D -G $APP_USER $APP_USER
RUN chown -R $APP_USER:$APP_USER  ${APP_HOME}/bin/webmethods_discovery_agent

USER $APP_USER

# alpine 3.19 linux/amd64 
FROM docker.io/alpine@sha256:13b7e62e8df80264dbb747995705a986aa530415763a6c58f84a3ca8af9a5bcd

ENV APP_USER axway
ENV APP_HOME /go/src/github.com/Axway/agents-webmethods

# Copy binary, user, config file and certs from previous build step
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder $APP_HOME/build/webmethods_discovery_agent.yml /webmethods_discovery_agent.yml
COPY --from=builder ${APP_HOME}/bin/webmethods_discovery_agent /webmethods_discovery_agent

RUN mkdir /keys && \
  chown -R axway /keys && \
  apk --no-cache add openssl libssl3 libcrypto3 musl musl-utils libc6-compat busybox curl && \
  find / -perm /6000 -type f -exec chmod a-s {} \; || true


USER $APP_USER
VOLUME ["/keys"]
HEALTHCHECK --retries=1 CMD curl --fail http://localhost:${STATUS_PORT:-8989}/status || exit 1
ENTRYPOINT ["/webmethods_discovery_agent"]
