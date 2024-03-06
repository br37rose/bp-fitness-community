package controller

import (
	"context"
	"time"

	"log/slog"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (impl *GatewayControllerImpl) PasswordReset(ctx context.Context, code string, password string) error {
	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := impl.UserStorer.GetByVerificationCode(ctx, code)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return err
	}
	if u == nil {
		impl.Logger.Warn("user does not exist validation error")
		return httperror.NewForBadRequestWithSingleField("code", "does not exist")
	}

	//TODO: Handle expiry dates.

	passwordHash, err := impl.Password.GenerateHashFromPassword(password)
	if err != nil {
		impl.Logger.Error("hashing error", slog.Any("error", err))
		return err
	}

	u.PasswordHash = passwordHash
	u.PasswordHashAlgorithm = impl.Password.AlgorithmName()
	u.EmailVerificationCode = "" // Remove email active code so it cannot be used agian.
	u.EmailVerificationExpiry = time.Now()
	u.ModifiedAt = time.Now()

	if err := impl.UserStorer.UpdateByID(ctx, u); err != nil {
		impl.Logger.Error("update error", slog.Any("err", err))
		return err
	}

	return nil
}
