# Build stage
FROM golang:1.16.7-buster AS build-env
RUN mkdir -p /go/src/ssr-golang
WORKDIR /go/src/ssr-golang
COPY . .
RUN GO111MODULE=on go mod vendor && CGO_ENABLED=0 GOOS=linux go build -o goserver ./cmd/dev-server/dev-server.go

FROM debian:buster-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY --from=build-env /go/src/ssr-golang/goserver /app/goserver
COPY ./web /app/web/
WORKDIR /app
EXPOSE 8080
ENTRYPOINT ["/app/goserver"]
