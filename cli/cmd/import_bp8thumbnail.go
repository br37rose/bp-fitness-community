package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/adapter/storage/mongodb"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/adapter/storage/s3"
	eq_d "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/equipment/datastore"
	e_d "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/exercise/datastore"
	o_d "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/organization/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/provider/uuid"
)

/*
    DEVELOPERS NOTE:
	YOU NEED TO HAVE ACCESS OT `FFmpeg` PROGRAM.

    PICK AN OPTION:
	(1) https://superuser.com/questions/624561/install-ffmpeg-on-os-x
	(2) https://json2video.com/how-to/ffmpeg-course/install-ffmpeg-mac.html

    MAC OPTION:
	1. DOWNLOAD `FFmpeg` HERE
	2. INSTALL ROSETTA 2. SEE: https://apple.stackexchange.com/questions/408375/zsh-bad-cpu-type-in-executable
	3. TRY RUNNING THIS APPLICATION BUT APPLE WILL REJECT...
	4. ... THEN DO THIS: https://www.lifewire.com/fix-developer-cannot-be-verified-error-5183898
	5. RUN THIS APPLICATION AGAIN AND IT WILL WORK.
*/

func init() {
	importBP8ThumbnailCmd.Flags().StringVarP(&customPrefix, "prefix", "p", "", "Prefix to append our root folder for all our files")
	rootCmd.AddCommand(importBP8ThumbnailCmd)
}

var importBP8ThumbnailCmd = &cobra.Command{
	Use:   "import_bp8thumbnail",
	Short: "Import the bp8 thumbnail for the exercises in the csv file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("\033[H\033[2J") // Clear screen
		log.Println("starting...")
		ctx := context.Background()
		loggerp := slog.Default()
		cfg := config.New()
		uuidp := uuid.NewProvider()
		s3 := s3.NewStorage(cfg, loggerp, uuidp)
		mc := mongodb.NewStorage(cfg, loggerp)
		orgDS := o_d.NewDatastore(cfg, loggerp, mc)
		eqDS := eq_d.NewDatastore(cfg, loggerp, mc)
		eDS := e_d.NewDatastore(cfg, loggerp, mc)

		// ALGORITHM:
		// 1. Iterate through all the CSV rows.
		// 2. For each row, find exercise with matching `ThumbnailURL` values.
		// 3. If match found then retrieve exercise record.
		// 4. Get the video filename from the `ThumbnailURL`.
		// 5. Find the video filename in the directory.
		// 6. Submit it to S3.
		// 7. Save the exercise with the new s3 uploaded file key.

		if err := importBP8Thumbnail(ctx, cfg, loggerp, s3, orgDS, eqDS, eDS); err != nil {
			loggerp.Error("failed initializing", slog.Any("err", err))
		}
	},
}

func importBP8Thumbnail(ctx context.Context, cfg *config.Conf, logger *slog.Logger, s3 s3.S3Storager, oDS o_d.OrganizationStorer, eqDS eq_d.EquipmentStorer, eDS e_d.ExerciseStorer) error {
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

		if err := importBP8ThumbnailIntoDB(ctx, logger, s3, eqDS, eDS, org, exercise); err != nil {
			return err
		}
	}
	return nil
}

