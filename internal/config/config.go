package config

import (
	"fmt"
)

const DBKindInMemory string = "IN_MEMORY"
const DBKindPostgres string = "POSTGRES"

type Config struct {
	AppName string     `json:"app_name"`
	BaseURL string     `json:"base_url"`
	HTTP    HTTPConfig `json:"http"`
	DB      DBConfig   `json:"db"`
}

var GlobalConfig Config

func (c *Config) LoadConfig() error {
	fmt.Println("Loading config...")

	c.AppName = "URL shortener"
	c.BaseURL = "http://localhost:8080"
	c.HTTP = *NewHTTPConfig(DefaultHTTPSchema, DefaultHTTPHost, DefaultHTTPPort)
	c.DB = *NewDBConfig(DBKindInMemory)

	return nil
}
