module otel-trace-id-propagation-example

go 1.16

require (
	github.com/labstack/echo/v4 v4.3.0
	github.com/rs/zerolog v1.22.0
	go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho v0.20.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.20.0
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/exporters/stdout v0.20.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.20.0
	go.opentelemetry.io/otel/sdk v0.20.0
	go.opentelemetry.io/otel/trace v0.20.0
)
