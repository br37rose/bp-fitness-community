package mailgun

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"log/slog"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/provider/uuid"
)

type Emailer interface {
	Send(ctx context.Context, sender, subject, recipient, htmlContent string) error
	GetSenderEmail() string
	GetBackendDomainName() string
	GetFrontendDomainName() string
	GetMaintenanceEmail() string
}

type mailgunEmailer struct {
	Mailgun          *mailgun.MailgunImpl
	UUID             uuid.Provider
	Logger           *slog.Logger
	senderEmail      string
	apiDomainName    string
	appDomainName    string
	maintenanceEmail string
}

func NewEmailer(cfg *c.Conf, logger *slog.Logger, uuidp uuid.Provider) Emailer {
	// Defensive code: Make sure we have access to the file before proceeding any further with the code.
	logger.Debug("mailgun emailer initializing...")
	mg := mailgun.NewMailgun(cfg.Emailer.Domain, cfg.Emailer.APIKey)
	logger.Debug("mailgun emailer was initialized.")

	mg.SetAPIBase(cfg.Emailer.APIBase) // Override to support our custom email requirements.

	return &mailgunEmailer{
		Mailgun:          mg,
		UUID:             uuidp,
		Logger:           logger,
		senderEmail:      cfg.Emailer.SenderEmail,
		apiDomainName:    cfg.AppServer.APIDomainName,
		appDomainName:    cfg.AppServer.AppDomainName,
		maintenanceEmail: cfg.Emailer.MaintenanceEmail,
	}
}

func (me *mailgunEmailer) Send(ctx context.Context, sender, subject, recipient, body string) error {
	me.Logger.Debug("sent email",
		slog.String("sender", sender),
		slog.String("subject", subject),
		slog.String("recipient", recipient),
		slog.String("body", body))

	message := me.Mailgun.NewMessage(sender, subject, "", recipient)
	message.SetHtml(body)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := me.Mailgun.Send(ctx, message)

	if err != nil {
		me.Logger.Error("emailer failed sending", slog.Any("err", err))
		return err
	}

	me.Logger.Debug("emailer sent with response", slog.Any("id", id), slog.Any("resp", resp))

	return nil
}

func (me *mailgunEmailer) GetSenderEmail() string {
	return me.senderEmail
}

func (me *mailgunEmailer) GetBackendDomainName() string {
	return me.apiDomainName
}

func (me *mailgunEmailer) GetFrontendDomainName() string {
	return me.appDomainName
}

func (me *mailgunEmailer) GetMaintenanceEmail() string {
	return me.maintenanceEmail
}
