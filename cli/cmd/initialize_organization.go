package cmd

import (
	"context"
	"log/slog"
	"time"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/adapter/storage/mongodb"
	o_d "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/organization/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config"
)

func init() {
	rootCmd.AddCommand(initializeOrganizationCmd)
}

var initializeOrganizationCmd = &cobra.Command{
	Use:   "initialize_organization",
	Short: "Initialize the default organization",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		loggerp := slog.Default()
		cfg := config.New()
		mc := mongodb.NewStorage(cfg, loggerp)
		orgDS := o_d.NewDatastore(cfg, loggerp, mc)
		if err := initializeOrganization(ctx, cfg, loggerp, orgDS); err != nil {
			loggerp.Error("failed initializing", slog.Any("err", err))
		}
	},
}

func initializeOrganization(ctx context.Context, cfg *config.Conf, logger *slog.Logger, orgS o_d.OrganizationStorer) error {
	logger.Debug("initiailizing organization...")

	id, err := primitive.ObjectIDFromHex(cfg.AppServer.InitialOrgID)
	if err != nil {
		logger.Error("primitive object id create from hex failed error",
			slog.Any("error", err),
			slog.Any("InitialOrgID", cfg.AppServer.InitialOrgID))
		return err
	}
	logger.Debug("Looking up organization id from environment variable",
		slog.Any("InitialOrgID", cfg.AppServer.InitialOrgID))
	u, err := orgS.GetByID(ctx, id)
	if err != nil {
		logger.Error("database check if exists error",
			slog.Any("error", err),
			slog.Any("InitialOrgID", cfg.AppServer.InitialOrgID))
		return err
	}

	// If already exists then exit without doing anything.
	if u != nil {
		logger.Debug("Found organization, skipping creation",
			slog.Any("InitialOrgID", cfg.AppServer.InitialOrgID))
		// Return `nil` instead of organization so the code will terminate
		// with creating more user accounts.
		return nil
	}
	logger.Debug("Organization does not exist, creating now",
		slog.Any("InitialOrgID", cfg.AppServer.InitialOrgID))

	org := &o_d.Organization{
		ID:           id,
		Name:         cfg.AppServer.InitialOrgName,
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
	err = orgS.UpsertByID(ctx, org)
	if err != nil {
		logger.Error("database create error",
			slog.Any("error", err),
			slog.Any("InitialOrgID", cfg.AppServer.InitialOrgID),
			slog.Any("InitialOrgName", cfg.AppServer.InitialOrgName))
		return err
	}
	logger.Debug("Organizational created.",
		slog.Any("_id", org.ID),
		slog.String("name", org.Name),
		slog.Any("created_by_user_id", org.CreatedByUserID),
		slog.Time("created_at", org.CreatedAt),
		slog.Time("modified_at", org.ModifiedAt))

	logger.Debug("initiailized organization")
	return nil
}

// Auto-generated comment for change 12
