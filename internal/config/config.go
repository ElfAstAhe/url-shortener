package config

import (
	"flag"
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

func init() {
	flag.StringVar(&GlobalConfig.AppName, "p", "URL shortener", "application name")
	flag.StringVar(&GlobalConfig.BaseURL, "b", "http://localhost:8080", "base url")
	flag.StringVar(&GlobalConfig.DB.Kind, "db", DBKindInMemory, "db kind")
	GlobalConfig.HTTP = *NewHTTPConfig(DefaultHTTPSchema, DefaultHTTPHost, DefaultHTTPPort)
	flag.Var(&GlobalConfig.HTTP, "a", "http interface")
}

func (c *Config) LoadConfig() error {
	fmt.Println("Loading config...")
	/*
	   c.AppName = "URL shortener"
	   c.BaseURL = "http://localhost:8080"
	   c.HTTP = *NewHTTPConfig(DefaultHTTPSchema, DefaultHTTPHost, DefaultHTTPPort)
	   c.DB = *NewDBConfig(DBKindInMemory)
	*/
	flag.Parse()

	fmt.Printf("Config: %+v\r\n", GlobalConfig)

	return nil
}
