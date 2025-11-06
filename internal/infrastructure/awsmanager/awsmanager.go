package awsmanager

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"

	appconfig "boilerplate-golang/internal/infrastructure/config"
)

var (
	s3Client   *s3.Client
	bucketName string
	region     string
)

// Init initializes the AWS S3 client with credentials from config
func Init() error {
	cfg := appconfig.Get()

	// Skip initialization if AWS is not configured
	if cfg.AWS.AccessKeyID == "" || cfg.AWS.SecretAccessKey == "" || cfg.AWS.S3Bucket == "" {
		log.Println("awsmanager: AWS not configured, skipping S3 client initialization")
		return nil
	}

	ctx := context.Background()

	// Load AWS configuration
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.AWS.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWS.AccessKeyID,
			cfg.AWS.SecretAccessKey,
			"", // session token (optional)
		)),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3Client = s3.NewFromConfig(awsCfg)
	bucketName = cfg.AWS.S3Bucket
	region = cfg.AWS.Region

	log.Printf("awsmanager: initialized with bucket '%s' in region '%s'", bucketName, region)
	return nil
}

// UploadFileResult contains the result of a file upload
type UploadFileResult struct {
	URL      string
	Key      string
	Bucket   string
	Size     int64
	MimeType string
}

// UploadFile uploads a file to S3 and returns upload details
func UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (*UploadFileResult, error) {
	if s3Client == nil {
		return nil, fmt.Errorf("S3 client not initialized")
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s-%s%s", uuid.New().String(), time.Now().Format("20060102150405"), ext)
	key := filepath.Join(folder, filename)

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload to S3
	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Generate URL
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
		bucketName,
		region,
		key,
	)

	log.Printf("awsmanager: uploaded file to %s", url)

	return &UploadFileResult{
		URL:      url,
		Key:      key,
		Bucket:   bucketName,
		Size:     header.Size,
		MimeType: contentType,
	}, nil
}

// DeleteFile removes a file from S3
func DeleteFile(ctx context.Context, key string) error {
	if s3Client == nil {
		return fmt.Errorf("S3 client not initialized")
	}

	_, err := s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	log.Printf("awsmanager: deleted file %s", key)
	return nil
}

// GeneratePresignedURL creates a temporary signed URL for private file access
func GeneratePresignedURL(ctx context.Context, key string, expiresIn time.Duration) (string, error) {
	if s3Client == nil {
		return "", fmt.Errorf("S3 client not initialized")
	}

	presignClient := s3.NewPresignClient(s3Client)

	presignedURL, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiresIn
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.URL, nil
}
