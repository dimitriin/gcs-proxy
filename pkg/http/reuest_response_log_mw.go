package http

import (
	"context"
	"log/slog"
	"net/http"
)

type RequestResponseLogMiddleware struct {
	handler http.Handler
	logger  *slog.Logger
	level   slog.Level
	enabled bool
}

func NewRequestResponseLogMiddleware(
	handler http.Handler,
	logger *slog.Logger,
	level slog.Level,
	enabled bool,
) *RequestResponseLogMiddleware {
	return &RequestResponseLogMiddleware{
		handler: handler,
		logger:  logger,
		level:   level,
		enabled: enabled,
	}
}

func (m *RequestResponseLogMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	m.log(
		request.Context(),
		"request received",
		slog.String("method", request.Method),
		slog.String("url", request.URL.String()),
	)

	wrapperWriter := NewStatusRecorder(writer)

	m.handler.ServeHTTP(wrapperWriter, request)

	m.log(
		request.Context(),
		"request handled",
		slog.String("method", request.Method),
		slog.String("url", request.URL.String()),
		slog.Int("status_code", wrapperWriter.GetStatusCode()),
	)
}

func (m *RequestResponseLogMiddleware) log(ctx context.Context, msg string, args ...any) {
	m.logger.Log(ctx, m.level, msg, args...)
}
