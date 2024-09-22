# GCS Proxy
Reverse proxy with authentication for Google Cloud Storage (GCS).
It allows to provide access to private GCS buckets, for example, use it as sidecar container in Kubernetes.

## Build

### Binary from source code
Build binary from source code:
```bash
make build
```
Optional command configuration env variables:
* `GOOS` (default `linux`);
* `GOARCH` (default `amd64`);
* `BIN_PATH` (default `./bin`).

### Docker image from source code
Build docker image from source code:
```bash
make build-image
```
Optional command configuration env variables:
* `GCS_PROXY_DOCKER_IMG_REPO` (default `ghcr.io/dimitriin/gcs-proxy`);
* `GCS_PROXY_DOCKER_IMG_TAG` (default `latest`).

### Prebuilt docker image
Prebuilt docker image could be found at GitHub Container Registry: 
`ghcr.io/dimitriin/gcs-proxy:${RELEASE_TAG}`, for example, `ghcr.io/dimitriin/gcs-proxy:v1.0.0`.

### Minimal configuration
By default, gcs-proxy tries to find default credentials to GCS. See https://cloud.google.com/docs/authentication/external/set-up-adc for more information.

Alternatively, custom credentials could be set with one of next environment variables:
* `GCS_PROXY_GOOGLE_CLOUD_STORAGE_CREDS_JSON` - JSON with GCS service account credentials;
* `GCS_PROXY_GOOGLE_CLOUD_STORAGE_CREDS_FILE` - path to file with GCS service account credentials.

### Advanced configuration
Also, the following environment variables could be set to change default configuration:
* `GCS_PROXY_LOG_LEVEL` - log level (default: `INFO`);
* `GCS_PROXY_SERVER_HOST` - host to bind proxy server (default: `localhost`);
* `GCS_PROXY_SERVER_PORT` - port to bind proxy server (default: `8787`);
* `GCS_PROXY_SERVER_READ_HEADER_TIMEOUT` - read header timeout for proxy server (default: `5s`);
* `GCS_PROXY_SERVER_ROUTES_PROXY` - route to be proxied to GCS (default: `/{bucket:[0-9a-zA-Z-_.]+}/{object:.*}`);
* `GCS_PROXY_SERVER_ROUTES_HEALTH` - healthcheck route (default: `/_health`);
* `GCS_PROXY_SERVER_ROUTES_METRICS` - metrics route (default: `/_metrics`);
* `GCS_PROXY_SERVER_REQUEST_RESPONSE_LOG_ENABLED` - enable request/response logging (default: `true`);
* `GCS_PROXY_SERVER_REQUEST_RESPONSE_LOG_LEVEL` - request/response log level (default: `INFO`);
* `GCS_PROXY_SERVER_OBSERVABILITY_METRICS_ENABLED` - enable middleware to collect additional proxy metrics (default: `true`);
* `GCS_PROXY_SERVER_OBSERVABILITY_METRICS_NAMESPACE` - prometheus namespace for metrics (default: `gcs`);
* `GCS_PROXY_SERVER_OBSERVABILITY_METRICS_SUBSYSTEM` - prometheus subsystem for metrics (default: `proxy`);
* `GCS_PROXY_GOOGLE_CLOUD_STORAGE_ENDPOINT` - GCS endpoint (default: `https://storage.googleapis.com`);
* `GCS_PROXY_GOOGLE_CLOUD_STORAGE_SCOPES` - GCS scopes (default: `https://www.googleapis.com/auth/devstorage.read_write`).