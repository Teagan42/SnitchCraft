package metrics

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/teagan42/snitchcraft/internal/models"
)

func TestSetupTracing_SetsTracerProvider(t *testing.T) {
	// Save the original provider to restore after test
	origProvider := otel.GetTracerProvider()
	defer otel.SetTracerProvider(origProvider)

	SetupTracing()

	tp := otel.GetTracerProvider()
	if tp == nil {
		t.Fatal("TracerProvider is nil after setupTracing")
	}

	// Check that the provider is of type *trace.TracerProvider
	_, ok := tp.(*trace.TracerProvider)
	if !ok {
		t.Errorf("TracerProvider is not of type *trace.TracerProvider, got %T", tp)
	}
}

func TestOpenTelemetryMetrics_Start_HeuristicCounters(t *testing.T) {
	origProvider := otel.GetMeterProvider()
	defer otel.SetMeterProvider(origProvider)
	otel.SetMeterProvider(sdkmetric.NewMeterProvider())
	var meter = otel.GetMeterProvider().Meter("test")
	requestCounter, err := meter.Int64Counter(
		"snitchcraft_http_requests_total",
		metric.WithDescription("Total HTTP requests received"),
	)
	if err != nil {
		panic(err)
	}
	histogram, err := meter.Int64Histogram(
		"snitchcraft_http_request_duration_milliseconds",
		metric.WithDescription("Duration of HTTP request processing prior to proxy in nanoseconds"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		panic(err)
	}
	metrics := &OpenTelemetryMetrics{
		meter:                meter,
		heuristicCounters:    make(map[string]metric.Int64Counter),
		requestCounter:       requestCounter,
		requestTimeHistogram: histogram, // Not used in this test
	}

	resultChan := make(chan models.RequestResult, 1)
	err = metrics.Start(resultChan)
	if err != nil {
		t.Fatalf("Start returned error: %v", err)
	}

	resultChan <- models.RequestResult{
		Request: &http.Request{
			Method: "POST",
			URL:    &url.URL{Path: "/foo"},
			Response: &http.Response{
				StatusCode: 200,
			},
		},
		HeuristicResults: []models.HeuristicResult{
			{Name: "testheuristic", Issue: "something"},
		},
		Duration: uint64(100 * time.Nanosecond),
	}

	time.Sleep(200 * time.Millisecond)

	if len(metrics.heuristicCounters) == 0 {
		t.Error("Heuristic counters map is empty")
	}
	if _, ok := metrics.heuristicCounters["testheuristic"]; !ok {
		t.Error("Heuristic counter for 'testheuristic' was not created")
	}
}
