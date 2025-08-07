package config

type DbConfig struct {
	Kind     string `json:"kind"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewDbConfig(kind string) *DbConfig {
	return &DbConfig{
		Kind: kind,
	}
}
