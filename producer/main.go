package main

import (
	"context"
	otelMW "go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"net/http"

	"github.com/labstack/echo/v4"

	// otel
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	// jaeger exporter
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
)

func main() {
	// https://opentelemetry.io/docs/go/getting-started/

	// otel setup
	exporter, err := jaeger.NewRawExporter(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		panic(err)
	}

	// trace provider
	ctx := context.Background()
	bsp := sdktrace.NewBatchSpanProcessor(exporter)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // should't be used in prod
	)

	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	// set global options
	otel.SetTracerProvider(tp)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	otel.SetTextMapPropagator(propagator)

	// echo setup
	e := echo.New()
	e.Use(otelMW.Middleware("producer"))

	e.GET("/", func(c echo.Context) error {
		// return result
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1324"))
}
