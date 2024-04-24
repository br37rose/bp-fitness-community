package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/adapter/storage/mongodb"
	eq_d "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/equipment/datastore"
	e_d "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/exercise/datastore"
	o_d "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/organization/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config"
)

func init() {
	rootCmd.AddCommand(importBP8ExerciseCmd)
}

var importBP8ExerciseCmd = &cobra.Command{
	Use:   "import_bp8exercise",
	Short: "Import the exercise from csv file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("starting...")
		ctx := context.Background()
		loggerp := slog.Default()
		cfg := config.New()
		mc := mongodb.NewStorage(cfg, loggerp)
		orgDS := o_d.NewDatastore(cfg, loggerp, mc)
		eqDS := eq_d.NewDatastore(cfg, loggerp, mc)
		eDS := e_d.NewDatastore(cfg, loggerp, mc)
		if err := importBP8Exercise(ctx, cfg, loggerp, orgDS, eqDS, eDS); err != nil {
			loggerp.Error("failed initializing", slog.Any("err", err))
		}
	},
}

type BP8CSVExerciseRow struct {
	ID          int
	Category    string
	Type        string
	Gender      string
	Location    string
	Equipment   string
	Movement    string
	Description string
	VideoName   string
	VideoURL    string
}

var bp8GenderMap = map[string]int8{
	"Unspecified": 1,
	"Male":        2,
	"Female":      3,
	"Both":        4,
}

func importBP8Exercise(ctx context.Context, cfg *config.Conf, logger *slog.Logger, oDS o_d.OrganizationStorer, eqDS eq_d.EquipmentStorer, eDS e_d.ExerciseStorer) error {
	logger.Debug("getting  default organization")
	orgID, err := primitive.ObjectIDFromHex(cfg.AppServer.InitialOrgID)
	if err != nil {
		logger.Error("primitive object id create from hex failed error",
			slog.Any("error", err),
			slog.Any("InitialOrgID", cfg.AppServer.InitialOrgID))
		return err
	}
	org, err := oDS.GetByID(ctx, orgID)

	logger.Debug("importing exercises...")

	// Open the CSV file
	file, err := os.Open("./static/bp8_exercises_2024.csv")
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read the CSV headers
	headers, err := reader.Read()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	logger.Debug("skipping headers", slog.Any("headers", headers))

	// Read and process each line
	for {
		line, err := reader.Read()
		if err != nil {
			break // Stop reading when end of file is reached or an error occurs
		}

		// Create an Exercise instance and populate its fields
		exercise := &BP8CSVExerciseRow{
			ID:          parseInt(line[0]),
			Category:    line[1],
			Type:        line[2],
			Gender:      line[3],
			Location:    line[4],
			Equipment:   line[5],
			Movement:    line[6],
			Description: line[7],
			VideoName:   line[8],
			VideoURL:    line[9],
		}

		logger.Debug("Opened", slog.Any("exercise", exercise))
		log.Println("row ID", exercise.ID, "VideoURL:", exercise.VideoURL)

		if err := importBP8ExerciseIntoDB(ctx, logger, eqDS, eDS, org, exercise); err != nil {
			return err
		}
	}
	return nil
}

func importBP8ExerciseIntoDB(ctx context.Context, logger *slog.Logger, eqDS eq_d.EquipmentStorer, es e_d.ExerciseStorer, org *o_d.Organization, row *BP8CSVExerciseRow) error {
	////
	//// Step 1: Find related equipment.
	////

	equipments := make([]*e_d.ExerciseEquipment, 0)
	equipmentArr := strings.Split(row.Equipment, "/")
	for _, name := range equipmentArr {
		e, err := eqDS.GetByName(ctx, strings.Title(name))
		if err != nil {
			return err
		}
		if e != nil {
			ee := &e_d.ExerciseEquipment{
				ID:     e.ID,
				Name:   e.Name,
				No:     e.No,
				Status: e.Status,
			}
			equipments = append(equipments, ee)
		}
	}

	////
	//// Step 2: Movement type.
	////

	mt := strings.TrimSpace(row.Movement)
	movementType, ok := movementTypeMap[mt]
	if ok == false {
		fmt.Println("skipping b/c not found movement type:", row.Movement)
		fmt.Println(movementTypeMap)
		fmt.Println(movementTypeMap[row.Movement])
		fmt.Println()
		os.Exit(0)
		return nil
	}

	////
	//// Step 3: Category
	////

	category, ok := categoryMap[row.Category]
	if ok == false {
		fmt.Println("skipping b/c not category type:", row.Category)
		fmt.Println()
		os.Exit(0)
		return nil
	}

	////
	//// Step 4: Define the unique ID of our exercise.
	////
	eID := primitive.NewObjectID()

	////
	//// Step 5: Tag
	////
	bp8TagID, _ := primitive.ObjectIDFromHex("659c1f2be0b31ef6ec8e288d") // Hard-coded from `import_tag` file.

	et := &e_d.ExerciseTag{
		ID:             bp8TagID,
		OrganizationID: org.ID,
		Text:           "BP8 Fitness",
		Status:         e_d.ExerciseStatusActive,
	}

	////
	//// Step 6: Create our entity.
	////

	e := &e_d.Exercise{
		OrganizationID:   org.ID,
		OrganizationName: org.Name,
		ID:               eID,
		Type:             e_d.ExerciseTypeCustom,
		Gender:           row.Gender,
		MovementType:     movementType,
		Category:         category,
		Name:             row.VideoName,
		Description:      row.Description,
		VideoType:        e_d.ExerciseVideoTypeS3,
		VideoName:        row.VideoName,
		VideoURL:         row.VideoURL,
		Status:           e_d.ExerciseStatusActive,
		Tags:             []*e_d.ExerciseTag{et},
	}

	if err := es.Create(ctx, e); err != nil {
		return err
	}
	logger.Debug("Imported exercise", slog.Any("id", e.ID), slog.Any("row_id", row.ID))
	fmt.Println("Imported exercise ID:", e.ID, "via Row ID:", row.ID)
	fmt.Println()
	return nil
}
