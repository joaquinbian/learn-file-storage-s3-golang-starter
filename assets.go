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

func (cfg apiConfig) getFilePath(path string, ext string) string {
	return filepath.Join(fmt.Sprintf("%v/%v.%v", cfg.assetsRoot, path, ext))

}

func (cfg apiConfig) getFileURL(filePath string) string {
	return fmt.Sprintf("http://localhost:%v/%v", cfg.port, filePath)

}

func getS3AssetURL(bucket, region, key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)
}

func generateRandID() (string, error) {
	randID := make([]byte, 32)

	_, err := rand.Read(randID)

	if err != nil {
		return "", errors.New("error reading random id")
	}

	encoding := base64.RawStdEncoding.EncodeToString(randID)

	return encoding, nil

}
