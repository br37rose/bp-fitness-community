package mongodbcache

import (
	"context"
	"time"

	"github.com/faabiosr/cachego"
	"github.com/faabiosr/cachego/mongo"
	mongo_client "go.mongodb.org/mongo-driver/mongo"
	"log/slog"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config"
)

type Cacher interface {
	Shutdown()
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, val string) error
	SetWithExpiry(ctx context.Context, key string, val string, expiry time.Duration) error
	Delete(ctx context.Context, key string) error
}

type cache struct {
	Client cachego.Cache
	Logger *slog.Logger
}

func NewCache(cfg *c.Conf, logger *slog.Logger, dbClient *mongo_client.Client) Cacher {
	logger.Debug("cache initializing...")

	cc := dbClient.Database(cfg.DB.Name).Collection("cache")

	c := mongo.New(cc)

	logger.Debug("cache initialized with mongodb as backend")
	return &cache{
		Client: c,
		Logger: logger,
	}
}

func (s *cache) Shutdown() {
	// Do nothing...
}

func (s *cache) Get(ctx context.Context, key string) (string, error) {
	val, err := s.Client.Fetch(key)
	if err != nil {
		s.Logger.Error("cache get failed", slog.Any("error", err))
		return "", err
	}
	return string(val), nil
}

func (s *cache) Set(ctx context.Context, key string, val string) error {
	err := s.Client.Save(key, val, 0)
	if err != nil {
		s.Logger.Error("cache set failed", slog.Any("error", err))
		return err
	}
	return nil
}

func (s *cache) SetWithExpiry(ctx context.Context, key string, val string, expiry time.Duration) error {
	err := s.Client.Save(key, val, expiry)
	if err != nil {
		s.Logger.Error("cache set with expiry failed", slog.Any("error", err))
		return err
	}
	return nil
}

func (s *cache) Delete(ctx context.Context, key string) error {
	err := s.Client.Delete(key)
	if err != nil {
		s.Logger.Error("cache delete failed", slog.Any("error", err))
		return err
	}
	return nil
}
