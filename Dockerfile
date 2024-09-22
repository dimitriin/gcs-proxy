FROM golang:1.22 AS build

COPY . .

RUN make build

FROM gcr.io/distroless/base

COPY --from=build /go/bin/gcs-proxy /gcs-proxy

CMD ["/gcs-proxy"]