func importBP8ThumbnailIntoDB(ctx context.Context, logger *slog.Logger, s3 s3.S3Storager, eqDS eq_d.EquipmentStorer, ex e_d.ExerciseStorer, org *o_d.Organization, row *BP8CSVExerciseRow) error {

	exercise, err := ex.GetByVideoURL(ctx, row.VideoURL)
	if err != nil {
		return err
	}
	if exercise == nil {
		fmt.Println("skipping row", row.ID, " b/c not found exercise in database for video URL:", row.VideoURL)
		fmt.Println()
		return nil
	}

	// DEVELOPERS NOTE: Format the https encoded url into something we can work with.
	encodedPath := strings.Replace(row.VideoURL, "https://www.dropbox.com/home/BP8%20finished%203rd%20shoot/", "", 1)
	filename := strings.Replace(encodedPath, "Cardio?preview=", "", -1)
	filename = strings.Replace(filename, "Upper%20body?preview=", "", -1)
	filename = strings.Replace(filename, "Core?preview=", "", -1)
	filename = strings.Replace(filename, "Lower%20body?preview=", "", -1)
	filename = strings.Replace(filename, "full%20body?preview=", "", -1)
	filename = strings.Replace(filename, "+", " ", -1)

	// unique case 1
	filename = strings.Replace(filename, "https://www.dropbox.com/preview/BP8%20finished%203rd%20shoot/Cardio/", "", -1)
	filename = strings.Replace(filename, "?context=content_suggestions&role=personal", "", -1)

	// unique case 2
	filename = strings.Replace(filename, "https://www.dropbox.com/scl/fo/lfw67r1b6r05ymg86hiwh/h/Upper%20body?dl=0&preview=", " ", -1)

	// unique case 3
	filename = strings.Replace(filename, "https://www.dropbox.com/home/Shoot%201%20Rough%20Cuts?preview=", " ", -1)

	// unique case 4
	filename = strings.Replace(filename, "https://www.dropbox.com/home/Shoot%201%20Rough%20Cuts?preview=", " ", -1)

	if strings.Contains(filename, "https") {
		fmt.Println("skipping row", row.ID, " b/c corrupt video url:", row.VideoURL)
		fmt.Println()
		// os.Exit(0)
		return nil
	}
	if filename == "" {
		fmt.Println("skipping row", row.ID, " b/c empty video url")
		fmt.Println()
		return nil
	}

	// For debugging purposes only.
	// fmt.Println("--->", exercise)
	// fmt.Println("encodedPath --->", encodedPath)
	// fmt.Println("filename --->", filename)

	matchChan := make(chan string)

	root := "./static"
	go func() {
		// Start the walk in a goroutine
		err := filepath.Walk(root, visit(filename, matchChan))
		if err != nil {
			fmt.Printf("error walking the path %v: %v\n", root, err)
		}
		close(matchChan)
	}()

	// Read from the channel
	for matchingFile := range matchChan {
		fmt.Printf("Found matching file at: %s\n", matchingFile)
		thumbnailFilename := strings.Replace(filename, ".mp4", "thumbnail.jpg", -1)
		thumbnailFilepath, err := createThumbnail(matchingFile)
		if err != nil {
			return err
		}

		if err := uploadThumbnailToS3(ctx, s3, ex, exercise, thumbnailFilepath, thumbnailFilename); err != nil {
			return err
		}

		// os.Exit(0) // For debugging purposes only.
	}

	return nil
}

func createThumbnail(videoPath string) (string, error) {
	// Generate the output file path for our thumbnail.
	outputPath := strings.Replace(videoPath, ".mp4", "thumbnail.jpg", -1)

	// Run FFmpeg command to extract a thumbnail
	cmd := exec.Command("./ffmpeg", "-i", videoPath, "-ss", "00:00:01", "-vframes", "1", outputPath)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error running FFmpeg: %v", err)
	}
	return outputPath, nil
}

func uploadThumbnailToS3(ctx context.Context, s3 s3.S3Storager, ex e_d.ExerciseStorer, e *e_d.Exercise, localFilepath string, localFilename string) error {
	//
	// Open and read file.
	//

	// Read the file contents into a []byte slice
	fileContent, err := ioutil.ReadFile(localFilepath)
	if err != nil {
		log.Fatal(err)
	}

	//
	// Upload the content to S3.
	//

	// Generate the key of our exercise video upload.
	videoObjectKey := fmt.Sprintf("org/%v/exercise/%v/%v", e.OrganizationID.Hex(), e.ID.Hex(), localFilename)

	// DEVELOPERS NOTE: If we are using any sort of prefix then apply it now.
	if customPrefix != "" {
		videoObjectKey = fmt.Sprintf("%v_%v", customPrefix, videoObjectKey)
	}

	// Generate a presigned URL for today.
	expiryDur := time.Hour * 12
	videoObjectURL, err := s3.GetPresignedURL(ctx, videoObjectKey, expiryDur)

	fmt.Println("Uploading", videoObjectKey, "for exercise ID", e.ID)

	// Upload to our s3 bucket.
	if err := s3.UploadContent(ctx, videoObjectKey, fileContent); err != nil {
		return err
	}

	fmt.Println("Uploaded successfully thumnail", localFilename, "for exercise ID", e.ID)

	// Update the exercise.
	e.ThumbnailType = e_d.ExerciseThumbnailTypeS3
	e.ThumbnailName = localFilename
	e.ThumbnailObjectKey = videoObjectKey
	e.ThumbnailObjectURL = videoObjectURL
	e.ThumbnailObjectExpiry = time.Now().Add(expiryDur)
	if err := ex.UpdateByID(context.Background(), e); err != nil {
		return err
	}
	return nil
}

// Auto-generated comment for change 15
