# Build image
# golang:1.20.6-alpine3.18 linux/amd64 
FROM docker.io/golang@sha256:6f592e0689192b7e477313264bb190024d654ef0a08fb1732af4f4b498a2e8ad AS builder
ENV APP_HOME /go/src/github.com/Axway/agents-webmethods
ENV APP_USER axway
ENV AGENT=${APP_HOME}/cmd/traceability


RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME

# Copy necessary files
COPY . .

RUN export time=`date +%Y%m%d%H%M%S` && \
  export commit_id=`git rev-parse --short HEAD` && \
  export version=`git tag -l --sort='version:refname' | grep -Eo '[0-9]{1,}\.[0-9]{1,}\.[0-9]{1,3}$' | tail -1` && \
  export sdk_version=`go list -m github.com/Axway/agent-sdk | awk '{print $2}' | awk -F'-' '{print substr($1, 2)}'` && \
  export GOOS=linux && \
  export CGO_ENABLED=0 && \
  export GOARCH=amd64 && \
  go build -tags static_all \
  -ldflags="-X 'github.com/Axway/agent-sdk/pkg/cmd.BuildTime=${time}' \
  -X 'github.com/Axway/agent-sdk/pkg/cmd.BuildVersion=${version}' \
  -X 'github.com/Axway/agent-sdk/pkg/cmd.BuildCommitSha=${commit_id}' \
  -X 'github.com/Axway/agent-sdk/pkg/cmd.SDKBuildVersion=${sdk_version}' \
  -X 'github.com/Axway/agent-sdk/pkg/cmd.BuildAgentName=webMethodsTraceabilityAgent'" \
  -a -o ${APP_HOME}/bin/webmethods_traceability_agent ${AGENT}/main.go


# Create non-root user
RUN addgroup -g 2500 $APP_USER && adduser -u 2500 -D -G $APP_USER $APP_USER
RUN chown -R $APP_USER:$APP_USER  ${APP_HOME}/bin/webmethods_traceability_agent

USER $APP_USER

# alpine 3.18.2
FROM docker.io/alpine@sha256:25fad2a32ad1f6f510e528448ae1ec69a28ef81916a004d3629874104f8a7f70

ENV APP_USER axway
ENV APP_HOME /go/src/github.com/Axway/agents-webmethods

# Copy binary, user, config file and certs from previous build step
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder $APP_HOME/build/webmethods_traceability_agent.yml /webmethods_traceability_agent.yml
COPY --from=builder ${APP_HOME}/bin/webmethods_traceability_agent /webmethods_traceability_agent

RUN mkdir /keys /data && \
  chown -R axway /keys /data && \
  apk --no-cache add openssl libssl libcrypto musl musl-utils libc6-compat busybox curl && \
  find / -perm /6000 -type f -exec chmod a-s {} \; || true


USER $APP_USER
VOLUME ["/keys", "/data"]
HEALTHCHECK --retries=1 CMD curl --fail http://localhost:${STATUS_PORT:-8989}/status || exit 1
ENTRYPOINT ["/webmethods_traceability_agent"]