FROM golang:alpine AS builder-stage

ENV CGO_ENABLED=0 \
    GOOS=linux

RUN apk update && apk add --no-cache git ca-certificates tzdata

WORKDIR /build

COPY . .

#RUN go mod download
RUN go build -ldflags="-s -w" -mod=mod -o /build/altor cmd/main/main.go

FROM scratch
ENV DEBUG=0 \
    LOG_LEVEL="ERROR" \
    ANILIBRIA_WS_SCHEME="wss" \
    ANILIBRIA_WS_HOST="api.anilibria.tv" \
    ANILIBRIA_WS_PATH="v2/ws/" \
    QBT_SCHEME="http" \
    QBT_HOST="127.0.0.1:8080" \
    QBT_USER="admin" \
    QBT_PASSWORD="" \
    QBT_SAVE_FOLDER="" \
    QBT_IGNORE_TLS_VERIFY=0 \
    QBT_CATEGORY="altor-bot"

COPY --from=builder-stage /build/altor /altor
COPY --from=builder-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder-stage /usr/share/zoneinfo /usr/share/zoneinfo/

ENTRYPOINT ["/altor"]
