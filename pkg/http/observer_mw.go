package http

import (
	"log/slog"
	"net/http"
	"time"
)

type observer interface {
	Observe(
		duration time.Duration,
		requestMethod string,
		requestPath string,
		responseCode int,
	) error
}

type requestPathGeneralizer interface {
	GetRequestGeneralizedPath(r *http.Request) string
}

type ObserverMiddleware struct {
	handler                http.Handler
	observer               observer
	requestPathGeneralizer requestPathGeneralizer
	logger                 *slog.Logger
	enabled                bool
}

func NewObserverMiddleware(
	handler http.Handler,
	observer observer,
	requestPathGeneralizer requestPathGeneralizer,
	logger *slog.Logger,
	enabled bool,
) *ObserverMiddleware {
	return &ObserverMiddleware{
		handler:                handler,
		observer:               observer,
		requestPathGeneralizer: requestPathGeneralizer,
		logger:                 logger,
		enabled:                enabled,
	}
}

func (m *ObserverMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	startTime := time.Now()

	wrapperWriter := NewStatusRecorder(writer)

	defer func() {
		if err := m.observer.Observe(
			time.Since(startTime),
			request.Method,
			m.requestPathGeneralizer.GetRequestGeneralizedPath(request),
			wrapperWriter.statusCode,
		); err != nil {
			m.logger.Warn("can not observe request", slog.String("error", err.Error()))
		}
	}()

	m.handler.ServeHTTP(wrapperWriter, request)
}
