ARG VERSION="dev"
ARG GO_VERSION=1.18
ARG ALPINE_VERSION=3.17

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

ARG VERSION
ARG APP
ENV VERSION_VALUE $VERSION
ENV APP_VALUE $APP

WORKDIR $GOPATH/src/mypackage/myapp/

COPY . .

# hadolint ignore=SC2086
RUN BUILDFLAGS="-w -s -extldflags '-static' -X 'main.Version=${VERSION_VALUE}' -X 'main.App=${APP_VALUE}'" CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags=${BUILDFLAGS} -a -o /go/bin/service ./cmd/${APP_VALUE}/main.go && \
  GRPC_HEALTH_PROBE_VERSION=v0.4.15 && \
  wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
  chmod +x /bin/grpc_health_probe

FROM alpine:${ALPINE_VERSION} AS alpine-builder

RUN apk update && apk add --no-cache ca-certificates=20220614-r4 tzdata=2022f-r1
RUN update-ca-certificates

ENV USER=appuser
ENV UID=1000
ENV GID=1000

RUN addgroup --gid "${GID}" appgroup

# hadolint ignore=DL3059
RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid "${UID}" \
  -G appgroup \
  "${USER}"

FROM scratch

COPY --from=alpine-builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=alpine-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=alpine-builder /etc/passwd /etc/passwd
COPY --from=alpine-builder /etc/group /etc/group

USER 1000:1000

COPY --from=builder /go/bin/service /go/bin/service
COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe

CMD ["/go/bin/service"]
