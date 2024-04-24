package templatedemailer

import (
	"bytes"
	"context"
	"path"
	"text/template"

	"log/slog"
)

func (impl *templatedEmailer) SendForgotPasswordEmail(email, verificationCode, firstName string) error {
	// FOR TESTING PURPOSES ONLY.
	fp := path.Join("templates", "forgot_password.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		impl.Logger.Error("parsing error", slog.Any("error", err))
		return err
	}

	var processed bytes.Buffer

	// Render the HTML template with our data.
	data := struct {
		Email            string
		VerificationLink string
		FirstName        string
	}{
		Email:            email,
		VerificationLink: "https://" + impl.Emailer.GetFrontendDomainName() + "/password-reset?q=" + verificationCode,
		FirstName:        firstName,
	}
	if err := tmpl.Execute(&processed, data); err != nil {
		impl.Logger.Error("template execution error", slog.Any("error", err))
		return err
	}
	body := processed.String() // DEVELOPERS NOTE: Convert our long sequence of data into a string.

	if err := impl.Emailer.Send(context.Background(), impl.Emailer.GetSenderEmail(), "Forgot Password", email, body); err != nil {
		impl.Logger.Error("sending error", slog.Any("error", err))
		return err
	}
	return nil
}
