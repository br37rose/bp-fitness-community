package mongodb

import (
	"context"
	"log"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
)

// NewProvider function will initiate a connection with MongoDB and supply an instance of the connection to our app. This code is useful for dependency injection.
func NewProvider(appCfg *c.Conf, logger *slog.Logger) *mongo.Client {
	logger.Debug("mongodb storage initializing...")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(appCfg.DB.URI))
	if err != nil {
		log.Fatal(err)
	}

	// The MongoDB client provides a Ping() method to tell you if a MongoDB database has been found and connected.
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	logger.Debug("mongodb storage initialized successfully")
	return client
}
