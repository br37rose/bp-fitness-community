package cmd

import (
	"context"
	"log"
	"log/slog"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/adapter/storage/mongodb"
	o_d "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/organization/datastore"
	tag_d "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/tag/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config"
)

func init() {
	rootCmd.AddCommand(importTagCmd)
}

var importTagCmd = &cobra.Command{
	Use:   "import_tag",
	Short: "Imports all standard tag",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		loggerp := slog.Default()
		cfg := config.New()
		mc := mongodb.NewStorage(cfg, loggerp)
		orgDS := o_d.NewDatastore(cfg, loggerp, mc)
		vcDS := tag_d.NewDatastore(cfg, loggerp, mc)

		loggerp.Debug("getting default organization")
		orgID, err := primitive.ObjectIDFromHex(cfg.AppServer.InitialOrgID)
		if err != nil {
			loggerp.Error("primitive object id create from hex failed error",
				slog.Any("error", err),
				slog.Any("InitialOrgID", cfg.AppServer.InitialOrgID))
			log.Fatal(err)
		}
		org, err := orgDS.GetByID(ctx, orgID)

		loggerp.Debug("importing tag...")

		if err := importTag(ctx, cfg, loggerp, orgDS, vcDS, org); err != nil {
			loggerp.Error("failed initializing", slog.Any("err", err))
		}
	},
}

func importTag(ctx context.Context, cfg *config.Conf, logger *slog.Logger, oDS o_d.OrganizationStorer, vcDS tag_d.TagStorer, org *o_d.Organization) error {
	// Hard code values.
	id1, _ := primitive.ObjectIDFromHex("659c1f2be0b31ef6ec8e288d")
	id2, _ := primitive.ObjectIDFromHex("659c1f2be0b31ef6ec8e288e")

	data := []*tag_d.Tag{
		&tag_d.Tag{
			OrganizationID:   org.ID,
			OrganizationName: org.Name,
			ID:               id1,
			Text:             "BP8 Fitness",
			Status:           tag_d.TagStatusActive,
		},
		&tag_d.Tag{
			OrganizationID:   org.ID,
			OrganizationName: org.Name,
			ID:               id2,
			Text:             "Pure Nutrition",
			Status:           tag_d.TagStatusActive,
		},
	}
	for _, datum := range data {
		if err := vcDS.Create(context.Background(), datum); err != nil {
			return err
		}
	}

	return nil
}
