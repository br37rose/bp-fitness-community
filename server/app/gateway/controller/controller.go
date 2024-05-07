package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cache/mongodbcache"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/distributedmutex"
	pm "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/paymentprocessor/stripe"
	s3_storage "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/storage/s3"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/templatedemailer"
	gateway_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/gateway/datastore"
	organization_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/datastore"
	rp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	u_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/jwt"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/password"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

type GatewayController interface {
	MemberRegister(ctx context.Context, req *MemberRegisterRequestIDO) (*gateway_s.LoginResponseIDO, error)
	Login(ctx context.Context, email, password string) (*gateway_s.LoginResponseIDO, error)
	GetUserBySessionID(ctx context.Context, sessionID string) (*user_s.User, error)
	RefreshToken(ctx context.Context, value string) (*user_s.User, string, time.Time, string, time.Time, error)
	Verify(ctx context.Context, code string) error
	Logout(ctx context.Context) error
	ForgotPassword(ctx context.Context, email string) error
	PasswordReset(ctx context.Context, code string, password string) error
	Account(ctx context.Context) (*user_s.User, error)
	AccountUpdate(ctx context.Context, nu *user_s.User) error
	AccountChangePassword(ctx context.Context, req *AccountChangePasswordRequestIDO) error
	AccountListLatestInvoices(ctx context.Context, cursor int64, limit int64) (*user_s.StripeInvoiceListResult, error)
	Avatar(ctx context.Context, req *AccountOperationAvatarRequest) (*u_s.User, error)
	GenerateOTP(ctx context.Context) (*OTPGenerateResponseIDO, error)
	GenerateOTPAndQRCodePNGImage(ctx context.Context) ([]byte, error)
	VerifyOTP(ctx context.Context, req *VerificationTokenRequestIDO) (*VerificationTokenResponseIDO, error)
	ValidateOTP(ctx context.Context, req *ValidateTokenRequestIDO) (*ValidateTokenResponseIDO, error)
	DisableOTP(ctx context.Context) (*u_d.User, error)
}

type GatewayControllerImpl struct {
	Config             *config.Conf
	Logger             *slog.Logger
	UUID               uuid.Provider
	JWT                jwt.Provider
	Password           password.Provider
	Cache              mongodbcache.Cacher
	S3                 s3_storage.S3Storager
	DbClient           *mongo.Client
	DistributedMutex   distributedmutex.Adapter
	TemplatedEmailer   templatedemailer.TemplatedEmailer
	PaymentProcessor   pm.PaymentProcessor
	UserStorer         user_s.UserStorer
	RankPointStorer    rp_s.RankPointStorer
	OrganizationStorer organization_s.OrganizationStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	jwtp jwt.Provider,
	passwordp password.Provider,
	cache mongodbcache.Cacher,
	s3 s3_storage.S3Storager,
	client *mongo.Client,
	dmux distributedmutex.Adapter,
	te templatedemailer.TemplatedEmailer,
	paymentProcessor pm.PaymentProcessor,
	rp_storer rp_s.RankPointStorer,
	usr_storer user_s.UserStorer,
	org_storer organization_s.OrganizationStorer,
) GatewayController {
	s := &GatewayControllerImpl{
		Config:             appCfg,
		Logger:             loggerp,
		UUID:               uuidp,
		JWT:                jwtp,
		Password:           passwordp,
		S3:                 s3,
		Cache:              cache,
		TemplatedEmailer:   te,
		DbClient:           client,
		DistributedMutex:   dmux,
		PaymentProcessor:   paymentProcessor,
		RankPointStorer:    rp_storer,
		UserStorer:         usr_storer,
		OrganizationStorer: org_storer,
	}
	// s.Logger.Debug("gateway controller initialization started...")
	// s.Logger.Debug("gateway controller initialized")
	return s
}

func (impl *GatewayControllerImpl) GetUserBySessionID(ctx context.Context, sessionID string) (*user_s.User, error) {

	userStr, err := impl.Cache.Get(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if userStr == "" {
		impl.Logger.Warn("record not found")
		return nil, errors.New("record not found")
	}
	var user user_s.User
	err = json.Unmarshal([]byte(userStr), &user)
	if err != nil {
		impl.Logger.Error("unmarshalling failed", slog.Any("err", err))
		return nil, err
	}
	return &user, nil
}
