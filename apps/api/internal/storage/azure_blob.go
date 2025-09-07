package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
)

type AzureBlobClient struct {
	client        *azblob.Client
	credential    *azblob.SharedKeyCredential
	containerName string
	accountName   string
}

func NewAzureBlobClient(containerName string) (*AzureBlobClient, error) {
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	if accountName == "" || accountKey == "" {
		return nil, fmt.Errorf("AZURE_STORAGE_ACCOUNT_NAME and AZURE_STORAGE_ACCOUNT_KEY must be set")
	}

	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create shared key credential: %w", err)
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob client: %w", err)
	}

	return &AzureBlobClient{
		client:        client,
		credential:    cred, // Store the credential for signing SAS tokens
		containerName: containerName,
		accountName:   accountName,
	}, nil
}

func (c *AzureBlobClient) generateSASURL(ctx context.Context, key string, expiry time.Duration, permissions sas.BlobPermissions, contentType string) (string, error) {
	expiryTime := time.Now().UTC().Add(expiry)

	sasValues := sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		ExpiryTime:    expiryTime,
		ContainerName: c.containerName,
		BlobName:      key,
		Permissions:   permissions.String(),
		ContentType:   contentType,
	}

	sasToken, err := sasValues.SignWithSharedKey(c.credential)
	if err != nil {
		return "", fmt.Errorf("failed to sign SAS token with shared key: %w", err)
	}

	finalURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s",
		c.accountName, c.containerName, key, sasToken.Encode())

	return finalURL, nil
}

func (c *AzureBlobClient) GenerateUploadURL(ctx context.Context, key string, contentType string, allowedTypes []string, maxUploadSize int, expiry time.Duration) (string, error) {
	if len(allowedTypes) > 0 {
		isAllowed := false
		for _, t := range allowedTypes {
			if t == contentType {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			return "", fmt.Errorf("ContentType is not allowed: %s", contentType)
		}
	}
	permissions := sas.BlobPermissions{Create: true, Write: true}
	return c.generateSASURL(ctx, key, expiry, permissions, contentType)
}

// func (c *AzureBlobClient) GenerateDownloadURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
// 	permissions := sas.BlobPermissions{Read: true}
// 	return c.generateSASURL(ctx, key, expiry, permissions)
// }
//
// func (c *AzureBlobClient) Delete(ctx context.Context, key string) error {
// 	_, err := c.client.DeleteBlob(ctx, c.containerName, key, nil)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete blob %s: %w", key, err)
// 	}
// 	return nil
// }
