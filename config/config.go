package config

import (
	"os"
	"time"
)

var Conf = new(Config)

var (
	ServerType    = GetEnv("SERVER_TYPE", "Development")
	MongoHost     = GetEnv("MONGO_HOST", "mongodb://127.0.0.1:27017")
	MongoDBName   = GetEnv("MONGO_DBNAME", "taskManager")
	MongoUser     = GetEnv("MONGO_USERNAME", "managerapp")
	MongoPassword = GetEnv("MONGO_PASSWORD", "1111")
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type Config struct {
	Server     Server
	MongoDB    MongoDB
	Pagination struct {
		MaxLimit           int
		MaxGetProfileLimit int
	}
}

type Server struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type MongoDB struct {
	Collections struct {
		Tasks    string
		Profiles string
		Comments string
	}
	Timeout               time.Duration
	DefaultContextTimeout time.Duration
	AppName               string
}
