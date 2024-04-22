package redis

import (
	"context"
	"crypto/tls"
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
			TLSConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
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
			slog.Bool("is_cluster", appCfg.Redis.IsClusterConfiguration),
			slog.Any("addrs", clusterOptions.Addrs),
			slog.String("username", clusterOptions.Username),
			slog.String("password", clusterOptions.Password),
		)

		// Ping Redis to check if the connection is working
		_, err := client.Ping(context.Background()).Result()
		if err != nil {
			logger.Error("failed connecting to redis cluster with error: %v", err)
			return nil
		}

		logger.Debug("redis cluster initialized successfully")
		return client
	}

	logger.Debug("redis initializing...")

	addr := appCfg.Redis.Addresses[0]

	// Configure the manditory options:
	opts := &redis.Options{
		Addr: addr,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
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
		slog.Bool("is_cluster", appCfg.Redis.IsClusterConfiguration),
		slog.String("addr", addr),
		slog.String("username", opts.Username),
		slog.String("password", opts.Password),
	)

	// Ping Redis to check if the connection is working
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		logger.Error("failed connecting to redis with error: %v", err)
		return nil
	}

	logger.Debug("redis initialized successfully")
	return client
}
