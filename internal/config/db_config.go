package config

type DBConfig struct {
	Kind     string `json:"kind"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewDBConfig(kind string) *DBConfig {
	return &DBConfig{
		Kind: kind,
	}
}
