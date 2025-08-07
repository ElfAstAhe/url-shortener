package config

import (
	"fmt"
)

const DBKindInMemory string = "IN_MEMORY"
const DBKindPostgres string = "POSTGRES"

type Config struct {
	AppName string       `json:"app_name"`
	HTTP    []HTTPConfig `json:"http"`
	DB      DBConfig     `json:"db"`
}

var GlobalConfig Config

func (c *Config) LoadConfig() error {
	fmt.Println("Loading config...")

	c.AppName = "URL shortener"
	c.HTTP = []HTTPConfig{*NewHTTPConfig(DefaultHTTPSchema, DefaultHTTPHost, DefaultHTTPPort)}
	c.DB = *NewDBConfig(DBKindInMemory)

	return nil
}
