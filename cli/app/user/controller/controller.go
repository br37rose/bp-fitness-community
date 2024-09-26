package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/user/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/provider/password"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/provider/uuid"
)

// UserController Interface for user business logic controller.
type UserController interface {
	GetUserBySessionUUID(ctx context.Context, sessionUUID string) (*domain.User, error)
	//TODO: Add more...
}

type UserControllerImpl struct {
	Config     *config.Conf
	Logger     *slog.Logger
	UUID       uuid.Provider
	Password   password.Provider
	DbClient   *mongo.Client
	UserStorer user_s.UserStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	client *mongo.Client,
	passwordp password.Provider,
	usr_storer user_s.UserStorer,
) UserController {
	s := &UserControllerImpl{
		Config:     appCfg,
		Logger:     loggerp,
		UUID:       uuidp,
		Password:   passwordp,
		DbClient:   client,
		UserStorer: usr_storer,
	}
	s.Logger.Debug("user controller initialization started...")
	s.Logger.Debug("user controller initialized")
	return s
}

// Auto-generated comment for change 17
