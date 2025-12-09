package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Kavehrafie.com/imgopt/internal/config"
	"github.com/Kavehrafie.com/imgopt/internal/processor"
	"github.com/Kavehrafie.com/imgopt/internal/storage"
)

func main() {
	cfg := config.Load()

	log.Printf("DEBUG: Storage Type configured as: '%s'", cfg.StorageType)

	var store storage.Provider
	var err error

	switch cfg.StorageType {
	case "bunny":
		log.Println("DEBUG: Initializing Bunny.net storage...")
		store, err = storage.NewBunnyStorage(cfg)
	case "b2":
		log.Println("DEBUG: Initializing Backblaze B2 storage...")
		store, err = storage.NewB2Storage(context.Background(), cfg)
	default:
		log.Fatalf("Unknown storage type: %s", cfg.StorageType)
	}

	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	proc := processor.NewService()

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			// This should be handled by the specific handler, but just in case
			return
		}

		// Parse path to get key and options
		// Format: /path/to/image.jpg/w_200/h_300/fit_cover/crop_smart
		path := strings.TrimPrefix(r.URL.Path, "/")
		parts := strings.Split(path, "/")

		var width, height int
		var fit, crop string
		var key string

		// Iterate backwards to find options
		endIdx := len(parts)
		for i := len(parts) - 1; i >= 0; i-- {
			part := parts[i]
			if strings.HasPrefix(part, "w_") {
				wStr := strings.TrimPrefix(part, "w_")
				width, _ = strconv.Atoi(wStr)
				endIdx = i
			} else if strings.HasPrefix(part, "h_") {
				hStr := strings.TrimPrefix(part, "h_")
				height, _ = strconv.Atoi(hStr)
				endIdx = i
			} else if strings.HasPrefix(part, "fit_") {
				fit = strings.TrimPrefix(part, "fit_")
				endIdx = i
			} else if strings.HasPrefix(part, "crop_") {
				crop = strings.TrimPrefix(part, "crop_")
				endIdx = i
			} else {
				// Not a recognized option, assume part of the key
				break
			}
		}

		if endIdx == 0 {
			http.Error(w, "Missing image key", http.StatusBadRequest)
			return
		}

		key = strings.Join(parts[:endIdx], "/")

		if key == "" {
			http.Error(w, "Missing 'key' parameter", http.StatusBadRequest)
			return
		}

		log.Printf("DEBUG: Processing request for key: '%s' (width: %d, height: %d, fit: %s, crop: %s)", key, width, height, fit, crop)

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		reader, contentType, err := store.GetFile(ctx, key)
		if err != nil {
			log.Printf("Error getting file %s: %v", key, err)
			http.Error(w, "Image not found", http.StatusNotFound)
			return
		}
		defer reader.Close()

		// If no processing needed, stream directly
		if width == 0 && height == 0 && fit == "" && crop == "" {
			w.Header().Set("Content-Type", contentType)
			w.Header().Set("Cache-Control", "public, max-age=31536000") // Cache for 1 year
			// Copy stream
			io.Copy(w, reader)
			return
		}

		// Process image
		// If width/height are 0, Resize will just decode and encode, effectively stripping metadata/optimizing slightly
		processedImg, err := proc.Resize(reader, width, height, fit, crop, contentType)

		if err != nil {
			log.Printf("Error processing image %s: %v", key, err)
			http.Error(w, "Failed to process image", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Cache-Control", "public, max-age=31536000") // Cache for 1 year
		w.Header().Set("Content-Length", strconv.Itoa(len(processedImg)))
		w.Write(processedImg)
	})

	addr := ":" + cfg.Port
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
