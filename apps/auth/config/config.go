package config

import "github.com/caarlos0/env/v11"

type Config struct {
	ServerConfig   ServerConfig
	TracerConfig   TracerConfig
	DatabaseConfig PostgresConfig
	MongoConfig    MongoConfig
	SMTPConfig     SMTPConfig
	RedisConfig    RedisConfig
	S3Config       S3Config
}

type ServerConfig struct {
	Port                int    `env:"PORT" envDefault:"8080"`
	SecretKeyJWT        string `env:"SECRET_KEY_JWT,required"`
	SecretKeyRefreshJWT string `env:"SECRET_KEY_REFRESH_JWT,required"`
	UsePostgres         bool   `env:"USE_POSTGRES,required"`
	UseMongo            bool   `env:"USE_MONGO,required"`
	AppName             string `env:"APP_NAME" envDefault:"auth-service"`
	Version             string `env:"VERSION" envDefault:"0.1.0"`
}

type TracerConfig struct {
	JaegerHost string `env:"JAEGER_HOST" envDefault:"localhost"`
	JaegerPort string `env:"JAEGER_PORT" envDefault:"4318"`
	Enabled    bool   `env:"JAEGER_ENABLED" envDefault:"false"`
}

type PostgresConfig struct {
	DatabaseHost     string `env:"DB_HOST,required"`
	DatabasePort     string `env:"DB_PORT,required"`
	DatabaseUser     string `env:"DB_USER,required"`
	DatabasePassword string `env:"DB_PASSWORD,required"`
	DatabaseName     string `env:"DB_NAME,required"`
}

type MongoConfig struct {
	MongoHost     string `env:"MONGO_HOST,required"`
	MongoPort     int    `env:"MONGO_PORT,required"`
	MongoUser     string `env:"MONGO_USER,required"`
	MongoPassword string `env:"MONGO_PASSWORD,required"`
	MongoDatabase string `env:"MONGO_DATABASE,required"`
}

type SMTPConfig struct {
	SMTPHost     string `env:"SMTP_HOST"`
	SMTPPort     int    `env:"SMTP_PORT"`
	SMTPUser     string `env:"SMTP_USER"`
	SMTPPassword string `env:"SMTP_PASSWORD"`
}

type RedisConfig struct {
	RedisHost     string `env:"REDIS_HOST"`
	RedisPort     int    `env:"REDIS_PORT"`
	RedisPassword string `env:"REDIS_PASSWORD"`
}

type S3Config struct {
	S3BucketName string `env:"S3_BUCKET_NAME"`
	S3Region     string `env:"S3_REGION"`
	S3AccessKey  string `env:"S3_ACCESS_KEY"`
	S3SecretKey  string `env:"S3_SECRET_KEY"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) GetJaegerHost() string {
	return c.TracerConfig.JaegerHost
}

func (c *Config) GetJaegerPort() string {
	return c.TracerConfig.JaegerPort
}

func (c *Config) GetAppName() string {
	return c.ServerConfig.AppName
}

func (c *Config) GetVersion() string {
	return c.ServerConfig.Version
}

func (c *Config) IsEnabled() bool {
	return c.TracerConfig.Enabled
}
