package processor

import (
	"fmt"
	"io"

	"github.com/h2non/bimg"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Resize(r io.Reader, width, height int, fit, crop string, format string) ([]byte, error) {
	buffer, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %w", err)
	}

	if width == 0 && height == 0 {
		return buffer, nil
	}

	newImage := bimg.NewImage(buffer)

	// Check if image format is supported
	if newImage.Type() == "unknown" {
		if len(buffer) == 0 {
			return nil, fmt.Errorf("image buffer is empty")
		}
		// Debug: return first 100 bytes as string to see if it's HTML or error message
		preview := string(buffer)
		if len(preview) > 100 {
			preview = preview[:100]
		}
		return nil, fmt.Errorf("unknown image format. Buffer preview: %q", preview)
	}

	options := bimg.Options{
		Width:  width,
		Height: height,
	}

	// Handle 'fit' parameter
	switch fit {
	case "cover":
		options.Crop = true
	case "contain":
		options.Embed = true
	case "fill":
		options.Force = true
	case "inside":
		// Default behavior, preserves aspect ratio, no crop, no padding (resulting image might be smaller than requested)
		options.Crop = false
		options.Embed = false
		options.Force = false
	default:
		// Default to cover if not specified, or inside?
		// User asked for "crop, fit too".
		// Let's default to 'cover' if crop is specified, otherwise 'inside'?
		// Or just default to 'inside' (standard resize).
		if crop != "" {
			options.Crop = true
		}
	}

	// Handle 'crop' parameter (gravity)
	if options.Crop {
		switch crop {
		case "top":
			options.Gravity = bimg.GravityNorth
		case "right":
			options.Gravity = bimg.GravityEast
		case "bottom":
			options.Gravity = bimg.GravitySouth
		case "left":
			options.Gravity = bimg.GravityWest
		case "smart":
			options.Gravity = bimg.GravitySmart
		default:
			options.Gravity = bimg.GravityCentre
		}
	}

	processed, err := newImage.Process(options)
	if err != nil {
		return nil, fmt.Errorf("failed to process image: %w. Detected Type: %s", err, newImage.Type())
	}

	return processed, nil
}

func (s *Service) SmartCrop(r io.Reader, width, height int, format string) ([]byte, error) {
	buffer, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %w", err)
	}

	newImage := bimg.NewImage(buffer)

	options := bimg.Options{
		Width:   width,
		Height:  height,
		Crop:    true,
		Gravity: bimg.GravitySmart,
	}

	processed, err := newImage.Process(options)
	if err != nil {
		return nil, fmt.Errorf("failed to process image: %w", err)
	}

	return processed, nil
}
