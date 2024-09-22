ARG GOOS=linux
ARG GOARCH=amd64

FROM golang:1.22 AS build

ARG GOOS
ARG GOARCH

COPY . .

RUN GOOS=${GOOS} GARCH=${GOARCH} make build

RUN ls -lah ./bin

FROM gcr.io/distroless/base

ARG GOOS
ARG GOARCH

COPY --from=build /go/bin/gcs-proxy-${GOOS}-${GOARCH} /gcs-proxy

CMD ["/gcs-proxy"]