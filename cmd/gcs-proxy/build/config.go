package build

import (
	"fmt"

	"github.com/samber/do"

	"github.com/dimitriin/gcs-proxy/cmd/gcs-proxy/config"
)

func Config(_ *do.Injector) (config.Config, error) {
	cfg, err := config.ReadFromENVAndValidate()
	if err != nil {
		return config.Config{}, fmt.Errorf("read config, %w", err)
	}

	return cfg, nil
}

func LogConfig(i *do.Injector) (config.LogConfig, error) {
	return do.MustInvoke[config.Config](i).Log, nil
}

func ServerConfig(i *do.Injector) (config.ServerConfig, error) {
	return do.MustInvoke[config.Config](i).Server, nil
}

func RequestResponseLogMwConfig(i *do.Injector) (config.RequestResponseLogConfig, error) {
	return do.MustInvoke[config.ServerConfig](i).RequestResponseLog, nil
}

func GCSConfig(i *do.Injector) (config.GoogleCloudStorageConfig, error) {
	return do.MustInvoke[config.Config](i).GoogleCloudStorage, nil
}

func RequestObservabilityConfig(i *do.Injector) (config.ObservabilityConfig, error) {
	return do.MustInvoke[config.ServerConfig](i).Observability, nil
}
