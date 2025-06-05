package metrics

import (
	"context"
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
func TestOpenTelemetryMetrics_Start_IncrementsRequestCounter(t *testing.T) {
	// Save and restore original MeterProvider
	origProvider := otel.GetMeterProvider()
	defer otel.SetMeterProvider(origProvider)

	// Use a no-op MeterProvider to avoid exporting metrics
	otel.SetMeterProvider(sdkmetric.NewMeterProvider())

	metrics := &OpenTelemetryMetrics{
		meter:             otel.GetMeterProvider().Meter("test"),
		heuristicCounters: make(map[string]metric.Int64Counter),
	}

	// Create a fake counter to observe Add calls
	metrics.requestCounter, _ = metrics.meter.Int64Counter("test_counter")

	resultChan := make(chan models.RequestResult, 1)
	err := metrics.Start(resultChan)
	if err != nil {
		t.Fatalf("Start returned error: %v", err)
	}

	var addCalled = false
	go func() {
		for {
			select {
			case result := <-resultChan:
				if result.Request == nil {
					t.Fatal("Received nil Request in result")
				} else {
					addCalled = true
				}
			}
		}
	}()

	resultChan <- models.RequestResult{
		Request: &http.Request{
			Method: "GET",
			URL:    &url.URL{Path: "/test"},
		},
	}

	// Wait for goroutine to process
	time.Sleep(200 * time.Millisecond)

	if !addCalled {
		t.Error("Request counter Add was not called")
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
	metrics := &OpenTelemetryMetrics{
		meter:             meter,
		heuristicCounters: make(map[string]metric.Int64Counter),
		requestCounter:    requestCounter,
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
		},
		HeuristicResults: []models.HeuristicResult{
			{Name: "testheuristic", Issue: "something"},
		},
	}

	time.Sleep(200 * time.Millisecond)

	if len(metrics.heuristicCounters) == 0 {
		t.Error("Heuristic counters map is empty")
	}
	if _, ok := metrics.heuristicCounters["testheuristic"]; !ok {
		t.Error("Heuristic counter for 'testheuristic' was not created")
	}
}

// fakeInt64Counter implements metric.Int64Counter for testing
type fakeInt64Counter struct {
	addFunc func(ctx context.Context, value int64, opts ...metric.AddOption)
}

func (f *fakeInt64Counter) Add(ctx context.Context, value int64, opts ...metric.AddOption) {
	if f.addFunc != nil {
		f.addFunc(ctx, value, opts...)
	}
}

func (f *fakeInt64Counter) int64Counter() {}
