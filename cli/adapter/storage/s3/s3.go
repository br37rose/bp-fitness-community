package s3 // Special thanks via https://docs.digitalocean.com/products/spaces/resources/s3-sdk-examples/

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"log/slog"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/provider/uuid"
)

type S3Storager interface {
	UploadContent(ctx context.Context, objectKey string, content []byte) error
	UploadContentFromMulipart(ctx context.Context, objectKey string, file multipart.File) error
	BucketExists(ctx context.Context, bucketName string) (bool, error)
	GetDownloadablePresignedURL(ctx context.Context, key string, duration time.Duration) (string, error)
	GetPresignedURL(ctx context.Context, key string, duration time.Duration) (string, error)
	DeleteByKeys(ctx context.Context, key []string) error
	Cut(ctx context.Context, sourceObjectKey string, destinationObjectKey string) error
	Copy(ctx context.Context, sourceObjectKey string, destinationObjectKey string) error
	GetBinaryData(ctx context.Context, objectKey string) (io.ReadCloser, error)
	DownloadToLocalfile(ctx context.Context, objectKey string, filePath string) (string, error)
	ListAllObjects(ctx context.Context) (*s3.ListObjectsOutput, error)
	FindMatchingObjectKey(s3Objects *s3.ListObjectsOutput, partialKey string) string
}

type s3Storager struct {
	S3Client      *s3.Client
	PresignClient *s3.PresignClient
	UUID          uuid.Provider
	Logger        *slog.Logger
	BucketName    string
}

// NewStorage connects to a specific S3 bucket instance and returns a connected
// instance structure.
func NewStorage(appConf *c.Conf, logger *slog.Logger, uuidp uuid.Provider) S3Storager {
	// DEVELOPERS NOTE:
	// How can I use the AWS SDK v2 for Go with DigitalOcean Spaces? via https://stackoverflow.com/a/74284205
	logger.Debug("s3 initializing...")

	// STEP 1: initialize the custom `endpoint` we will connect to.
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: appConf.AWS.Endpoint,
		}, nil
	})

	// STEP 2: Configure.
	sdkConfig, err := config.LoadDefaultConfig(
		context.TODO(), config.WithRegion(appConf.AWS.Region),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(appConf.AWS.AccessKey, appConf.AWS.SecretKey, "")),
	)
	if err != nil {
		log.Fatal(err) // We need to crash the program at start to satisfy google wire requirement of having no errors.
	}

	// STEP 3\: Load up s3 instance.
	s3Client := s3.NewFromConfig(sdkConfig)

	// Create our storage handler.
	s3Storage := &s3Storager{
		S3Client:      s3Client,
		PresignClient: s3.NewPresignClient(s3Client),
		Logger:        logger,
		UUID:          uuidp,
		BucketName:    appConf.AWS.BucketName,
	}

	logger.Debug("s3 checking remote connection...")

	// STEP 4: Connect to the s3 bucket instance and confirm that bucket exists.
	doesExist, err := s3Storage.BucketExists(context.TODO(), appConf.AWS.BucketName)
	if err != nil {
		log.Fatal(err) // We need to crash the program at start to satisfy google wire requirement of having no errors.
	}
	if !doesExist {
		log.Fatal("bucket name does not exist") // We need to crash the program at start to satisfy google wire requirement of having no errors.
	}

	logger.Debug("s3 initialized")

	// Return our s3 storage handler.
	return s3Storage
}

func (s *s3Storager) UploadContent(ctx context.Context, objectKey string, content []byte) error {
	_, err := s.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(content),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *s3Storager) UploadContentFromMulipart(ctx context.Context, objectKey string, file multipart.File) error {
	// Create the S3 upload input parameters
	params := &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	}

	// Perform the file upload to S3
	_, err := s.S3Client.PutObject(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (s *s3Storager) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	// Note: https://docs.aws.amazon.com/code-library/latest/ug/go_2_s3_code_examples.html#actions

	_, err := s.S3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	exists := true
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				log.Printf("Bucket %v is available.\n", bucketName)
				exists = false
				err = nil
			default:
				log.Printf("Either you don't have access to bucket %v or another error occurred. "+
					"Here's what happened: %v\n", bucketName, err)
			}
		}
	}

	return exists, err
}

