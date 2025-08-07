package config

import (
	"fmt"
)

const DbKindInMemory string = "IN_MEMORY"
const DbKindPostgres string = "POSTGRES"

type Config struct {
	AppName string       `json:"app_name"`
	Http    []HttpConfig `json:"http"`
	Db      DbConfig     `json:"db"`
}

var GlobalConfig Config

func (c *Config) LoadConfig() error {
	fmt.Println("Loading config...")

	c.AppName = "URL shortener"
	c.Http = []HttpConfig{*NewHttpConfig(DefaultHttpSchema, DefaultHttpHost, DefaultHttpPort)}
	c.Db = *NewDbConfig(DbKindInMemory)

	return nil
}
