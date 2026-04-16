package storage

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
)

type BlobService interface {
	UploadKTP(ctx context.Context, fileBytes []byte, filename string, contentType string) (string, error)
}

type azureBlobService struct {
	client        *azblob.Client
	containerName string
}

func NewAzureBlobService(connectionString string) (BlobService, error) {
	// Untuk demo, jika dummy string diberikan, kita return struct yang "berpura-pura" sukses
	if connectionString == "" || connectionString == "DefaultEndpointsProtocol=https;AccountName=dummy;AccountKey=dummy;EndpointSuffix=core.windows.net" {
		return &azureBlobService{client: nil, containerName: "blobhghal2026"}, nil
	}

	client, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob client: %w", err)
	}

	return &azureBlobService{
		client:        client,
		containerName: "blobhghal2026",
	}, nil
}

func (s *azureBlobService) UploadKTP(ctx context.Context, fileBytes []byte, filename string, contentType string) (string, error) {
	// Prefix / key attachment: "uploads"
	blobName := path.Join("uploads", fmt.Sprintf("%d_%s", time.Now().Unix(), filename))

	// Jika ini di tahap "belajar" dan belum set connection string asli:
	if s.client == nil {
		// Mock sukses, return dummy URL
		fmt.Printf("MOCK UPLOAD SUCCESS: Container=%s, BlobName=%s, Size=%d bytes\n", s.containerName, blobName, len(fileBytes))
		return fmt.Sprintf("https://dummy.blob.core.windows.net/%s/%s", s.containerName, blobName), nil
	}

	// 1. Dapatkan referensi container, dan buat bila belum ada, atau abaikan bila sudah ada
	// Ini bisa di-improve dengan mengecek exists() atau ditaruh waktu init, tapi untuk kemudahan:
	s.client.CreateContainer(ctx, s.containerName, nil)

	// 2. Metadata dummy sesuai request (3-4 metadata yang relevan)
	metadata := map[string]*string{
		"uploadedBy":   Ptr("test_user_only"),
		"documentType": Ptr("KTP"),
		"status":       Ptr("pending_verification"),
	}

	// 3. Upload file
	options := &azblob.UploadBufferOptions{
		Metadata: metadata,
		HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType: &contentType,
		},
	}

	_, err := s.client.UploadBuffer(ctx, s.containerName, blobName, fileBytes, options)
	if err != nil {
		return "", fmt.Errorf("failed to upload blob: %w", err)
	}

	fileURL := fmt.Sprintf("%s%s/%s", s.client.URL(), s.containerName, blobName)
	return fileURL, nil
}

func Ptr[T any](v T) *T {
	return &v
}
