package config

import "strconv"

type HttpConfig struct {
	Schema string `json:"schema"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
}

func NewHttpConfig(schema string, host string, port int) *HttpConfig {
	return &HttpConfig{
		Schema: schema,
		Host:   host,
		Port:   port,
	}
}

func (hc *HttpConfig) GetHost() string {
	return hc.Host + ":" + strconv.Itoa(hc.Port)
}
