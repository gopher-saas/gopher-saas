package tracer

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

var (
	tracer     *sdktrace.TracerProvider
	mainTracer trace.Tracer
	isEnabled  bool
)

func Initialize(cfg TracerConfig) error {
	// Si le tracing n'est pas activé, définir un no-op tracer
	if !cfg.IsEnabled() {
		isEnabled = false
		slog.Info("tracing disabled, using no-op tracer")
		mainTracer = noop.NewTracerProvider().Tracer("no-op")
		return nil
	}

	isEnabled = true
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(cfg.GetJaegerHost()+":"+cfg.GetJaegerPort()),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(
			exporter,
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
			sdktrace.WithBatchTimeout(sdktrace.DefaultScheduleDelay*time.Millisecond),
		),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(cfg.GetAppName()),
				semconv.ServiceVersionKey.String(cfg.GetVersion()),
			),
		),
	)

	otel.SetTracerProvider(tracerProvider)
	tracer = tracerProvider
	mainTracer = tracerProvider.Tracer(cfg.GetAppName())
	slog.Info("tracing initialized", "app", cfg.GetAppName(), "jaeger_endpoint", cfg.GetJaegerHost()+":"+cfg.GetJaegerPort())
	return nil
}

// StartSpan starts a new span with the given context
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	// Si le tracing n'est pas activé, retourner un span noop
	if !isEnabled {
		return ctx, trace.SpanFromContext(ctx)
	}

	// Ajoute le request ID au span si disponible
	if requestID, ok := ctx.Value(middleware.RequestIDKey).(string); ok {
		opts = append(opts, trace.WithAttributes(
			attribute.String("request_id", requestID),
		))
	}
	return mainTracer.Start(ctx, name, opts...)
}

// Close properly shuts down the tracer
func Close(ctx context.Context) error {
	if !isEnabled || tracer == nil {
		return nil
	}
	return tracer.Shutdown(ctx)
}

// WithSpan is a helper to wrap a function with a span
func WithSpan(ctx context.Context, name string, fn func(context.Context) error) error {
	ctx, span := StartSpan(ctx, name)
	defer span.End()

	if err := fn(ctx); err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}