func (s *s3Storager) GetDownloadablePresignedURL(ctx context.Context, key string, duration time.Duration) (string, error) {
	// DEVELOPERS NOTE:
	// AWS S3 Bucket — presigned URL APIs with Go (2022) via https://ronen-niv.medium.com/aws-s3-handling-presigned-urls-2718ab247d57

	presignedUrl, err := s.PresignClient.PresignGetObject(context.Background(),
		&s3.GetObjectInput{
			Bucket:                     aws.String(s.BucketName),
			Key:                        aws.String(key),
			ResponseContentDisposition: aws.String("attachment"), // This field allows the file to download it directly from your browser
		},
		s3.WithPresignExpires(duration))
	if err != nil {
		return "", err
	}
	return presignedUrl.URL, nil
}

func (s *s3Storager) GetPresignedURL(ctx context.Context, objectKey string, duration time.Duration) (string, error) {
	// DEVELOPERS NOTE:
	// AWS S3 Bucket — presigned URL APIs with Go (2022) via https://ronen-niv.medium.com/aws-s3-handling-presigned-urls-2718ab247d57

	presignedUrl, err := s.PresignClient.PresignGetObject(context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(s.BucketName),
			Key:    aws.String(objectKey),
		},
		s3.WithPresignExpires(duration))
	if err != nil {
		return "", err
	}
	return presignedUrl.URL, nil
}

func (s *s3Storager) DeleteByKeys(ctx context.Context, objectKeys []string) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var objectIds []types.ObjectIdentifier
	for _, key := range objectKeys {
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	}
	_, err := s.S3Client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(s.BucketName),
		Delete: &types.Delete{Objects: objectIds},
	})
	if err != nil {
		log.Printf("Couldn't delete objects from bucket %v. Here's why: %v\n", s.BucketName, err)
	}
	return err
}

func (s *s3Storager) Cut(ctx context.Context, sourceObjectKey string, destinationObjectKey string) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second) // Increase timout so it runs longer then usual to handle this unique case.
	defer cancel()

	_, copyErr := s.S3Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(s.BucketName),
		CopySource: aws.String(s.BucketName + "/" + sourceObjectKey),
		Key:        aws.String(destinationObjectKey),
	})
	if copyErr != nil {
		s.Logger.Error("Failed to copy object:", slog.Any("copyErr", copyErr))
		return copyErr
	}

	s.Logger.Debug("Object copied successfully.")

	// Delete the original object
	_, deleteErr := s.S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(sourceObjectKey),
	})
	if deleteErr != nil {
		s.Logger.Error("Failed to delete original object:", slog.Any("deleteErr", deleteErr))
		return deleteErr
	}

	s.Logger.Debug("Original object deleted.")

	return nil
}

func (s *s3Storager) Copy(ctx context.Context, sourceObjectKey string, destinationObjectKey string) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second) // Increase timout so it runs longer then usual to handle this unique case.
	defer cancel()

	_, copyErr := s.S3Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(s.BucketName),
		CopySource: aws.String(s.BucketName + "/" + sourceObjectKey),
		Key:        aws.String(destinationObjectKey),
	})
	if copyErr != nil {
		s.Logger.Error("Failed to copy object:", slog.Any("copyErr", copyErr))
		return copyErr
	}

	s.Logger.Debug("Object copied successfully.")

	return nil
}

// GetBinaryData function will return the binary data for the particular key.
func (s *s3Storager) GetBinaryData(ctx context.Context, objectKey string) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(objectKey),
	}

	s3object, err := s.S3Client.GetObject(ctx, input)
	if err != nil {
		return nil, err
	}
	return s3object.Body, nil
}

func (s *s3Storager) DownloadToLocalfile(ctx context.Context, objectKey string, filePath string) (string, error) {
	responseBin, err := s.GetBinaryData(ctx, objectKey)
	if err != nil {
		return filePath, err
	}
	out, err := os.Create(filePath)
	if err != nil {
		return filePath, err
	}
	defer out.Close()

	_, err = io.Copy(out, responseBin)
	if err != nil {
		return "", err
	}
	return filePath, err
}

func (s *s3Storager) ListAllObjects(ctx context.Context) (*s3.ListObjectsOutput, error) {
	input := &s3.ListObjectsInput{
		Bucket: aws.String(s.BucketName),
	}

	objects, err := s.S3Client.ListObjects(ctx, input)
	if err != nil {
		return nil, err
	}

	return objects, nil
}

// Function will iterate over all the s3 objects to match the partial key with
// the actual key found in the S3 bucket.
func (s *s3Storager) FindMatchingObjectKey(s3Objects *s3.ListObjectsOutput, partialKey string) string {
	for _, obj := range s3Objects.Contents {

		match := strings.Contains(*obj.Key, partialKey)

		// If a match happens then it means we have found the ACTUAL KEY in the
		// s3 objects inside the bucket.
		if match == true {
			return *obj.Key
		}
	}
	return ""
}
