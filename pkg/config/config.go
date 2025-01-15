package config

import "os"

type Config struct {
	MongoDB string
	Kafka   []string
}

func LoadConfig() Config {
	return Config{
		MongoDB: os.Getenv("MONGODB_DATABASE"),
		Kafka:   []string{"localhost:9092"},
	}
}
