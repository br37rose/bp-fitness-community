package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	s_d "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/exercise/datastore"
)

func (c *ExerciseControllerImpl) Create(ctx context.Context, m *s_d.Exercise) (*s_d.Exercise, error) {
	// // Modify the exercise based on role.
	// userRole, ok := ctx.Value(constants.SessionUserRole).(int8)
	// if ok {
	// 	switch userRole {
	// 	case u_d.RetailerRole:
	// 		// Override status.
	// 		m.Status = s_d.ExerciseStatusPending
	//
	// 		// Auto-assign the user-if
	// 		m.UserID = ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// 		m.UserFirstName = ctx.Value(constants.SessionUserFirstName).(string)
	// 		m.UserLastName = ctx.Value(constants.SessionUserLastName).(string)
	// 		m.UserCompanyName = ctx.Value(constants.SessionUserCompanyName).(string)
	// 		m.ServiceType = s_d.PreScreeningServiceType
	// 	case u_d.StaffRole:
	// 		m.Status = s_d.ExerciseStatusActive
	// 	default:
	// 		m.Status = s_d.ExerciseStatusError
	// 	}
	// }

	// Add defaults.
	m.ID = primitive.NewObjectID()
	m.CreatedAt = time.Now()
	m.ModifiedAt = time.Now()

	// Save to our database.
	err := c.ExerciseStorer.Create(ctx, m)
	if err != nil {
		c.Logger.Error("database create error", slog.Any("error", err))
		return nil, err
	}

	return m, nil
}
