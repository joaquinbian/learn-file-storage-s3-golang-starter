package main

import (
	"errors"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/google/uuid"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func getExtension(mimeType string) (string, error) {
	mediaType, _, err := mime.ParseMediaType(mimeType)

	if err != nil {
		return "", err
	}

	return mediaType, nil
}

func validateAndGetExtension(mimeType string, validExts []string) (string, error) {
	mediaType, err := getExtension(mimeType)
	if err != nil {
		return "", nil
	}

	if !slices.Contains(validExts, mediaType) {
		return "", errors.New("invalid media type")
	}

	ext := strings.Split(mediaType, "/")[1]
	return ext, nil
}

func (cfg apiConfig) getFilePath(id uuid.UUID, ext string) string {
	return filepath.Join(fmt.Sprintf("%v/%v.%v", cfg.assetsRoot, id.String(), ext))

}

func (cfg apiConfig) getFileURL(filePath string) string {
	return fmt.Sprintf("http://localhost:%v/%v", cfg.port, filePath)

}
