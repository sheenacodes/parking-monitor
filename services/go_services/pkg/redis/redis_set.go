package redis

import (
	"context"
	"fmt"

	"go_services/pkg/logger"
)

func (r *RedisClient) IsSetNotEmpty(redisSetName string) (bool, error) {
	// Use the SCARD command to get the number of members in the set
	ctx := context.Background()
	card, err := r.Client.SCard(ctx, redisSetName).Result()
	if err != nil {
		logger.Log.Error().Err(err).Msg("error checking set size")
		return false, fmt.Errorf("error checking set size: %v", err)
	}

	logger.Log.Debug().Msgf("%d items available in Redis Set", card)
	// Return true if the number of members is greater than 0
	return card > 0, nil
}

func (r *RedisClient) AddItemToSet(item string, redisSetName string) error {
	ctx := context.Background()
	_, err := r.Client.SAdd(ctx, redisSetName, item).Result()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to make entry to Redis Set")
		return err
	}
	logger.Log.Debug().Msgf(" %s added to Redis Set %s", item, redisSetName)
	return nil
}

func (r *RedisClient) RemoveItemFromSet(item string, redisSetName string) error {
	ctx := context.Background()
	_, err := r.Client.SRem(ctx, redisSetName, item).Result()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to remove entry from Redis Set")
		return err
	}
	logger.Log.Debug().Msgf("%s removed from Redis Set %s", item, redisSetName)
	return nil
}

func (r *RedisClient) GetRandomItemFromSet(redisSetName string) (string, error) {
	// Use SRANDMEMBER to get a random member from the set
	ctx := context.Background()
	item, err := r.Client.SRandMember(ctx, redisSetName).Result()
	if err != nil {
		return "", fmt.Errorf("could not get random member from set: %w", err)
	}
	return item, nil
}
