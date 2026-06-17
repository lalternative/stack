// Package observability wires the skalpai sdk-go (OTLP traces + metrics +
// logs) at boot and exposes an Echo middleware that emits HTTP spans.
//
// When SKALPAI_ENDPOINT or SKALPAI_API_KEY are missing, Init returns a
// no-op shutdown and a nil error so local dev and `go test ./...` keep
// working without telemetry.
package observability

import (
	"context"

	skalpai "github.com/digstack/skalpai/packages/sdk-go"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

type Config struct {
	ServiceName string
	Endpoint    string
	APIKey      string
	ProjectID   string
}

// Init configures global OTEL providers via the skalpai SDK. The returned
// shutdown function is always safe to call (no-op if telemetry is off).
func Init(ctx context.Context, cfg Config) (func(context.Context) error, error) {
	if cfg.Endpoint == "" || cfg.APIKey == "" {
		return func(context.Context) error { return nil }, nil
	}
	return skalpai.Init(ctx, skalpai.Config{
		ServiceName: cfg.ServiceName,
		Endpoint:    cfg.Endpoint,
		APIKey:      cfg.APIKey,
	})
}

// EchoMiddleware returns the OTEL HTTP middleware bound to the service name.
func EchoMiddleware(serviceName string) echo.MiddlewareFunc {
	return otelecho.Middleware(serviceName)
}
