package build

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/samber/do"

	"github.com/dimitriin/gcs-proxy/cmd/gcs-proxy/config"
)

func ProxyServer(i *do.Injector) (*http.Server, error) {
	cfg := do.MustInvoke[config.ServerConfig](i)
	proxy := do.MustInvokeNamed[http.Handler](i, ReverseProxyHTTPHandler)

	router := mux.NewRouter()

	router.HandleFunc(cfg.Routes.Health, func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}).Methods("GET", "HEAD")

	router.Handle(cfg.Routes.Metrics, promhttp.Handler()).Methods("GET", "HEAD")

	router.HandleFunc(cfg.Routes.Proxy, proxy.ServeHTTP)

	//nolint:exhaustruct
	return &http.Server{
		Addr:              fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler:           router,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}, nil
}
