package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string     `yaml:"env" env-default:"local"`
	DatabaseURL    string     `yaml:"db_url" env-default:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	RedisAddr      string     `yaml:"redis_addr" env-default:"localhost:6379"`
	MigrationsPath string     `yaml:"migrations_path" env-default:"file://migrations"`
	HTTPAddr       HTTPServer `yaml:"http_addr"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath, ok := os.LookupEnv("CONFIG_PATH")
	if !ok || configPath == "" {
		configPath = "./services/user/config/local.yaml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", configPath)
	}

	cfg := Config{}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal(err)
	}
	return &cfg
}
