package config

import (
	"fmt"
	"strconv"
	"strings"
)

type HTTPConfig struct {
	Schema string `json:"schema"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
}

func (hc *HTTPConfig) String() string {
	return fmt.Sprintf("%s:%v", hc.Host, hc.Port)
}

func (hc *HTTPConfig) Set(s string) error {
	temp := s
	params := strings.Split(temp, ":")
	if len(params) != 2 {
		temp = "localhost:8080"
		params = strings.Split(temp, ":")
	}
	hc.Host = params[0]
	var err error
	hc.Port, err = strconv.Atoi(params[1])
	if err != nil {
		return err
	}

	return nil
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
