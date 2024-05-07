package redis

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
)

func NewProvider(appCfg *c.Conf, logger *slog.Logger) redis.UniversalClient {
	// Extract the Redis URL from the configuration
	redisURL := appCfg.Redis.URL
	if redisURL == "" {
		log.Fatal("Redis URL is empty")
	}

	// Parse the Redis URL
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed connecting to Redis with error: %v", err)
	}

	// Extra options
	opt.DialTimeout = 1 * time.Minute

	// Create the Redis client
	client := redis.NewClient(opt)

	// Ping Redis to check if the connection is working
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed connecting to Redis with error: %v", err)
	}

	logger.Debug("Redis client initialized successfully")
	return client
}
