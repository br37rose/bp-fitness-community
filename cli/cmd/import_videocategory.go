package cmd

import (
	"context"
	"log"
	"log/slog"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/adapter/storage/mongodb"
	o_d "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/organization/datastore"
	vc_d "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/videocategory/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config"
)

func init() {
	rootCmd.AddCommand(ivcCmd)
}

var ivcCmd = &cobra.Command{
	Use:   "import_videocategory",
	Short: "Imports all the video categories",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		loggerp := slog.Default()
		cfg := config.New()
		mc := mongodb.NewStorage(cfg, loggerp)
		orgDS := o_d.NewDatastore(cfg, loggerp, mc)
		vcDS := vc_d.NewDatastore(cfg, loggerp, mc)

		loggerp.Debug("getting default organization")
		orgID, err := primitive.ObjectIDFromHex(cfg.AppServer.InitialOrgID)
		if err != nil {
			loggerp.Error("primitive object id create from hex failed error",
				slog.Any("error", err),
				slog.Any("InitialOrgID", cfg.AppServer.InitialOrgID))
			log.Fatal(err)
		}
		org, err := orgDS.GetByID(ctx, orgID)

		loggerp.Debug("importing video categories...")

		if err := importVideoCategory(ctx, cfg, loggerp, orgDS, vcDS, org); err != nil {
			loggerp.Error("failed initializing", slog.Any("err", err))
		}
	},
}

func importVideoCategory(ctx context.Context, cfg *config.Conf, logger *slog.Logger, oDS o_d.OrganizationStorer, vcDS vc_d.VideoCategoryStorer, org *o_d.Organization) error {
	data := []*vc_d.VideoCategory{
		&vc_d.VideoCategory{
			OrganizationID:   org.ID,
			OrganizationName: org.Name,
			ID:               primitive.NewObjectID(),
			Name:             "Off-Ice",
			No:               1,
			Status:           vc_d.StatusActive,
		},
		&vc_d.VideoCategory{
			OrganizationID:   org.ID,
			OrganizationName: org.Name,
			ID:               primitive.NewObjectID(),
			Name:             "On-Ice",
			No:               2,
			Status:           vc_d.StatusActive,
		},
		&vc_d.VideoCategory{
			OrganizationID:   org.ID,
			OrganizationName: org.Name,
			ID:               primitive.NewObjectID(),
			Name:             "Exclusive Content",
			No:               3,
			Status:           vc_d.StatusActive,
		},
	}
	for _, datum := range data {
		if err := vcDS.Create(context.Background(), datum); err != nil {
			return err
		}
	}

	return nil
}
