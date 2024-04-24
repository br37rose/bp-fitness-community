package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"log/slog"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config"
)

type Cacher interface {
	Shutdown()
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, val []byte) error
	SetWithExpiry(ctx context.Context, key string, val []byte, expiry time.Duration) error
	Delete(ctx context.Context, key string) error
}

type cache struct {
	Client *redis.Client
	Logger *slog.Logger
}

func NewCache(cfg *c.Conf, logger *slog.Logger) Cacher {
	logger.Debug("cache initializing...")
	opt, err := redis.ParseURL(cfg.Cache.URI)
	if err != nil {
		logger.Error("cache failed parsing url error", slog.Any("err", err), slog.String("URI", cfg.Cache.URI))
		log.Fatal(err)
	}
	rdb := redis.NewClient(opt)

	// Confirm connection made with our application and redis.
	response, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logger.Error("cache failed pinging redis", slog.Any("err", err), slog.String("URI", cfg.Cache.URI))
		log.Fatal(err)
	}

	logger.Debug("cache initialized with successful redis ping response", slog.Any("resp", response))
	return &cache{
		Client: rdb,
		Logger: logger,
	}
}

func (s *cache) Shutdown() {
	s.Client.Close()
}

func (s *cache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := s.Client.Get(ctx, key).Result()
	if err != nil {
		s.Logger.Error("cache get failed", slog.Any("error", err))
		return nil, err
	}
	return []byte(val), nil
}

func (s *cache) Set(ctx context.Context, key string, val []byte) error {
	err := s.Client.Set(ctx, key, val, 0).Err()
	if err != nil {
		s.Logger.Error("cache set failed", slog.Any("error", err))
		return err
	}
	return nil
}

func (s *cache) SetWithExpiry(ctx context.Context, key string, val []byte, expiry time.Duration) error {
	err := s.Client.Set(ctx, key, val, expiry).Err()
	if err != nil {
		s.Logger.Error("cache set with expiry failed", slog.Any("error", err))
		return err
	}
	return nil
}

func (s *cache) Delete(ctx context.Context, key string) error {
	err := s.Client.Del(ctx, key).Err()
	if err != nil {
		s.Logger.Error("cache delete failed", slog.Any("error", err))
		return err
	}
	return nil
}
