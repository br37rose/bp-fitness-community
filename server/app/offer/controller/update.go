package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	e_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/datastore"
	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/offer/datastore"
	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	vcon_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocontent/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (c *OfferControllerImpl) UpdateByID(ctx context.Context, ns *domain.Offer) (*domain.Offer, error) {
	// Extract from our session the following data.
	urole := ctx.Value(constants.SessionUserRole).(int8)
	// uid := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// uname := ctx.Value(constants.SessionUserName).(string)
	oid := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	oname := ctx.Value(constants.SessionUserOrganizationName).(string)

	switch urole { // Security.
	case u_d.UserRoleRoot:
		return nil, httperror.NewForForbiddenWithSingleField("message", "you did not saasify offer")
	case u_d.UserRoleTrainer:
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	case u_d.UserRoleMember:
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	}

	// Fetch the original organization.
	os, err := c.OfferStorer.GetByID(ctx, ns.ID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if os == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "workout program type does not exist")
	}

	// Modify our original organization.
	os.OrganizationID = oid
	os.OrganizationName = oname
	os.ModifiedAt = time.Now()
	os.Status = ns.Status
	os.BusinessFunction = ns.BusinessFunction
	os.MembershipRank = ns.MembershipRank

	// Save to the database the modified organization.
	if err := c.OfferStorer.UpdateByID(ctx, os); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return nil, err
	}

	if err := c.updateRelatedExercises(ctx, os); err != nil {
		c.Logger.Error("update related exercises error", slog.Any("error", err))
		return nil, err
	}

	if err := c.updateRelatedVideoContents(ctx, os); err != nil {
		c.Logger.Error("update related video content error", slog.Any("error", err))
		return nil, err
	}

	return os, nil
}

func (c *OfferControllerImpl) updateRelatedExercises(ctx context.Context, offer *domain.Offer) error {
	f := &e_s.ExerciseListFilter{
		Cursor:    primitive.NilObjectID,
		PageSize:  1_000_000_000,
		SortField: "_id",
		SortOrder: 1,
		OfferID:   offer.ID,
	}
	res, err := c.ExerciseStorer.ListByFilter(ctx, f)
	if err != nil {
		return err
	}
	for _, vc := range res.Results {
		vc.OfferID = offer.ID
		vc.OfferName = offer.Name
		vc.OfferMembershipRank = offer.MembershipRank
		if err := c.ExerciseStorer.UpdateByID(ctx, vc); err != nil {
			c.Logger.Error("database update by id error",
				slog.Any("error", err))
			return err
		}
		c.Logger.Error("updated exercise",
			slog.Any("ID", vc.ID))
	}
	return nil
}

func (c *OfferControllerImpl) updateRelatedVideoContents(ctx context.Context, offer *domain.Offer) error {
	f := &vcon_s.VideoContentListFilter{
		Cursor:    primitive.NilObjectID,
		PageSize:  1_000_000_000,
		SortField: "_id",
		SortOrder: 1,
		OfferID:   offer.ID,
	}
	res, err := c.VideoContentStorer.ListByFilter(ctx, f)
	if err != nil {
		return err
	}
	for _, vc := range res.Results {
		vc.OfferID = offer.ID
		vc.OfferName = offer.Name
		vc.OfferMembershipRank = offer.MembershipRank
		if err := c.VideoContentStorer.UpdateByID(ctx, vc); err != nil {
			c.Logger.Error("database update by id error",
				slog.Any("error", err))
			return err
		}
		c.Logger.Error("updated video content",
			slog.Any("ID", vc.ID))
	}
	return nil
}
