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

func (impl *templatedEmailer) SendMemberWaitlistersAvailableOpeningEmail(waitlisterEmails []string, branchID string, trainerName string, workoutSessionID string, workoutProgramName string, startAt time.Time) error {
	fp := path.Join("templates", "member_waitlister_available_opening.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		impl.Logger.Error("parsing error", slog.Any("error", err))
		return err
	}

	var processed bytes.Buffer

	// Render the HTML template with our data.
	data := struct {
		DetailsLink        string
		WorkoutProgramName string
		TrainerName        string
		StartAt            string
	}{
		DetailsLink:        "https://" + impl.Emailer.GetFrontendDomainName() + "/session/" + workoutSessionID,
		WorkoutProgramName: workoutProgramName,
		TrainerName:        trainerName,
		StartAt:            timekit.ToAmericanDateTimeString(startAt),
	}
	if err := tmpl.Execute(&processed, data); err != nil {
		impl.Logger.Error("template execution error", slog.Any("error", err))
		return err
	}
	body := processed.String() // DEVELOPERS NOTE: Convert our long sequence of data into a string.

	for _, waitlisterEmail := range waitlisterEmails {
		if err := impl.Emailer.Send(context.Background(), impl.Emailer.GetSenderEmail(), "Class schedule opening!", waitlisterEmail, body); err != nil {
			impl.Logger.Error("sending error", slog.Any("error", err))
			// return err // To simplify let us not error.
		}
	}

	return nil
}
