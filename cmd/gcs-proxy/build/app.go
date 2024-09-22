package build

import (
	"log/slog"
	"net/http"

	"github.com/samber/do"

	"github.com/dimitriin/gcs-proxy/cmd/gcs-proxy/app"
)

func App(i *do.Injector) (*app.App, error) {
	return app.NewApp(
		do.MustInvoke[*http.Server](i),
		do.MustInvoke[*slog.Logger](i),
	), nil
}
