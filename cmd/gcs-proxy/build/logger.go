package build

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/samber/do"

	"github.com/dimitriin/gcs-proxy/cmd/gcs-proxy/config"
)

func Logger(i *do.Injector) (*slog.Logger, error) {
	cfg := do.MustInvoke[config.LogConfig](i)

	lvl := slog.Level(0)
	if err := lvl.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level, %w", err)
	}

	logger := slog.New(
		slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     lvl,
			ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
				if attr.Key == slog.LevelKey {
					return slog.String(attr.Key, strings.ToLower(attr.Value.String()))
				}

				if attr.Key == slog.TimeKey {
					return slog.String("app-time", attr.Value.Time().Format("2006-01-02 15:04:05.00Z0700"))
				}

				if attr.Value.Kind() == slog.KindDuration {
					return slog.Float64(attr.Key, float64(attr.Value.Duration())/float64(time.Second))
				}

				return attr
			},
		}),
	)

	return logger, nil
}
