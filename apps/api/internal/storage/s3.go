package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	client  *s3.Client
	presign *s3.PresignClient
	bucket  string
}

func NewS3Client(ctx context.Context, bucket string) (*S3Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	return &S3Client{
		client:  s3Client,
		presign: s3.NewPresignClient(s3Client),
		bucket:  bucket,
	}, nil
}

// GenerateUploadURL creates a presigned POST URL that allows a client to upload a file
// directly to S3. It enforces a maximum file size and can optionally enforce a specific
// content type.
//
// This function is the recommended way to handle secure, client-side uploads as it
// offloads the bandwidth from the server and provides fine-grained control over uploads.
//
// Parameters:
//   - ctx: The context for the AWS API call.
//   - key: The full object key (path and filename) where the file will be stored in S3.
//   - contentType: The MIME type of the file the client intends to upload (e.g., "application/pdf").
//   - allowedTypes: A slice of allowed MIME types. If this slice is not nil and not empty,
//     the function will validate the provided contentType against this list and generate a
//     policy that enforces an exact match. If the slice is nil or empty, the function will
//     not add a content type restriction to the policy, allowing any file type to be uploaded.
//   - maxUploadSize: The maximum allowable file size in bytes.
//   - expirySecs: The duration for which the generated URL will be valid.
//
// Returns:
//   - A pointer to a PresignedPost struct containing the URL and the required form fields
//     for the upload.
//   - An error if the content type is disallowed or if the URL generation fails.
//
// Example (Enforcing specific file types):
//
//	allowed := []string{"application/pdf", "application/msword"}
//	// contentType comes from the client's request, e.g., "application/pdf"
//	presignedPost, err := s3Client.GenerateUploadURL(ctx, "my-document.pdf", "application/pdf", allowed, 10*1024*1024, 15*time.Minute)
//
// Example (Allowing any file type):
//
//	// Pass nil for allowedTypes
//	presignedPost, err := s3Client.GenerateUploadURL(ctx, "any-file.zip", "application/zip", nil, 10*1024*1024, 15*time.Minute)
func (s *S3Client) GenerateUploadURL(ctx context.Context, key string, contentType string, allowedTypes []string, maxUploadSize int, expirySecs time.Duration) (*s3.PresignedPostRequest, error) {
	input := &s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	}
	conditions := [][]interface{}{
		{"content-length-range", 1, maxUploadSize},
	}

	if len(allowedTypes) > 0 {
		isAllowed := false
		for _, t := range allowedTypes {
			if t == contentType {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			return &s3.PresignedPostRequest{}, fmt.Errorf("disallowed content type: %s", contentType)
		}

		conditions = append(conditions, []interface{}{"eq", "$Content-Type", contentType})
	}

	post, err := s.presign.PresignPostObject(ctx, input, func(opts *s3.PresignPostOptions) {
		opts.Expires = expirySecs
		// opts.Conditions = append(opts.Conditions,
		// 	[]interface{}{"content-length-range", 1, maxUploadSize},
		// )
		// opts.Conditions = append(opts.Conditions,
		// 	conditions,
		// )
		for _, c := range conditions {
			opts.Conditions = append(opts.Conditions, c)
		}
	})
	if err != nil {
		return &s3.PresignedPostRequest{}, fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	return post, nil
}

func (s *S3Client) GenerateDownloadURL(ctx context.Context, key string, expirySecs time.Duration) (string, error) {
	req, err := s.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expirySecs
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned download URL: %w", err)
	}
	return req.URL, nil
}

func (s *S3Client) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	return nil
}
