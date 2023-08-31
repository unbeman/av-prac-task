package config

import "github.com/caarlos0/env/v8"

const (
	AddressDefault       = "0.0.0.0:8080"
	LogLevelDefault      = "info"
	DSNDefault           = "postgresql://postgres:1211@localhost:5432/dus"
	FileDirectoryDefault = "store"
	WorkersCountDefault  = 2
	TasksSizeDefault     = 2
)

type WorkerPoolConfig struct {
	WorkersCount int
	TasksSize    int
}

func NewWorkerPoolConfig() WorkerPoolConfig {
	return WorkerPoolConfig{WorkersCount: WorkersCountDefault, TasksSize: TasksSizeDefault}
}

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
	WorkersPool   WorkerPoolConfig
}

func GetAppConfig() (AppConfig, error) {
	cfg := AppConfig{
		DB:            NewPostgresConfig(),
		Logger:        NewLoggerConfig(),
		Address:       AddressDefault,
		FileDirectory: FileDirectoryDefault,
		WorkersPool:   NewWorkerPoolConfig(),
	}

	if err := cfg.parseEnv(); err != nil {
		return cfg, err
	}
	return cfg, nil
}
