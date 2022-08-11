package configuration

import (
	"flag"
	"github.com/caarlos0/env"
	"log"
)

const (
	PostgresString             = "user=pavelchuykin password=1234 dbname=pavelchuykin sslmode=disable"
	GrpcPort                   = 50051
	AccessTokenLiveTimeMinutes = 15
	RefreshTokenLiveTimeDays   = 7
	AccessTokenSecret          = "jdnfksdmfksd"
	RefreshTokenSecret         = "mcmvmkmsdnfsdmfdsjf"
)

// NewConfig функция для создания нового экземпляра конфигурации
func NewConfig() *Config {
	cfg := Config{
		PostgresString:             PostgresString,
		GrpcPort:                   GrpcPort,
		AccessTokenLiveTimeMinutes: AccessTokenLiveTimeMinutes,
		RefreshTokenLiveTimeDays:   RefreshTokenLiveTimeDays,
		AccessTokenSecret:          AccessTokenSecret,
		RefreshTokenSecret:         RefreshTokenSecret,
	}

	flagGrpcPort := flag.Int("a", GrpcPort, "grpc-client port")
	flagPostgresString := flag.String("d", PostgresString, "database string")

	if *flagGrpcPort != GrpcPort {
		cfg.GrpcPort = *flagGrpcPort
	}
	if *flagPostgresString != PostgresString {
		cfg.PostgresString = *flagPostgresString
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	return &cfg
}

// Config струкутра для хранения конфигурации
type Config struct {
	PostgresString             string `env:"POSTGRES_STRING"`
	GrpcPort                   int    `env:"GRPC_PORT"`
	AccessTokenLiveTimeMinutes int    `env:"ACCESS_TOKEN_LIVE_TIME_MINUTES"`
	RefreshTokenLiveTimeDays   int    `env:"REFRESH_TOKEN_LIVE_TIME_DAYS"`
	AccessTokenSecret          string `env:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret         string `env:"REFRESH_TOKEN_SECRET"`
}
