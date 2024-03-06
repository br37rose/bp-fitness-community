package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	o_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/offer/datastore"
)

func (impl *OfferControllerImpl) createDefaults(ctx context.Context) error {
	impl.Logger.Debug("offer createDefaults started...")
	// DEVELOPERS NOTE: ALL OF THESE ARE PURPOPSEFULLY HARD-CODE AND ARE TO BE DELETED UPON RELEASE.

	orgID, _ := primitive.ObjectIDFromHex(impl.Config.AppServer.InitialOrgID)

	// ---- Monthly Subscription ----
	o1ID, _ := primitive.ObjectIDFromHex("64da7f02443e34e206d03fa6") // Hard-coded
	o1, _ := impl.GetByID(ctx, o1ID)
	if o1 == nil || true { // Temporary always true, delete when ready to launch.
		impl.Logger.Debug("creating `Monthly`...")
		o1 = &o_ds.Offer{
			OrganizationID:       orgID,
			OrganizationName:     "BP8 Fitness",
			ID:                   o1ID,
			Name:                 "Monthly Membership",
			Description:          "Get unlimited access to all our content except VIP content.",
			Price:                5.99,
			PriceCurrency:        "CAD",
			PayFrequency:         o_ds.PayFrequencyMonthly,
			Status:               o_ds.StatusActive,
			Type:                 o_ds.OfferTypeService,
			BusinessFunction:     o_ds.BusinessFunctionProvideMembershipAccessToContentAccess,
			CreatedAt:            time.Now(),
			ModifiedAt:           time.Now(),
			StripeProductID:      "prod_OfpvwCae7C6OIm",            // Hard-coded
			StripePriceID:        "price_1NsUFHH2wIbBWH08wa7XqAa8", // Hard-coded
			IsSubscription:       true,
			MembershipRank:       4, //4=regular
			PaymentProcessorName: "Stripe, Inc.",
			// IncludesOfferIDs: []primitive.ObjectID{o1ID},
		}
		if err := impl.OfferStorer.Upsert(ctx, o1); err != nil {
			return err
		}
	}

	// ---- VIP Legendary (Monthly) Membership ----
	o2ID, _ := primitive.ObjectIDFromHex("64da7ef78d1c8fef4b76c9ad") // Hard-coded
	o2, _ := impl.GetByID(ctx, o2ID)
	if o2 == nil || true { // Temporary always true, delete when ready to launch.
		impl.Logger.Debug("creating `VIP Legendary Membership`...")
		o2 = &o_ds.Offer{
			OrganizationID:       orgID,
			OrganizationName:     "BP8 Fitness",
			ID:                   o2ID,
			Name:                 "VIP Legendary Membership",
			Description:          "Join the elite fitness enthusiasts with our VIP Legendary Membership for full access to all the content.",
			Price:                8.99,
			PriceCurrency:        "CAD",
			PayFrequency:         o_ds.PayFrequencyMonthly,
			Status:               o_ds.StatusActive,
			Type:                 o_ds.OfferTypeService,
			BusinessFunction:     o_ds.BusinessFunctionProvideMembershipAccessToContentAccess,
			CreatedAt:            time.Now(),
			ModifiedAt:           time.Now(),
			StripeProductID:      "prod_OfrehRlxukLDpb",
			StripePriceID:        "price_1NsVuwH2wIbBWH082AMPO1Qj",
			IsSubscription:       true,
			MembershipRank:       5, // 5=elite
			PaymentProcessorName: "Stripe, Inc.",
		}
		if err := impl.OfferStorer.Upsert(ctx, o2); err != nil {
			return err
		}
	}

	// uID := primitive.NewObjectID()
	// impl.Logger.Warn("id.", slog.Any("id", uID))

	impl.Logger.Debug("offer createDefaults finished")
	return nil
}
