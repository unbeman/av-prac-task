package config

import "github.com/caarlos0/env/v8"

const (
	AddressDefault       = "0.0.0.0:8080"
	LogLevelDefault      = "info"
	DSNDefault           = "postgresql://postgres:1211@localhost:5432/dus"
	FileDirectoryDefault = "store"
)

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL"`
}

func NewLoggerConfig() LoggerConfig {
	return LoggerConfig{Level: LogLevelDefault}
}

type PostgresConfig struct {
	DSN string `env:"POSTGRES_DSN"`
}

func NewPostgresConfig() PostgresConfig {
	return PostgresConfig{DSN: DSNDefault}
}

func (cfg *AppConfig) parseEnv() error {
	if err := env.Parse(cfg); err != nil {
		return err
	}
	return nil
}

type AppConfig struct {
	DB            PostgresConfig
	Logger        LoggerConfig
	Address       string `env:"RUN_ADDRESS"`
	FileDirectory string `env:"FILE_DIRECTORY"`
}

func GetAppConfig() (AppConfig, error) {
	cfg := AppConfig{
		DB:            NewPostgresConfig(),
		Logger:        NewLoggerConfig(),
		Address:       AddressDefault,
		FileDirectory: FileDirectoryDefault,
	}

	if err := cfg.parseEnv(); err != nil {
		return cfg, err
	}
	return cfg, nil
}
