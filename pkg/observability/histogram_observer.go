package observability

import (
	"fmt"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type HistogramObserver struct {
	histogramVec *prometheus.HistogramVec
}

func NewHistogramObserver(namespace string, subsystem string) *HistogramObserver {
	//nolint:exhaustruct
	histogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "http_request_duration_seconds",
		Help:      "Http request duration seconds",
	}, []string{
		"request_method",
		"request_path",
		"response_code",
	})

	return &HistogramObserver{
		histogramVec: histogramVec,
	}
}

func (h *HistogramObserver) Register() error {
	if err := prometheus.Register(h.histogramVec); err != nil {
		return fmt.Errorf("can not register http_request_duration_seconds metric, %w", err)
	}

	return nil
}

func (h *HistogramObserver) Observe(
	duration time.Duration,
	requestMethod string,
	requestPath string,
	responseCode int,
) error {
	metric, err := h.histogramVec.GetMetricWith(map[string]string{
		"request_method": requestMethod,
		"request_path":   requestPath,
		"response_code":  strconv.Itoa(responseCode),
	})
	if err != nil {
		return fmt.Errorf("can not get metric with labels, %w", err)
	}

	metric.Observe(duration.Seconds())

	return nil
}
