package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/adapter/storage/mongodb"
	e_d "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/exercise/datastore"
	o_d "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/organization/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config"
)

func init() {
	rootCmd.AddCommand(importExerciseCmd)
}

var importExerciseCmd = &cobra.Command{
	Use:   "import_exercise",
	Short: "Import the exercise from directory",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		loggerp := slog.Default()
		cfg := config.New()
		mc := mongodb.NewStorage(cfg, loggerp)
		orgDS := o_d.NewDatastore(cfg, loggerp, mc)
		eDS := e_d.NewDatastore(cfg, loggerp, mc)
		if err := importExercise(ctx, cfg, loggerp, orgDS, eDS); err != nil {
			loggerp.Error("failed initializing", slog.Any("err", err))
		}
	},
}

type CSVExerciseRow struct {
	ID             int
	UUID           string
	Category       string
	Type           string
	Name           string
	Gender         string
	MovementType   string
	Description    string
	VideoType      string
	VideoS3Key     string
	VideoName      string
	VideoURL       string
	ThumbnailType  string
	ThumbnailS3Key string
	ThumbnailName  string
	ThumbnailURL   string
	State          string
}

var categoryMap = map[string]int8{
	"Anterior":               1,
	"Anti-Extension":         2,
	"Anti-extension":         2,
	"Anti-Lateral Flexion":   3,
	"Anti-Rotation":          4,
	"Band":                   5,
	"Bar":                    6,
	"Barbell":                7,
	"Biceps":                 8,
	"Bodyweight":             9,
	"Cable":                  10,
	"Cable or band":          11,
	"Carries":                12,
	"Combination":            13,
	"Dumbbell":               14,
	"Dumbell":                15,
	"Dynamic":                16,
	"Flexion":                17,
	"Ground-Based Exercises": 18,
	"Hip Belt":               19,
	"Lunge":                  20,
	"Plate":                  21,
	"Posterior":              22,
	"Pushup":                 23,
	"Ring":                   24,
	"Static":                 25,
	"Triceps":                26,
}

var movementTypeMap = map[string]int8{
	"Arms":                       1,
	"Core Work":                  2,
	"Core work":                  2,
	"Corrective Work":            3,
	"Hip-Hinge":                  4,
	"Horizontal Pressing":        5,
	"Horizontal pressing":        5,
	"Horizontal Pushing":         5,
	"Horizontal pushing":         5,
	"Horizontal Pulling":         6,
	"Horizontal pulling":         6,
	"Jump":                       7,
	"Jumps":                      7,
	"Single-Leg":                 8,
	"Squat":                      9,
	"Squats":                     9,
	"Vertical Pressing":          10,
	"Vertical pressing":          10,
	"Vertical Pushing":           10,
	"Vertical Pulling":           11,
	"Vertical pulling":           11,
	"Warmups & Mobility Fillers": 12,
	"Warmup":                     12,
	"Work":                       13,
}

func importExercise(ctx context.Context, cfg *config.Conf, logger *slog.Logger, oDS o_d.OrganizationStorer, eDS e_d.ExerciseStorer) error {
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
	file, err := os.Open("./static/exercises_2023.csv")
	if err != nil {
		fmt.Println("Error:", err)
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
		exercise := &CSVExerciseRow{
			ID:             parseInt(line[0]),
			UUID:           line[1],
			Category:       line[2],
			Type:           line[3],
			Name:           line[4],
			Gender:         line[5],
			MovementType:   line[6],
			Description:    line[7],
			VideoType:      line[8],
			VideoS3Key:     line[9],
			VideoName:      line[10],
			VideoURL:       line[11],
			ThumbnailType:  line[12],
			ThumbnailS3Key: line[13],
			ThumbnailName:  line[14],
			ThumbnailURL:   line[15],
			State:          line[16],
		}

		if err := importExerciseIntoDB(ctx, logger, eDS, org, exercise); err != nil {
			return err
		}
	}
	return nil
}

func parseInt(s string) int {
	i := 0
	fmt.Sscanf(s, "%d", &i)
	return i
}

func ConvertCSVRowToExercise(org *o_d.Organization, row *CSVExerciseRow) *e_d.Exercise {
	////
	//// Step 1: Movement type.
	////

	movementType, ok := movementTypeMap[row.MovementType]
	if ok == false {
		fmt.Println("skipping b/c not found movement type:", row.MovementType)
		fmt.Println()
		os.Exit(0)
		return nil
	}

	////
	//// Step 2: Category
	////

	category, ok := categoryMap[row.Category]
	if ok == false {
		fmt.Println("skipping b/c not category type:", row.Category)
		fmt.Println()
		os.Exit(0)
		return nil
	}

	////
	//// Step 3: Define the unique ID of our exercise.
	////

	eID := primitive.NewObjectID()

	////
	//// Step 4: Tag
	////
	bp8TagID, _ := primitive.ObjectIDFromHex("659c1f2be0b31ef6ec8e288e") // Hard-coded from `import_tag` file.

	et := &e_d.ExerciseTag{
		ID:             bp8TagID,
		OrganizationID: org.ID,
		Text:           "Pure Nutrition",
		Status:         e_d.ExerciseStatusActive,
	}

	////
	//// Step 5: Create our entity.
	////

	e := &e_d.Exercise{
		ID:                    eID,
		Type:                  e_d.ExerciseTypeSystem,
		Name:                  row.Name,
		Gender:                row.Gender,
		MovementType:          movementType,
		Category:              category,
		Description:           row.Description,
		VideoType:             e_d.ExerciseVideoTypeVimeo,
		VideoName:             row.VideoName,
		VideoURL:              row.VideoURL,
		ThumbnailType:         e_d.ExerciseThumbnailTypeS3,
		ThumbnailObjectKey:    fmt.Sprintf("exercises/thumbnails/%v", row.ThumbnailName),
		ThumbnailObjectExpiry: time.Now(),
		ThumbnailName:         row.ThumbnailName,
		ThumbnailObjectURL:    row.ThumbnailURL,
		Status:                e_d.ExerciseStatusActive,
		Tags:                  []*e_d.ExerciseTag{et},
	}

	e.OrganizationID = org.ID
	e.OrganizationName = org.Name
	e.ThumbnailURL = strings.Replace(e.ThumbnailURL, "brandonprust", "bp8fitnesscommunity", 1)
	e.ThumbnailURL = strings.Replace(e.ThumbnailURL, "exercise/thumbnail", "exercises/thumbnails", 1)

	return e
}

func importExerciseIntoDB(ctx context.Context, logger *slog.Logger, es e_d.ExerciseStorer, org *o_d.Organization, exercise *CSVExerciseRow) error {
	e := ConvertCSVRowToExercise(org, exercise)
	if err := es.Create(ctx, e); err != nil {
		return err
	}
	logger.Debug("Imported exercise", slog.Any("id", e.ID))
	return nil
}
