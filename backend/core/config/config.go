package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Database Database
	Server   Server
	Logger   Logger
}

type Database struct {
	Host   string
	Port   int
	DBName string `toml:"db_name"`
	DBUser string `toml:"db_user"`
}

type Server struct {
	Host string
	Port int
	Mode string
}

type Logger struct {
	Level string
}

func NewConfig(path string) (*Config, error) {
	cfg := &Config{}
	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
