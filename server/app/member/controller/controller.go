package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	pm "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/paymentprocessor/stripe"
	s3_storage "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/storage/s3"
	org_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/datastore"
	rp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
	u_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/password"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

// MemberController Interface for member business logic controller.
type MemberController interface {
	Create(ctx context.Context, requestData *MemberCreateRequestIDO) (*user_s.User, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*user_s.User, error)
	UpdateByID(ctx context.Context, requestData *MemberUpdateRequestIDO) (*user_s.User, error)
	ListByFilter(ctx context.Context, f *user_s.UserListFilter) (*user_s.UserListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *user_s.UserListFilter) ([]*user_s.UserAsSelectOption, error)
	ArchiveByID(ctx context.Context, id primitive.ObjectID) (*user_s.User, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CreateComment(ctx context.Context, memberID primitive.ObjectID, content string) (*user_s.User, error)
	Avatar(ctx context.Context, req *MemberOperationAvatarRequest) (*u_s.User, error)
}

type MemberControllerImpl struct {
	Config             *config.Conf
	Logger             *slog.Logger
	UUID               uuid.Provider
	S3                 s3_storage.S3Storager
	Password           password.Provider
	PaymentProcessor   pm.PaymentProcessor
	DbClient           *mongo.Client
	OrganizationStorer org_s.OrganizationStorer
	RankPointStorer    rp_s.RankPointStorer
	UserStorer         user_s.UserStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	s3 s3_storage.S3Storager,
	passwordp password.Provider,
	paymentProcessor pm.PaymentProcessor,
	client *mongo.Client,
	orgS org_s.OrganizationStorer,
	rp_storer rp_s.RankPointStorer,
	usrS user_s.UserStorer,
) MemberController {
	s := &MemberControllerImpl{
		Config:             appCfg,
		Logger:             loggerp,
		UUID:               uuidp,
		S3:                 s3,
		Password:           passwordp,
		PaymentProcessor:   paymentProcessor,
		DbClient:           client,
		OrganizationStorer: orgS,
		RankPointStorer:    rp_storer,
		UserStorer:         usrS,
	}
	s.Logger.Debug("member controller initialization started...")
	s.Logger.Debug("member controller initialized")
	return s
}
