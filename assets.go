package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func validateImageMediaType(mimeType string) (string, error) {

	VALID_MEDIA_TYPES := []string{"image/png", "image/jpeg"}
	mediatype, _, err := mime.ParseMediaType(mimeType)
	if err != nil {
		return "", err
	}

	if !slices.Contains(VALID_MEDIA_TYPES, mediatype) {
		return "", errors.New("file is not an image")
	}

	return mediatype, nil
}

func validateVideoMediaType(mimeType string) (string, error) {
	VALID_MEDIA_TYPES := []string{"video/mp4", "video/mov"}
	mediatype, _, err := mime.ParseMediaType(mimeType)
	if err != nil {
		return "", err
	}

	if !slices.Contains(VALID_MEDIA_TYPES, mediatype) {
		return "", errors.New("file is not a video")
	}

	return mediatype, nil
}

func getAssetPath(assetID string, mediaType string) string {
	ext := mediaTypeToExt(mediaType)

	return fmt.Sprintf("%s%s", assetID, ext)
}
func (cfg *apiConfig) getAssetFilepath(path string) string {
	return filepath.Join(cfg.assetsRoot, path)
}

func (cfg *apiConfig) getAssetURL(path string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, path)
}
func mediaTypeToExt(mimeType string) string {
	parts := strings.Split(mimeType, "/")
	if len(parts) != 2 {
		return ".bin"
	}
	return "." + parts[1]
}

func getRandomPath() string {
	url := make([]byte, 32)
	rand.Read(url)

	encodedURL := base64.RawURLEncoding.EncodeToString(url)
	return encodedURL
}
