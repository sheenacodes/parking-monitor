package config

import (
	"os"
	"strconv"
)

type Config struct {
	RabbitMQURL    string
	EntryQueueName string
	ExitQueueName  string
	LogLevel       string
	RedisAddress   string
	RedisPassword  string
	RedisDB        int
}

func LoadConfig() *Config {
	return &Config{
		RabbitMQURL:    getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		EntryQueueName: getEnv("RABBITMQ_ENTRY_QUEUE_NAME", ""),
		ExitQueueName:  getEnv("RABBITMQ_EXIT_QUEUE_NAME", ""),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		RedisAddress:   getEnv("REDIS_ADDR", "redis"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		RedisDB:        getEnvAsInt("REDIS_DB", 0),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
