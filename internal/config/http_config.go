package config

import "strconv"

type HTTPConfig struct {
	Schema string `json:"schema"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
}

func NewHTTPConfig(schema string, host string, port int) *HTTPConfig {
	return &HTTPConfig{
		Schema: schema,
		Host:   host,
		Port:   port,
	}
}

func (hc *HTTPConfig) GetHost() string {
	return hc.Host + ":" + strconv.Itoa(hc.Port)
}
