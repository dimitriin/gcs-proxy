package main

import (
	"github.com/samber/do"

	"github.com/dimitriin/gcs-proxy/cmd/gcs-proxy/app"
	"github.com/dimitriin/gcs-proxy/cmd/gcs-proxy/build"
)

func main() {
	injector := do.New()

	do.Provide(injector, build.Config)
	do.Provide(injector, build.LogConfig)
	do.Provide(injector, build.ServerConfig)
	do.Provide(injector, build.RequestResponseLogMwConfig)
	do.Provide(injector, build.RequestObservabilityConfig)
	do.Provide(injector, build.GCSConfig)

	do.Provide(injector, build.Logger)

	do.Provide(injector, build.ReverseProxy)
	do.Provide(injector, build.ReverseProxyWithRequestResponseLog)
	do.ProvideNamed(injector, build.ReverseProxyHTTPHandler, build.ObservableReverseProxy)

	do.Provide(injector, build.ProxyServer)
	do.Provide(injector, build.App)

	do.MustInvoke[*app.App](injector).Run()
}
