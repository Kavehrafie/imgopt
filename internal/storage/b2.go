package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/Backblaze/blazer/b2"
	"github.com/Kavehrafie.com/imgopt/internal/config"
)

type Provider interface {
	GetFile(ctx context.Context, key string) (io.ReadCloser, string, error)
}

type B2Storage struct {
	client *b2.Client
	bucket *b2.Bucket
}

func NewB2Storage(ctx context.Context, cfg *config.Config) (*B2Storage, error) {
	client, err := b2.NewClient(ctx, cfg.B2AccountID, cfg.B2ApplicationKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create B2 client: %w", err)
	}

	bucket, err := client.Bucket(ctx, cfg.B2BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket: %w", err)
	}

	return &B2Storage{
		client: client,
		bucket: bucket,
	}, nil
}

func (s *B2Storage) GetFile(ctx context.Context, key string) (io.ReadCloser, string, error) {
	obj := s.bucket.Object(key)

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get file attributes for key '%s': %w", key, err)
	}

	reader := obj.NewReader(ctx)

	return reader, attrs.ContentType, nil
}
