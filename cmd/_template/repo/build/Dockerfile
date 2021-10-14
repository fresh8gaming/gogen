ARG VERSION="dev"

FROM golang:1.17-alpine AS builder

ARG VERSION
ARG APP
ENV VERSION_VALUE $VERSION
ENV APP_VALUE $APP

WORKDIR $GOPATH/src/mypackage/myapp/

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags \"-static\" -X \"main.Version=${VERSION_VALUE}\" -X \"main.App=${APP_VALUE}\"" -a -o /go/bin/service ./cmd/${APP_VALUE}/main.go && \
  GRPC_HEALTH_PROBE_VERSION=v0.1.0-alpha.1 && \
  wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
  chmod +x /bin/grpc_health_probe

FROM fresh8gaming/docker-images:base-scratch

COPY --from=builder /go/bin/service /go/bin/service
COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe

CMD ["/go/bin/service"]