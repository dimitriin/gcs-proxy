# GCS proxy

Reverse proxy with authentication for Google Cloud Storage (GCS). 
The proxy provides access to private GCS buckets, making it ideal for use as sidecar container in Kubernetes.

## Table of contents
1. [Configuration](#configuration)
2. [Running](#running)
3. [Build binary from source code](#build-binary-from-source-code)
4. [Run from binary](#run-from-binary)
5. [Contributing](#contributing)
6. [License](#license)

## Configuration
The proxy could be configured with environment variables. 
By default, the proxy tries to find default credentials to GCS. 
See https://cloud.google.com/docs/authentication/external/set-up-adc for more information.

Alternatively, custom credentials could be set with one of next environment variables:
* `GCS_PROXY_GOOGLE_CLOUD_STORAGE_CREDS_JSON` - JSON string with GCS service account credentials;
* `GCS_PROXY_GOOGLE_CLOUD_STORAGE_CREDS_FILE` - path to JSON file with GCS service account credentials.

For more configuration options see [Advanced configuration](#advanced-configuration) section.

## Running
Run the proxy with docker by the following command:
```bash
docker run \
  -p 8787:8787 \
  -v ${HOST_PATH_TO_SERVICE_ACCOUNT_JSON_WITH_ACCESS_TO_GCS}:/service_account.json 
  -e GCS_PROXY_GOOGLE_CLOUD_STORAGE_CREDS_PATH /service_account.json \ 
  ghcr.io/dimitriin/gcs-proxy:v1.0.0
```
Prebuilt docker image `ghcr.io/dimitriin/gcs-proxy:${RELEASE_TAG}` could be found at [GitHub Container Registry](https://github.com/dimitriin/gcs-proxy/pkgs/container/gcs-proxy).


Then access to GCS bucket objects with:
```
GET http://localhost:8787/${BUCKET_NAME}/${OBJECT_NAME}
```

Also, write operations provided by [XML-API](https://cloud.google.com/storage/docs/xml-api/overview) are available, 
but do not forget to [set proper scopes to the service account](https://cloud.google.com/storage/docs/oauth-scopes).

## Build binary from source code
Run make command to build binary from source code:
```bash
make build
```
Optional command configuration environment variables:
* `GOOS` (default `linux`);
* `GOARCH` (default `amd64`);
* `BIN_PATH` (default `./bin`).

## Run from binary

Run the proxy with the following command:
```bash
GCS_PROXY_GOOGLE_CLOUD_STORAGE_CREDS_PATH=${HOST_PATH_TO_SERVICE_ACCOUNT_JSON_WITH_ACCESS_TO_GCS} \
./bin/gcs-proxy-${GOOS}-${GOARCH}
```

### Build docker image from source code
Build docker image from source code:
```bash
make build-image
```
Optional command configuration env variables:
* `GCS_PROXY_DOCKER_IMG_REPO` (default `ghcr.io/dimitriin/gcs-proxy`);
* `GCS_PROXY_DOCKER_IMG_TAG` (default `latest`).

### Advanced configuration

Advanced configuration environment variables:

| Variable | Description | Default |
| -------- | ----------- | ------- |
| `GCS_PROXY_LOG_LEVEL` | Log level | `INFO` |
| `GCS_PROXY_SERVER_HOST` | Proxy server host | `localhost` |
| `GCS_PROXY_SERVER_PORT` | Proxy server port | `8787` |
| `GCS_PROXY_SERVER_READ_HEADER_TIMEOUT` | Read header timeout | `5s` |
| `GCS_PROXY_SERVER_ROUTES_PROXY` | Route proxied to GCS | `/{bucket:[0-9a-zA-Z-_.]+}/{object:.*}` |
| `GCS_PROXY_SERVER_ROUTES_HEALTH` | Health check route | `/_health` |
| `GCS_PROXY_SERVER_ROUTES_METRICS` | Metrics route | `/_metrics` |
| `GCS_PROXY_SERVER_REQUEST_RESPONSE_LOG_ENABLED` | Enable request/response logging | `true` |
| `GCS_PROXY_SERVER_REQUEST_RESPONSE_LOG_LEVEL` | Request/response log level | `INFO` |
| `GCS_PROXY_SERVER_OBSERVABILITY_METRICS_ENABLED` | Enable proxy metrics | `true` |
| `GCS_PROXY_SERVER_OBSERVABILITY_METRICS_NAMESPACE` | Prometheus metrics namespace | `gcs` |
| `GCS_PROXY_SERVER_OBSERVABILITY_METRICS_SUBSYSTEM` | Prometheus metrics subsystem | `proxy` |
| `GCS_PROXY_GOOGLE_CLOUD_STORAGE_ENDPOINT` | GCS endpoint | `https://storage.googleapis.com` |
| `GCS_PROXY_GOOGLE_CLOUD_STORAGE_SCOPES` | GCS scopes | `https://www.googleapis.com/auth/devstorage.read_write` |

## Contributing

Feel free to submit [issues](https://github.com/dimitriin/gcs-proxy/issues) or [pull requests](https://github.com/dimitriin/gcs-proxy/pulls).

## License

GCS proxy is licensed under the MIT License. See the [LICENSE](./LICENSE) for more details.

