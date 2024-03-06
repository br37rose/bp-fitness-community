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

func (impl *templatedEmailer) SendMemberCancelledBookingEmailToTrainer(trainerEmail string, trainerFirstName string, workoutProgramName string, branchID string, memberID string, memberName string, startAt time.Time) error {
	////
	//// Send notification email to trainer.
	////

	fp := path.Join("templates", "trainer_member_booking_cancelled.html")
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
		MemberName         string
		StartAt            string
	}{
		DetailsLink:        "https://" + impl.Emailer.GetFrontendDomainName() + "/admin/branch/" + branchID + "/member/" + memberID,
		TrainerFirstName:   trainerFirstName,
		WorkoutProgramName: workoutProgramName,
		MemberName:         memberName,
		StartAt:            timekit.ToAmericanDateTimeString(startAt),
	}
	if err := tmpl.Execute(&processed, data); err != nil {
		impl.Logger.Error("template execution error", slog.Any("error", err))
		return err
	}
	body := processed.String() // DEVELOPERS NOTE: Convert our long sequence of data into a string.

	if err := impl.Emailer.Send(context.Background(), impl.Emailer.GetSenderEmail(), "You has a cancelion", trainerEmail, body); err != nil {
		impl.Logger.Error("sending error", slog.Any("error", err))
		// return err // To simplify let us not error.
	}
	return nil
}
