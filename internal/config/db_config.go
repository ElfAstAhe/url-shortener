package config

import (
	"fmt"
	"strconv"
	"strings"
)

const DBKindInMemory string = "IN_MEMORY"
const DBKindPostgres string = "POSTGRES"

type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewDBConfig(host string, port int, database string, username string, password string) *DBConfig {
	return &DBConfig{
		Host:     host,
		Port:     port,
		Database: database,
		Username: username,
		Password: password,
	}
}

func DefaultDBConfig() *DBConfig {
	return NewDBConfig(DefaultDBHost, DefaultDBPort, DefaultDBDatabase, DefaultDBUsername, DefaultDBPassword)
}

// flag.Value ==================================
func (DB *DBConfig) String() string {
	return fmt.Sprintf("%s:%v:%s:%s:%s", DB.Host, DB.Port, DB.Database, DB.Username, DB.Password)
}

func (DB *DBConfig) Set(s string) error {
	params := strings.Split(s, ":")
	if len(s) < 5 {
		return fmt.Errorf("invalid db config format: %s, example: host:port:database:username:password", s)
	}

	DB.Host = params[0]
	DB.Port, _ = strconv.Atoi(params[1])
	DB.Database = params[2]
	DB.Username = params[3]
	DB.Password = params[4]

	return nil
}

// =============================================
