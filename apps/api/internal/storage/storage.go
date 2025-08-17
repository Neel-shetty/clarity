package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage interface {
	GenerateUploadURL(ctx context.Context, key string, contentType string, allowedTypes []string, maxUploadSize int, expirySecs time.Duration) (*s3.PresignedPostRequest, error)
	GenerateDownloadURL(ctx context.Context, key string, expiry time.Duration) (string, error)
	Delete(ctx context.Context, key string) error
}
