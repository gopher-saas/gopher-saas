package tracer

type TracerConfig interface {
	GetJaegerHost() string
	GetJaegerPort() string
	GetAppName() string
	GetVersion() string
	IsEnabled() bool
}
