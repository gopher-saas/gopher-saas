package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
	With(fields ...zapcore.Field) Logger
	WithContext(ctx context.Context) Logger
}

type ZapLogger struct {
	logger *zap.Logger
}

func NewLogger(serviceName string, isProduction bool) Logger {
	var config zap.Config

	config = zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.FunctionKey = "function"
	config.EncoderConfig.StacktraceKey = "stacktrace"
	config.EncoderConfig.LineEnding = zapcore.DefaultLineEnding
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.EncoderConfig.EncodeName = zapcore.FullNameEncoder

	logger, err := config.Build()
	if err != nil {
		os.Exit(1)
	}

	return &ZapLogger{
		logger: logger.Named(serviceName).With(zap.String("service", serviceName)),
	}
}

func (l *ZapLogger) Debug(msg string, fields ...zapcore.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *ZapLogger) Info(msg string, fields ...zapcore.Field) {
	l.logger.Info(msg, fields...)
}

func (l *ZapLogger) Warn(msg string, fields ...zapcore.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *ZapLogger) Error(msg string, fields ...zapcore.Field) {
	l.logger.Error(msg, fields...)
}

func (l *ZapLogger) Fatal(msg string, fields ...zapcore.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *ZapLogger) With(fields ...zapcore.Field) Logger {
	return &ZapLogger{
		logger: l.logger.With(fields...),
	}
}

type contextKey string

const loggerKey contextKey = "logger"

func (l *ZapLogger) WithContext(ctx context.Context) Logger {
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return l.With(zap.String("request_id", requestID))
	}
	return l
}

func FromContext(ctx context.Context, defaultLogger Logger) Logger {
	if logger, ok := ctx.Value(loggerKey).(Logger); ok {
		return logger
	}
	return defaultLogger
}

func ContextWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
