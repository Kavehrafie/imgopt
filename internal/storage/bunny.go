package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/Kavehrafie.com/imgopt/internal/config"
	bunnystorage "github.com/l0wl3vel/bunny-storage-go-sdk"
)

type BunnyStorage struct {
	client bunnystorage.Client
	zone   string
}

func NewBunnyStorage(cfg *config.Config) (*BunnyStorage, error) {
	var endpointStr string
	switch cfg.BunnyEndpoint {
	case "ny":
		endpointStr = bunnystorage.ENDPOINT_NEW_YORK_US
	case "la":
		endpointStr = bunnystorage.ENDPOINT_LOS_ANGELES_US
	case "sg":
		endpointStr = bunnystorage.ENDPOINT_SINGAPORE_SG
	case "syd":
		endpointStr = bunnystorage.ENDPOINT_SYDNEY_SYD
	case "uk":
		endpointStr = bunnystorage.ENDPOINT_LONDON_UK
	case "se":
		endpointStr = bunnystorage.ENDPOINT_STOCKHOLM_SE
	default:
		endpointStr = bunnystorage.ENDPOINT_FALKENSTEIN_DE
	}

	u := url.URL{
		Scheme: "https",
		Host:   endpointStr,
	}

	// Prefer ReadOnlyKey for downloads
	apiKey := cfg.BunnyAccessKey
	if cfg.BunnyReadOnlyKey != "" {
		apiKey = cfg.BunnyReadOnlyKey
	}

	client := bunnystorage.NewClient(u, apiKey)

	return &BunnyStorage{
		client: client,
		zone:   cfg.BunnyZoneName,
	}, nil
}

func (s *BunnyStorage) GetFile(ctx context.Context, key string) (io.ReadCloser, string, error) {
	// Sanitize key
	key = strings.Trim(key, "/")

	// Construct full path including zone
	fullPath := fmt.Sprintf("/%s/%s", s.zone, key)

	fmt.Printf("DEBUG: Bunny Download Path: '%s'\n", fullPath)

	content, err := s.client.Download(fullPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download file from Bunny: %w", err)
	}

	ext := filepath.Ext(key)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return io.NopCloser(bytes.NewReader(content)), contentType, nil
}
