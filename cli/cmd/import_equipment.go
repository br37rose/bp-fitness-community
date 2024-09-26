package cmd

import (
	"context"
	"log"
	"log/slog"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/adapter/storage/mongodb"
	eq_d "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/equipment/datastore"
	o_d "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/organization/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config"
)

func init() {
	rootCmd.AddCommand(importEquipmentCmd)
}

var importEquipmentCmd = &cobra.Command{
	Use:   "import_equipment",
	Short: "Imports all standard equipment",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		loggerp := slog.Default()
		cfg := config.New()
		mc := mongodb.NewStorage(cfg, loggerp)
		orgDS := o_d.NewDatastore(cfg, loggerp, mc)
		vcDS := eq_d.NewDatastore(cfg, loggerp, mc)

		loggerp.Debug("getting default organization")
		orgID, err := primitive.ObjectIDFromHex(cfg.AppServer.InitialOrgID)
		if err != nil {
			loggerp.Error("primitive object id create from hex failed error",
				slog.Any("error", err),
				slog.Any("InitialOrgID", cfg.AppServer.InitialOrgID))
			log.Fatal(err)
		}
		org, err := orgDS.GetByID(ctx, orgID)

		loggerp.Debug("importing equipment...")

		if err := importEquipment(ctx, cfg, loggerp, orgDS, vcDS, org); err != nil {
			loggerp.Error("failed initializing", slog.Any("err", err))
		}
	},
}

func importEquipment(ctx context.Context, cfg *config.Conf, logger *slog.Logger, oDS o_d.OrganizationStorer, vcDS eq_d.EquipmentStorer, org *o_d.Organization) error {
	data := []*eq_d.Equipment{
		&eq_d.Equipment{
			OrganizationID:   org.ID,
			OrganizationName: org.Name,
			ID:               primitive.NewObjectID(),
			Name:             "Bench",
			No:               1,
			Status:           eq_d.StatusActive,
		},
		&eq_d.Equipment{
			OrganizationID:   org.ID,
			OrganizationName: org.Name,
			ID:               primitive.NewObjectID(),
			Name:             "Box",
			No:               2,
			Status:           eq_d.StatusActive,
		},
		&eq_d.Equipment{
			OrganizationID:   org.ID,
			OrganizationName: org.Name,
			ID:               primitive.NewObjectID(),
			Name:             "Stationary Bike",
			No:               3,
			Status:           eq_d.StatusActive,
		},
		&eq_d.Equipment{
			OrganizationID:   org.ID,
			OrganizationName: org.Name,
			ID:               primitive.NewObjectID(),
			Name:             "Bosu Ball",
			No:               4,
			Status:           eq_d.StatusActive,
		},
		&eq_d.Equipment{
			OrganizationID:   org.ID,
			OrganizationName: org.Name,
			ID:               primitive.NewObjectID(),
			Name:             "Cables",
			No:               5,
			Status:           eq_d.StatusActive,
		},
		&eq_d.Equipment{
			OrganizationID:   org.ID,
			OrganizationName: org.Name,
			ID:               primitive.NewObjectID(),
			Name:             "Rowing Machine",
			No:               6,
			Status:           eq_d.StatusActive,
		},
		&eq_d.Equipment{
			OrganizationID:   org.ID,
			OrganizationName: org.Name,
			ID:               primitive.NewObjectID(),
			Name:             "Band",
			No:               7,
			Status:           eq_d.StatusActive,
		},
		&eq_d.Equipment{
			OrganizationID:   org.ID,
			OrganizationName: org.Name,
			ID:               primitive.NewObjectID(),
			Name:             "Sled",
			No:               8,
			Status:           eq_d.StatusActive,
		},
		&eq_d.Equipment{
			OrganizationID:   org.ID,
			OrganizationName: org.Name,
			ID:               primitive.NewObjectID(),
			Name:             "Barbell",
			No:               9,
			Status:           eq_d.StatusActive,
		},
		&eq_d.Equipment{
			OrganizationID:   org.ID,
			OrganizationName: org.Name,
			ID:               primitive.NewObjectID(),
			Name:             "Dumbbells",
			No:               10,
			Status:           eq_d.StatusActive,
		},
		&eq_d.Equipment{
			OrganizationID:   org.ID,
			OrganizationName: org.Name,
			ID:               primitive.NewObjectID(),
			Name:             "Swiss Ball",
			No:               11,
			Status:           eq_d.StatusActive,
		},
		&eq_d.Equipment{
			OrganizationID:   org.ID,
			OrganizationName: org.Name,
			ID:               primitive.NewObjectID(),
			Name:             "Dead Ball",
			No:               13,
			Status:           eq_d.StatusActive,
		},
		&eq_d.Equipment{
			OrganizationID:   org.ID,
			OrganizationName: org.Name,
			ID:               primitive.NewObjectID(),
			Name:             "Kettle Bell",
			No:               14,
			Status:           eq_d.StatusActive,
		},
	}
	for _, datum := range data {
		if err := vcDS.Create(context.Background(), datum); err != nil {
			return err
		}
	}

	return nil
}

// Auto-generated comment for change 2
