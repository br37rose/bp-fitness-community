package templatedemailer

import (
	"bytes"
	"context"
	"log/slog"
	"path"
	"text/template"
)

func (impl *templatedEmailer) SendMemberVerificationEmail(email, verificationCode, firstName, orgName string) error {
	// FOR TESTING PURPOSES ONLY.
	fp := path.Join("templates", "gateway_email_verification.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		impl.Logger.Error("parsing error", slog.Any("error", err))
		return err
	}

	var processed bytes.Buffer

	// Render the HTML template with our data.
	data := struct {
		OrganizationName string
		Email            string
		VerificationLink string
		FirstName        string
	}{
		OrganizationName: orgName,
		Email:            email,
		VerificationLink: "https://" + impl.Emailer.GetFrontendDomainName() + "/verify?q=" + verificationCode,
		FirstName:        firstName,
	}
	if err := tmpl.Execute(&processed, data); err != nil {
		impl.Logger.Error("template execution error", slog.Any("error", err))
		return err
	}
	body := processed.String() // DEVELOPERS NOTE: Convert our long sequence of data into a string.

	if err := impl.Emailer.Send(context.Background(), impl.Emailer.GetSenderEmail(), "Activate your Account", email, body); err != nil {
		impl.Logger.Error("sending error", slog.Any("error", err))
		// return err // To simplify let us not error.
	}
	return nil
}
