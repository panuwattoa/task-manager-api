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
	Server  Server
	MongoDB MongoDB
}

type Server struct {
	AppName         string
	Env             string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type MongoDB struct {
	Collections struct {
		Tasks string
	}
	Timeout               time.Duration
	DefaultContextTimeout time.Duration
	AppName               string
}
