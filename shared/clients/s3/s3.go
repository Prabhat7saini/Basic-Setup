package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/zap"
)

type Config struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	SignedURLExpiry int
	MaxFileSize     int
}

type s3Client struct {
	client   *s3.S3
	uploader *s3manager.Uploader
	logger   *zap.Logger
	config   Config
}

type UploadResult struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

type GeneratePresignedURLResult struct {
	URL string `json:"url"`
}

type Client interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader) (*UploadResult, error)
	GeneratePresignedURL(ctx context.Context, key string) (*GeneratePresignedURLResult, error)
}

func NewClient(cfg Config, logger *zap.Logger) (Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &s3Client{
		client:   s3.New(sess),
		uploader: s3manager.NewUploader(sess),
		logger:   logger.With(zap.String("component", "s3-client")),
		config:   cfg,
	}, nil
}

func (s *s3Client) generateKey(filename string) string {
	base := filepath.Base(filename)
	extension := filepath.Ext(base)
	name := strings.TrimSuffix(base, extension)

	// Clean the name
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "'", "")
	name = strings.ReplaceAll(name, "\"", "")

	// final key: <name>_<timestamp><ext>
	return fmt.Sprintf("%s_%d%s", name, time.Now().UnixNano(), extension)
}

func (s *s3Client) UploadFile(ctx context.Context, file *multipart.FileHeader) (*UploadResult, error) {
	if file == nil {
		s.logger.Error("nil file provided")
		return nil, fmt.Errorf("file not provided")
	}
	maxFileSize := s.config.MaxFileSize
	if file.Size > int64(maxFileSize) {
		s.logger.Warn("file too large",
			zap.String("filename", file.Filename),
			zap.Int64("size", file.Size))
		return nil, fmt.Errorf("file size exceeds  limit")
	}

	src, err := file.Open()
	if err != nil {
		s.logger.Error("failed to open file",
			zap.String("filename", file.Filename),
			zap.Error(err))
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		s.logger.Error("failed to read file",
			zap.String("filename", file.Filename),
			zap.Error(err))
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	key := s.generateKey(file.Filename)
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(fileBytes)
	}

	_, err = s.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(s.config.BucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		s.logger.Error("failed to upload file",
			zap.String("key", key),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	s.logger.Info("file uploaded successfully",
		zap.String("key", key),
	)

	url, err := s.GeneratePresignedURL(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("file uploaded, but failed to generate URL: %w", err)
	}

	return &UploadResult{Key: key, URL: url.URL}, nil
}

func (s *s3Client) GeneratePresignedURL(ctx context.Context, key string) (*GeneratePresignedURLResult, error) {
	if key == "" {
		s.logger.Error("empty key provided")
		return nil, fmt.Errorf("key not provided")
	}

	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.config.BucketName),
		Key:    aws.String(key),
	})

	url, err := req.Presign(time.Duration(s.config.SignedURLExpiry) * time.Minute)
	if err != nil {
		s.logger.Error("failed to generate presigned URL",
			zap.String("key", key),
			zap.Error(err))

		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return &GeneratePresignedURLResult{URL: url}, nil
}
