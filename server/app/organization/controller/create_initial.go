package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	o_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/datastore"
	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
)

// createInitialRootAdmin function creates the initial root administrator if not previously created.
func (c *OrganizationControllerImpl) createInitialRootAdmin(ctx context.Context) error {
	uid, err := primitive.ObjectIDFromHex(c.Config.AppServer.InitialRootAdminID)
	if err != nil {
		c.Logger.Error("primitive object id create from hex failed error", slog.Any("error", err))
		return err
	}

	org, err := c.UserStorer.GetByID(ctx, uid)
	if err != nil {
		c.Logger.Error("database check if exists error", slog.Any("error", err), slog.Any("InitialRootAdminID", uid))
		return err
	}
	if org == nil {
		c.Logger.Debug("No root user detected, proceeding to create now...")
		passwordHash, err := c.Password.GenerateHashFromPassword(c.Config.AppServer.InitialRootAdminPassword)
		if err != nil {
			c.Logger.Error("hashing error", slog.Any("error", err))
			return err
		}
		m := &u_d.User{
			ID:                    uid,
			FirstName:             "Root",
			LastName:              "Administrator",
			Name:                  "Root Administrator",
			LexicalName:           "Administrator, Root",
			Email:                 c.Config.AppServer.InitialRootAdminEmail,
			Phone:                 "+1 (647) 967-2269",
			Country:               "Canada",
			Region:                "Ontario",
			City:                  "London",
			PostalCode:            "N6H 0B7",
			AddressLine1:          "1828 Blue Heron Dr Unit #26",
			PasswordHash:          passwordHash,
			PasswordHashAlgorithm: c.Password.AlgorithmName(),
			Role:                  u_d.UserRoleRoot,
			WasEmailVerified:      true,
			CreatedAt:             time.Now(),
			ModifiedAt:            time.Now(),
			Status:                u_d.UserStatusActive,
			AgreeTOS:              true,
		}
		err = c.UserStorer.Create(ctx, m)
		if err != nil {
			c.Logger.Error("database create error", slog.Any("error", err))
			return err
		}
		c.Logger.Debug("Root user created.",
			slog.Any("_id", m.ID),
			slog.String("name", m.Name),
			slog.String("email", m.Email),
			slog.String("password_hash_algorithm", m.PasswordHashAlgorithm),
			slog.String("password_hash", m.PasswordHash))
	} else {
		c.Logger.Debug("Root user already exists, skipping creation.")
	}
	return nil
}

// createInitialOrg function creates the initial organization and the administrator if not previously created.
func (c *OrganizationControllerImpl) createInitialOrg(ctx context.Context) (*o_d.Organization, error) {
	id, err := primitive.ObjectIDFromHex(c.Config.AppServer.InitialOrgID)
	if err != nil {
		c.Logger.Error("primitive object id create from hex failed error",
			slog.Any("error", err),
			slog.Any("InitialOrgID", c.Config.AppServer.InitialOrgID))
		return nil, err
	}

	c.Logger.Debug("Looking up organization id from environment variable",
		slog.Any("InitialOrgID", c.Config.AppServer.InitialOrgID))
	u, err := c.OrganizationStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database check if exists error",
			slog.Any("error", err),
			slog.Any("InitialOrgID", c.Config.AppServer.InitialOrgID))
		return nil, err
	}

	// If already exists then exit without doing anything.
	if u != nil {
		c.Logger.Debug("Found organization, skipping creation",
			slog.Any("InitialOrgID", c.Config.AppServer.InitialOrgID))
		// Return `nil` instead of organization so the code will terminate
		// with creating more user accounts.
		return nil, nil
	}
	c.Logger.Debug("Organization does not exist, creating now",
		slog.Any("InitialOrgID", c.Config.AppServer.InitialOrgID))

	org := &o_d.Organization{
		ID:           id,
		Name:         c.Config.AppServer.InitialOrgName,
		Email:        "mike@bp8fitness.com",
		Phone:        "+1 (647) 967-2269",
		Country:      "Canada",
		Region:       "Ontario",
		City:         "London",
		PostalCode:   "N6H 0B7",
		AddressLine1: "1828 Blue Heron Dr Unit #26",
		Status:       o_d.OrganizationStatusActive,
		Type:         o_d.GymType,
		CreatedAt:    time.Now(),
		ModifiedAt:   time.Now(),
	}
	err = c.OrganizationStorer.Create(ctx, org)
	if err != nil {
		c.Logger.Error("database create error",
			slog.Any("error", err),
			slog.Any("InitialOrgID", c.Config.AppServer.InitialOrgID),
			slog.Any("InitialOrgName", c.Config.AppServer.InitialOrgName))
		return nil, err
	}
	c.Logger.Debug("Organizational created.",
		slog.Any("_id", org.ID),
		slog.String("name", org.Name),
		slog.Any("created_by_user_id", org.CreatedByUserID),
		slog.Time("created_at", org.CreatedAt),
		slog.Time("modified_at", org.ModifiedAt))
	return org, nil
}

// createInitialOrgAdmin function creates the initial organization and the administrator if not previously created.
func (c *OrganizationControllerImpl) createInitialOrgAdmin(ctx context.Context, org *o_d.Organization) error {

	//
	// STEP 1: Check to see if we have our default organization.
	//

	doesExist, err := c.UserStorer.CheckIfExistsByEmail(ctx, c.Config.AppServer.InitialOrgAdminEmail)
	if err != nil {
		c.Logger.Error("database check if exists error", slog.Any("error", err))
		return err
	}
	if doesExist == false {
		//
		// STEP 2: Create our default organization administrator.
		//

		c.Logger.Debug("No org admin user detected, proceeding to create now...")
		passwordHash, err := c.Password.GenerateHashFromPassword(c.Config.AppServer.InitialOrgAdminPassword)
		if err != nil {
			c.Logger.Error("hashing error", slog.Any("error", err))
			return err
		}
		m := &u_d.User{
			OrganizationID:        org.ID,
			OrganizationName:      org.Name,
			ID:                    primitive.NewObjectID(),
			FirstName:             "Organization",
			LastName:              "Administrator",
			Name:                  "Organization Administrator",
			LexicalName:           "Administrator, Organization",
			Email:                 c.Config.AppServer.InitialOrgAdminEmail,
			PasswordHash:          passwordHash,
			PasswordHashAlgorithm: c.Password.AlgorithmName(),
			Role:                  u_d.UserRoleAdmin,
			WasEmailVerified:      true,
			CreatedAt:             time.Now(),
			ModifiedAt:            time.Now(),
			Status:                u_d.UserStatusActive,
			AgreeTOS:              true,
		}
		err = c.UserStorer.Create(ctx, m)
		if err != nil {
			c.Logger.Error("database create error", slog.Any("error", err))
			return err
		}
		c.Logger.Debug("Org admin user created.",
			slog.Any("_id", m.ID),
			slog.String("name", m.Name),
			slog.String("email", m.Email),
			slog.String("password_hash_algorithm", m.PasswordHashAlgorithm),
			slog.String("password_hash", m.PasswordHash))
	} else {
		c.Logger.Debug("Org admin user already exists, skipping creation.")
	}
	return nil
}
