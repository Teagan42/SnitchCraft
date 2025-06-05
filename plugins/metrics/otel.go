package metrics

import (
	"context"
	"fmt"

	"github.com/teagan42/snitchcraft/internal/interfaces"
	"github.com/teagan42/snitchcraft/internal/models"
	"github.com/teagan42/snitchcraft/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

type OpenTelemetryMetrics struct {
	meter                metric.Meter
	requestCounter       metric.Int64Counter
	heuristicCounters    map[string]metric.Int64Counter
	requestTimeHistogram metric.Int64Histogram
}

func SetupTracing() {
	exp, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.Empty()),
	)
	otel.SetTracerProvider(tp)
}

func NewOpenTelemetryMetrics(cfg models.Config) interfaces.MetricsPlugin {
	if cfg.OTELMetricUrl == "" {
		fmt.Println("[metrics] OTEL_METRIC_URL not set, skipping OpenTelemetry initialization")
		return nil
	}

	ctx := context.Background()
	exporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint(cfg.OTELMetricUrl),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		fmt.Printf("[metrics] Failed to create OTEL metric exporter: %v\n", err)
		return nil
	}

	reader := sdkmetric.NewPeriodicReader(exporter)
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	otel.SetMeterProvider(provider)

	meter := provider.Meter("snitchcraft")

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

	fmt.Println("[metrics] OpenTelemetryMetrics initialized")
	return &OpenTelemetryMetrics{
		meter:                meter,
		requestCounter:       requestCounter,
		heuristicCounters:    make(map[string]metric.Int64Counter),
		requestTimeHistogram: histogram,
	}
}

func (m *OpenTelemetryMetrics) Name() string {
	return "OpenTelemetryMetrics"
}

func (m *OpenTelemetryMetrics) Start(resultChannel chan models.RequestResult) error {
	go func() {
		for result := range resultChannel {
			m.requestCounter.Add(context.Background(), 1,
				metric.WithAttributes(
					attribute.String("method", result.Request.Method),
					attribute.String("path", result.Request.URL.Path),
				),
			)
			var atributes = append([]attribute.KeyValue{
				attribute.String("method", result.Request.Method),
				attribute.String("path", result.Request.URL.Path)},
				utils.Map(result.HeuristicResults, func(r models.HeuristicResult) attribute.KeyValue {
					return attribute.Bool(fmt.Sprintf("heuristic_%s", r.Name), r.Issue != "")
				})...,
			)
			m.requestTimeHistogram.Record(context.Background(), int64(result.Duration),
				metric.WithAttributes(
					atributes...,
				),
			)

			for _, heuristicResult := range utils.Filter(result.HeuristicResults, func(r models.HeuristicResult) bool {
				return r.Issue != ""
			}) {
				counter, ok := m.heuristicCounters[heuristicResult.Name]
				if !ok {
					newCounter, err := m.meter.Int64Counter(
						"snitchcraft_heuristic_"+heuristicResult.Name+"_hits",
						metric.WithDescription("Heuristic matches"),
					)
					if err != nil {
						fmt.Printf("[metrics] Failed to create counter for %s: %v\n", heuristicResult.Name, err)
						continue
					}
					m.heuristicCounters[heuristicResult.Name] = newCounter
					counter = newCounter
				}
				counter.Add(context.Background(), 1,
					metric.WithAttributes(
						[]attribute.KeyValue{
							attribute.String("method", result.Request.Method),
							attribute.String("path", result.Request.URL.Path),
							attribute.String("heuristic", heuristicResult.Name),
						}...,
					),
				)
			}
		}
	}()
	fmt.Println("[metrics] OpenTelemetryMetrics started")
	return nil
}

func init() {
	RegisterMetricsPlugin(NewOpenTelemetryMetrics)
	fmt.Println("[metrics] OpenTelemetry metrics plugin registered")
}
