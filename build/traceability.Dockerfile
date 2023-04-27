# Build image
# golang:1.19.8-alpine3.17 linux/amd64
FROM docker.io/golang@sha256:841c160ed35923d96c95c52403c4e6db5decd9cbce034aa851e412ade5d4b74f as builder

ENV APP_HOME /build
ENV APP_USER axway

RUN mkdir -p $APP_HOME /app

WORKDIR $APP_HOME
# Copy necessary files
COPY . .

RUN rm -rf bin

RUN make download
RUN make verify
RUN CGO_ENABLED=0  GOOS=linux GOARCH=amd64  make build-trace-docker

# Create non-root user
RUN addgroup $APP_USER && adduser --system $APP_USER --ingroup $APP_USER

RUN mkdir /app/data && \
  apk add ca-certificates && apk update && update-ca-certificates \
  apk --no-cache add curl=7.69.1-r0 && \
  chown -R $APP_USER /data && \
  find / -perm /6000 -type f -exec chmod a-s {} \; || true

RUN chgrp -R 0 /app && chmod -R g=u /app && chown -R $APP_USER /app

# alpine 3.17.3
FROM docker.io/alpine@sha256:b6ca290b6b4cdcca5b3db3ffa338ee0285c11744b4a6abaa9627746ee3291d8d
ENV APP_HOME /build
ENV APP_USER axway

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app /app
COPY --from=builder $APP_HOME/build/webmethods_traceability_agent.yml /app/webmethods_traceability_agent.yml
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

USER $APP_USER
VOLUME ["/tmp"]
HEALTHCHECK --retries=1 CMD curl --fail http://localhost:${STATUS_PORT:-8989}/status || exit 1
ENTRYPOINT ["/app/traceability","--path.config", "/app"]
