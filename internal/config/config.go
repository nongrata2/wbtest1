package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	HttpServerAddress string        `env:"HTTP_SERVER_ADDRESS" env-default:"localhost:8081"`
	HttpServerTimeout time.Duration `env:"HTTP_SERVER_TIMEOUT" env-default:"5s"`
	LogLevel          string        `env:"LOG_LEVEL" env-default:"DEBUG"`
	DBHost            string        `env:"DB_HOST" env-default:"db"`
	DBUser            string        `env:"POSTGRES_USER" env-default:"postgres"`
	DBPassword        string        `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	DBName            string        `env:"POSTGRES_NAME" env-default:"postgres"`
	DBPort            string        `env:"POSTGRES_PORT" env-default:"5432"`
	KafkaBrokers      string        `env:"KAFKA_BROKERS" env-required:"true"`
	KafkaTopic        string        `env:"KAFKA_TOPIC" env-required:"true"`
	KafkaGroupID      string        `env:"KAFKA_GROUP_ID" env-required:"true"`
}

func MustLoadCfg(configPath string) Config {
	if err := godotenv.Load(configPath); err != nil {
		log.Fatalf("failed to load .env file: %s", err)
	}

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("failed to read environment variables: %s", err)
	}

	return cfg
}
