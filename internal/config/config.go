// Package config
/*
  Iteration 5

  Configuration params priority :

  1 - ENV vars
  2 - CLI params
  3 - Default values
*/
package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

type Config struct {
	AppName      string      `json:"app_name,omitempty"`
	ProjectStage string      `json:"project_stage,omitempty"`
	LogLevel     string      `json:"log_level,omitempty"`
	BaseURL      string      `json:"base_url,omitempty" env:"BASE_URL"`
	HTTP         *HTTPConfig `json:"http,omitempty"`
	DBKind       string      `json:"db_kind,omitempty"`
	DBDsn        string      `json:"db_dsn,omitempty" env:"DATABASE_DSN"`
	StoragePath  string      `json:"storage_path,omitempty" env:"FILE_STORAGE_PATH"`
}

// Flags
const (
	FlagAppName       = "p"
	FlagProjectStage  = "s"
	FlagLogLevel      = "l"
	FlagBaseURL       = "b"
	FlagDBKind        = "k"
	FlagHTTPInterface = "a"
	FlagDBInterface   = "d"
	FlagStoragePath   = "f"
)

// Environment variables
const (
	EnvBaseURL         = "BASE_URL"
	EnvHTTPInterface   = "SERVER_ADDR"
	EnvStorageFilename = "FILE_STORAGE_PATH"
	EnvDatabaseDSN     = "DATABASE_DSN"
)

var AppConfig *Config

func NewConfig() *Config {
	var cfg = defaultConfig()

	cfg.initFlags()

	return cfg
}

func newConfig(appName string, projectStage string, logLevel string, baseURL string, HTTP *HTTPConfig, DBKind string, DBDsn string, storagePath string) *Config {
	return &Config{
		AppName:      appName,
		ProjectStage: projectStage,
		LogLevel:     logLevel,
		BaseURL:      baseURL,
		HTTP:         HTTP,
		DBKind:       DBKind,
		DBDsn:        DBDsn,
		StoragePath:  storagePath,
	}
}

func defaultConfig() *Config {
	return newConfig(DefaultAppName, DefaultStage, DefaultLogLevel, DefaultBaseURL, DefaultHTTPConfig(), DefaultDBKind, DefaultDBDsn, DefaultStoragePath)
}

func (c *Config) LoadConfig() error {
	fmt.Println("Parse cli params")
	var err = c.loadCli()
	if err != nil {
		return err
	}

	fmt.Println("Parse env params")
	err = c.loadEnv()
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) loadCli() error {
	flag.Parse()

	fmt.Printf("Config: %+v\r\n", c)

	return nil
}

func (c *Config) loadEnv() error {
	err := env.Parse(c)
	if err != nil {
		return err
	}

	err = parseFlag(EnvHTTPInterface, c.HTTP)
	if err != nil {
		return err
	}

	fmt.Printf("Config: %+v\r\n", c)

	return nil
}

func parseFlag(env string, value flag.Value) error {
	var envVar = os.Getenv(env)
	if envVar == "" {
		return nil
	}

	fmt.Printf("[DEBUG] Config: ENV [%s] VALUE [%+v]\r\n", env, value)

	err := value.Set(envVar)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) initFlags() {
	flag.StringVar(&c.AppName, FlagAppName, DefaultAppName, "application name")
	flag.StringVar(&c.ProjectStage, FlagProjectStage, ProjectStageDevelopment, "project stage")
	flag.StringVar(&c.LogLevel, FlagLogLevel, zap.InfoLevel.CapitalString(), "log level")
	flag.StringVar(&c.BaseURL, FlagBaseURL, DefaultBaseURL, "base url")
	flag.StringVar(&c.DBKind, FlagDBKind, DefaultDBKind, "db kind")
	flag.Var(c.HTTP, FlagHTTPInterface, "http interface")
	flag.StringVar(&c.DBDsn, FlagDBInterface, DefaultDBDsn, "database dsn")
	flag.StringVar(&c.StoragePath, FlagStoragePath, DefaultStoragePath, "storage path")
}
