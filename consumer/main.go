package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	otelMW "go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"io/ioutil"
	"net/http"

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
	e.Use(otelMW.Middleware("consumer"))

	client := http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithPropagators(propagator),
		),
	}

	e.GET("/", func(c echo.Context) error {
		ctx := c.Request().Context()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, `http://localhost:1324`, nil)
		if err != nil {
			return err
		}
		// set context to propagate trace id
		req.WithContext(ctx)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return c.String(http.StatusOK, string(body))
	})
	e.Logger.Fatal(e.Start(":1323"))
}
