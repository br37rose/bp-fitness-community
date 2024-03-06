package templatedemailer

import (
	"log/slog"
	"time"

	mg "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/emailer/mailgun"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

// TemplatedEmailer Is adapter for responsive HTML email templates sender.
type TemplatedEmailer interface {
	GetBackendDomainName() string
	GetFrontendDomainName() string
	SendMemberVerificationEmail(email, verificationCode, firstName, orgName string) error
	SendForgotPasswordEmail(email, verificationCode, firstName string) error
	SendMemberBookedSessionEmailToMember(memberEmail string, memberFirstName string, bookingID string, workoutProgramName string, trainerName string, startAt time.Time) error
	SendMemberCancelledBookingEmailToMember(memberEmail string, memberFirstName string, bookingID string, workoutProgramName string, trainerName string, startAt time.Time) error
	SendMemberWaitlistersAvailableOpeningEmail(waitlisterEmails []string, branchID string, trainerName string, workoutSessionID string, workoutProgramName string, startAt time.Time) error
	SendMemberCancelledBookingEmailToTrainer(trainerEmail string, trainerFirstName string, workoutProgramName string, branchID string, memberID string, memberName string, startAt time.Time) error
	SendMemberBookedSessionEmailToTrainer(trainerEmail string, trainerFirstName string, workoutProgramName string, branchID string, memberID string, memberName string, startAt time.Time) error
	SendMemberCancelledSubscriptionEmailToMember(memberEmail string, memberFirstName string, subscriptionName string, cancelationDate time.Time) error
	SendMemberSubscriptionStartedEmailToMember(memberEmail string, memberFirstName string, subscriptionName string, startedAt time.Time) error
	SendTrainerSessionFullEmail(trainerEmail string, trainerFirstName string, branchID string, workoutSessionID string, workoutProgramName string, startAt time.Time) error
	SendTrainerCancelledSessionEmailToMember(memberEmail string, memberFirstName string, trainerName string, workoutProgramName string, workoutSessionID string, startAt time.Time) error
}

type templatedEmailer struct {
	UUID    uuid.Provider
	Logger  *slog.Logger
	Emailer mg.Emailer
}

func NewTemplatedEmailer(cfg *c.Conf, logger *slog.Logger, uuidp uuid.Provider, emailer mg.Emailer) TemplatedEmailer {
	// Defensive code: Make sure we have access to the file before proceeding any further with the code.
	logger.Debug("templated emailer initializing...")
	logger.Debug("templated emailer initialized")

	return &templatedEmailer{
		UUID:    uuidp,
		Logger:  logger,
		Emailer: emailer,
	}
}

func (impl *templatedEmailer) GetBackendDomainName() string {
	return impl.Emailer.GetBackendDomainName()
}

func (impl *templatedEmailer) GetFrontendDomainName() string {
	return impl.Emailer.GetFrontendDomainName()
}
