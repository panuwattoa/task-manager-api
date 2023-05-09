package config

import (
	"os"
	"time"
)

var Conf = new(Config)

var (
	ServerType    = GetEnv("SERVER_TYPE", "Development")
	MongoHost     = GetEnv("MONGO_HOST", "127.0.0.1")
	MongoDBName   = GetEnv("MONGO_DBNAME", "")
	MongoUser     = GetEnv("MONGO_USERNAME", "")
	MongoPassword = GetEnv("MONGO_PASSWORD", "")
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type Config struct {
	Server  Server  `yaml:"server" json:"server"`
	MongoDB MongoDB `yaml:"mongo_db" json:"mongo_db"`
}

type Server struct {
	Port            string        `yaml:"port" json:"port"`
	ReadTimeout     time.Duration `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout" json:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" json:"shutdown_timeout"`
}

type MongoDB struct {
	Collections struct {
		Tasks string `yaml:"tasks" json:"tasks"`
	} `yaml:"collection" json:"collection"`
	Timeout               time.Duration `yaml:"timeout" json:"timeout"`
	DefaultContextTimeout time.Duration `yaml:"default_context_timeout" json:"default_context_timeout"`
	AppName               string        `yaml:"app_name" json:"app_name"`
}
