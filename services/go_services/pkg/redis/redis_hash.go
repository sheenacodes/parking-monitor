package redis

import (
	"context"
	"fmt"
	"time"

	"go_services/pkg/logger"

	"github.com/redis/go-redis/v9"
)

func (r *RedisClient) AddFieldToHash(hashKey string, fieldName string, fieldValue time.Time) error {
	ctx := context.Background()
	err := r.Client.HSet(ctx, hashKey, fieldName, fieldValue).Err()
	if err != nil {
		logger.Log.Error().Err(err).Msgf("Error Setting Hash Field %s for key %s", fieldName, hashKey)
		return err
	}
	logger.Log.Debug().Msgf(" %s:%s added to Redis Hash for key %s", fieldName, fieldValue, hashKey)
	return nil
}

func (r *RedisClient) GetFieldAsTime(hashKey string, fieldName string, layout string) (time.Time, error) {

	ctx := context.Background()
	value, err := r.Client.HGet(ctx, hashKey, fieldName).Result()
	if err != nil {
		if err == redis.Nil {
			// The field does not exist
			return time.Time{}, fmt.Errorf("field %s does not exist in hash %s", fieldName, hashKey)
		}
		return time.Time{}, fmt.Errorf("failed to get field value: %v", err)
	}

	// Parse the value as a time.Time
	parsedTime, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time: %v", err)
	}

	return parsedTime, nil
}
