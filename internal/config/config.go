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
)

type Config struct {
	AppName string      `json:"app_name"`
	BaseURL string      `json:"base_url" env:"BASE_URL"`
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

	err = parseFlag("SERVER_ADDRESS", c.HTTP)
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

	err := value.Set(envVar)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) initFlags() {
	flag.StringVar(&c.AppName, "p", DefaultAppName, "application name")
	flag.StringVar(&c.BaseURL, "b", DefaultBaseURL, "base url")
	flag.StringVar(&c.DBKind, "db", DefaultDBKind, "db kind")
	flag.Var(c.HTTP, "a", "http interface")
	flag.Var(c.DB, "d", "db interface")
}
