package config

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

var ErrJSONAndFileGCSSAProvided = errors.New("either creds json or file must be provided")

type Config struct {
	Log                LogConfig
	Server             ServerConfig
	GoogleCloudStorage GoogleCloudStorageConfig `split_words:"true"`
}

type LogConfig struct {
	Level string `default:"INFO" validate:"required"`
}

type ServerConfig struct {
	Host               string
	Port               string        `default:"8787" validate:"required"`
	ReadHeaderTimeout  time.Duration `default:"5s" validate:"required" split_words:"true"`
	Routes             RoutesConfig
	RequestResponseLog RequestResponseLogConfig `split_words:"true"`
	Observability      ObservabilityConfig      `split_words:"true"`
}

type ObservabilityConfig struct {
	Metrics struct {
		Enabled   bool   `default:"true"`
		Namespace string `default:"gcs" validate:"required"`
		Subsystem string `default:"proxy" validate:"required"`
	}
}

type RequestResponseLogConfig struct {
	Enabled bool   `default:"true"`
	Level   string `default:"INFO"`
}

func (c RequestResponseLogConfig) SLogLevel() (slog.Level, error) {
	lvl := slog.Level(0)
	if err := lvl.UnmarshalText([]byte(c.Level)); err != nil {
		return lvl, fmt.Errorf("unmarshal log level, %w", err)
	}

	return lvl, nil
}

type RoutesConfig struct {
	Health  string `default:"/_health" validate:"required"`
	Metrics string `default:"/_metrics" validate:"required"`
	Proxy   string `default:"/{bucket:[0-9a-zA-Z-_.]+}/{object:.*}" validate:"required"`
}

type GoogleCloudStorageConfig struct {
	Endpoint string   `default:"https://storage.googleapis.com" validate:"required"`
	Scopes   []string `default:"https://www.googleapis.com/auth/devstorage.read_write" validate:"required"`
	Creds    GoogleCloudStorageCredsConfig
}

type GoogleCloudStorageCredsConfig struct {
	JSON string
	File string
}

func ReadFromENVAndValidate() (Config, error) {
	cfg := Config{} //nolint:exhaustruct

	if err := envconfig.Process("GCS_PROXY", &cfg); err != nil {
		return cfg, fmt.Errorf("process env config, %w", err)
	}

	if err := validator.New().Struct(&cfg); err != nil {
		return cfg, fmt.Errorf("not valid config, %w", err)
	}

	if cfg.GoogleCloudStorage.Creds.JSON != "" && cfg.GoogleCloudStorage.Creds.File != "" {
		return cfg, ErrJSONAndFileGCSSAProvided
	}

	return cfg, nil
}
