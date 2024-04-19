package redis

import (
	"context"
	"log"
	"log/slog"

	"github.com/redis/go-redis/v9"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
)

func NewProvider(appCfg *c.Conf, logger *slog.Logger) redis.UniversalClient {
	// Extract the optional fields from our configuration.
	var username = ""
	if appCfg.Redis.Username != "" {
		username = appCfg.Redis.Username
	}
	var password = ""
	if appCfg.Redis.Password != "" {
		password = appCfg.Redis.Password
	}

	if appCfg.Redis.IsClusterConfiguration {
		logger.Debug("redis cluster initializing...")

		// Configure the manditory options:
		clusterOptions := &redis.ClusterOptions{
			Addrs: appCfg.Redis.Addresses,
		}

		// Configure the optional options:
		if username != "" {
			clusterOptions.Username = username
		}
		if password != "" {
			clusterOptions.Password = password
		}

		// Create our redis client.
		client := redis.NewClusterClient(clusterOptions)

		logger.Debug("redis cluster checking connection...",
			slog.Any("addrs", clusterOptions.Addrs),
			slog.String("username", clusterOptions.Username),
		)

		// Ping Redis to check if the connection is working
		_, err := client.Ping(context.Background()).Result()
		if err != nil {
			log.Fatalf("failed connecting to redis with error: %v", err)
		}

		logger.Debug("redis cluster initialized successfully")
		return client
	}

	logger.Debug("redis initializing...")

	// Configure the manditory options:
	opts := &redis.Options{
		Addr: appCfg.Redis.Addresses[0],
	}

	// Configure the optional options:
	if username != "" {
		opts.Username = username
	}
	if password != "" {
		opts.Password = password
	}

	client := redis.NewClient(opts)

	logger.Debug("redis checking connection...",
		slog.String("addr", opts.Addr),
		slog.String("username", opts.Username),
	)

	// Ping Redis to check if the connection is working
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed connecting to redis with error: %v", err)
	}

	logger.Debug("redis initialized successfully")
	return client
}
