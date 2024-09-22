package build

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/samber/do"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"

	"github.com/dimitriin/gcs-proxy/cmd/gcs-proxy/config"
	"github.com/dimitriin/gcs-proxy/internal/observability"
	ihttp "github.com/dimitriin/gcs-proxy/pkg/http"
	iobservability "github.com/dimitriin/gcs-proxy/pkg/observability"
)

const (
	ReverseProxyHTTPHandler string = "ReverseProxyHTTPHandler"
)

func ObservableReverseProxy(i *do.Injector) (http.Handler, error) {
	proxy := do.MustInvoke[*ihttp.RequestResponseLogMiddleware](i)
	logger := do.MustInvoke[*slog.Logger](i)
	cfg := do.MustInvoke[config.ObservabilityConfig](i)

	metric := iobservability.NewHistogramObserver(cfg.Metrics.Namespace, cfg.Metrics.Subsystem)

	if err := metric.Register(); err != nil {
		return nil, fmt.Errorf("register metric, %w", err)
	}

	return ihttp.NewObserverMiddleware(
		proxy,
		iobservability.NewHistogramObserver(cfg.Metrics.Namespace, cfg.Metrics.Subsystem),
		observability.NewBucketRequestPathGeneralizer(),
		logger,
		cfg.Metrics.Enabled,
	), nil
}

func ReverseProxyWithRequestResponseLog(i *do.Injector) (*ihttp.RequestResponseLogMiddleware, error) {
	cfg := do.MustInvoke[config.RequestResponseLogConfig](i)
	logger := do.MustInvoke[*slog.Logger](i)
	proxy := do.MustInvoke[*httputil.ReverseProxy](i)

	logLevel, err := cfg.SLogLevel()
	if err != nil {
		logger.Warn(
			"can`t parse log level for request response log middleware",
			slog.String("level", cfg.Level),
			slog.String("error", err.Error()),
		)
	}

	return ihttp.NewRequestResponseLogMiddleware(
		proxy,
		logger,
		logLevel,
		cfg.Enabled,
	), nil
}

func ReverseProxy(i *do.Injector) (*httputil.ReverseProxy, error) {
	cfg := do.MustInvoke[config.GoogleCloudStorageConfig](i)

	opts := make([]option.ClientOption, 0)

	if len(cfg.Scopes) > 0 {
		opts = append(opts, option.WithScopes(cfg.Scopes...))
	}

	if cfg.Creds.JSON != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(cfg.Creds.JSON)))
	}

	if cfg.Creds.File != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.Creds.File))
	}

	googleHTTPClient, _, err := transport.NewHTTPClient(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("create http client, %w", err)
	}

	googleStorageEndpoint, err := url.Parse(cfg.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse google storage url, %w", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(googleStorageEndpoint)
	proxy.Transport = googleHTTPClient.Transport

	originalProxyDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalProxyDirector(req)

		req.Host = googleStorageEndpoint.Host
		req.Header.Set("Host", googleStorageEndpoint.Host)
	}

	return proxy, nil
}
