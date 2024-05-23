package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Storage string           `yaml:"storage" env-required:"true"`
	Server  HTTPServerConfig `yaml:"http_server"`
}

type DBConfig struct {
	Username string `env:"DATABASE_USER" env-required:"true"`
	Password string `env:"DATABASE_PASSWORD" env-required:"true"`
	Database string `env:"DATABASE_NAME" env-required:"true"`
}

type HTTPServerConfig struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

// MustLoad reads config from config path and panics
// on error.
func MustLoad() (*Config, *DBConfig) {
	config, dbConfig, err := Load()
	if err != nil {
		panic(err)
	}
	return config, dbConfig
}

// Load reads config from config path. If `storage` is set to
// `postgres` it also reads DBConfig from environment variables.
// If `storage` is set to `memory` returns nil as second value.
func Load() (*Config, *DBConfig, error) {
	configPath := fetchConfigPath()
	if configPath == "" {
		return nil, nil, errors.New("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil, errors.New("config file doesn't exist: " + configPath)
	}

	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		return nil, nil, fmt.Errorf("failed to read config: %w", err)
	}

	// TODO: Maybe should also read .env, but I don't care right now.
	var dbConfig DBConfig
	if config.Storage == "postgres" {
		if err := cleanenv.ReadEnv(&dbConfig); err != nil {
			return nil, nil, fmt.Errorf("failed to read db config: %w", err)
		}

		return &config, &dbConfig, nil
	}

	return &config, nil, nil
}

// fetchConfigPath gets path to config.yaml file from --config flag
// or CONFIG_PATH environment variable. Returns empty string if neither
// are set.
func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "config file path")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
