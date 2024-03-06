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

func (impl *templatedEmailer) SendTrainerSessionFullEmail(trainerEmail string, trainerFirstName string, branchID string, workoutSessionID string, workoutProgramName string, startAt time.Time) error {
	fp := path.Join("templates", "trainer_session_full.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		impl.Logger.Error("parsing error", slog.Any("error", err))
		return err
	}

	var processed bytes.Buffer

	// Render the HTML template with our data.
	data := struct {
		DetailsLink        string
		TrainerFirstName   string
		WorkoutProgramName string
		StartAt            string
	}{
		DetailsLink:        "https://" + impl.Emailer.GetFrontendDomainName() + "/admin/branch/" + branchID + "/session/" + workoutSessionID + "/bookings",
		TrainerFirstName:   trainerFirstName,
		WorkoutProgramName: workoutProgramName,
		StartAt:            timekit.ToAmericanDateTimeString(startAt),
	}
	if err := tmpl.Execute(&processed, data); err != nil {
		impl.Logger.Error("template execution error", slog.Any("error", err))
		return err
	}
	body := processed.String() // DEVELOPERS NOTE: Convert our long sequence of data into a string.

	if err := impl.Emailer.Send(context.Background(), impl.Emailer.GetSenderEmail(), "Class schedule fully booked", trainerEmail, body); err != nil {
		impl.Logger.Error("sending error", slog.Any("error", err))
		// return err // To simplify let us not error.
	}
	return nil
}
