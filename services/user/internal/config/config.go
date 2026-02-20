package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string     `yaml: "env" env-default:"local"`
	DatabaseURL    string     `yaml: "db_url" env-default:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	RedisAddr      string     `yaml: "redis_addr" env-default:"localhost:6379"`
	MigrationsPath string     `yaml: "migrations_path" env-default:"file://migrations"`
	HTTPAddr       HTTPServer `yaml: "http_addr" env-default:"localhost:8080"`
}

type HTTPServer struct {
	Address     string        `yaml: address" env-default: "localhost:8080"`
	Timeout     time.Duration `yaml: "timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml: "idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	cfg := Config{}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal("Failed to read config: ", err)
	}

	return &cfg
}
