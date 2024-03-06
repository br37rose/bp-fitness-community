package templatedemailer

import (
	"bytes"
	"context"
	"log/slog"
	"path"
	"text/template"
	"time"

	"github.com/bartmika/timekit"
)

func (impl *templatedEmailer) SendMemberCancelledSubscriptionEmailToMember(memberEmail string, memberFirstName string, subscriptionName string, cancelationDate time.Time) error {
	////
	//// Send notification email to member.
	////

	fp := path.Join("templates", "member_subscription_cancelled.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		impl.Logger.Error("parsing error", slog.Any("error", err))
		return err
	}

	var processed bytes.Buffer

	// Render the HTML template with our data.
	data := struct {
		DetailsLink      string
		MemberFirstName  string
		SubscriptionName string
		CancelationDate  string
	}{
		DetailsLink:      "https://" + impl.Emailer.GetFrontendDomainName() + "/account/subscription",
		MemberFirstName:  memberFirstName,
		SubscriptionName: subscriptionName,
		CancelationDate:  timekit.ToAmericanDateTimeString(cancelationDate),
	}
	if err := tmpl.Execute(&processed, data); err != nil {
		impl.Logger.Error("template execution error", slog.Any("error", err))
		return err
	}
	body := processed.String() // DEVELOPERS NOTE: Convert our long sequence of data into a string.

	if err := impl.Emailer.Send(context.Background(), impl.Emailer.GetSenderEmail(), "You canceled a subscription", memberEmail, body); err != nil {
		impl.Logger.Error("sending error", slog.Any("error", err))
		// return err // To simplify let us not error.
	}
	return nil
}
