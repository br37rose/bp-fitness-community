package templatedemailer

import (
	"time"

	mg "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/adapter/emailer/mailgun"
	"log/slog"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/provider/uuid"
)

// TemplatedEmailer Is adapter for responsive HTML email templates sender.
type TemplatedEmailer interface {
	SendMemberVerificationEmail(email, verificationCode, firstName, orgName string) error
	SendForgotPasswordEmail(email, verificationCode, firstName string) error
	SendMemberBookedSessionEmailToMember(memberEmail string, memberFirstName string, bookingID string, workoutProgramName string, trainerName string, startAt time.Time) error
	SendMemberCancelledBookingEmailToMember(memberEmail string, memberFirstName string, bookingID string, workoutProgramName string, trainerName string, startAt time.Time) error
	SendMemberWaitlistersAvailableOpeningEmail(waitlisterEmails []string, branchID string, trainerName string, workoutSessionID string, workoutProgramName string, startAt time.Time) error
	SendMemberCancelledBookingEmailToTrainer(trainerEmail string, trainerFirstName string, workoutProgramName string, branchID string, memberID string, memberName string, startAt time.Time) error
	SendMemberBookedSessionEmailToTrainer(trainerEmail string, trainerFirstName string, workoutProgramName string, branchID string, memberID string, memberName string, startAt time.Time) error
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
