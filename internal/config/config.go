package config

import (
	"flag"
	"fmt"
)

type Config struct {
	AppName string      `json:"app_name"`
	BaseURL string      `json:"base_url"`
	HTTP    *HTTPConfig `json:"http"`
	DBKind  string      `json:"db_kind"`
	DB      *DBConfig   `json:"db"`
}

var AppConfig *Config

func NewConfig() *Config {
	var cfg = defaultConfig()

	cfg.initFlags()

	return cfg
}

func newConfig(appName string, baseURL string, HTTP *HTTPConfig, DBKind string, DB *DBConfig) *Config {
	return &Config{
		AppName: appName,
		BaseURL: baseURL,
		HTTP:    HTTP,
		DBKind:  DBKind,
		DB:      DB,
	}
}

func defaultConfig() *Config {
	return newConfig(DefaultAppName, DefaultBaseURL, DefaultHTTPConfig(), DefaultDBKind, DefaultDBConfig())
}

func (c *Config) LoadConfig() error {
	flag.Parse()

	fmt.Printf("Config: %+v\r\n", c)

	return nil
}

func (c *Config) initFlags() {
	flag.StringVar(&c.AppName, "p", DefaultAppName, "application name")
	flag.StringVar(&c.BaseURL, "b", DefaultBaseURL, "base url")
	flag.StringVar(&c.DBKind, "db", DefaultDBKind, "db kind")
	flag.Var(c.HTTP, "a", "http interface")
	flag.Var(c.DB, "d", "db interface")
}
