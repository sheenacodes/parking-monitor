package config

import (
	"os"
	"strconv"
)

type Config struct {
	RabbitMQURL   string
	QueueName     string
	LogLevel      string
	GeneratorMode string
	RedisAddress  string
	RedisPassword string
	RedisDB       int
}

func LoadConfig() *Config {
	return &Config{
		RabbitMQURL:   getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		QueueName:     getEnv("ENTRY_QUEUE_NAME", "entry_events_queue"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		GeneratorMode: getEnv("GENERATOR_MODE", "entry"),
		RedisAddress:  getEnv("REDIS_ADDR", "redis"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),
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
