package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func getExtension(mimeType string) string {
	return strings.Split(mimeType, "/")[1]
}

func (cfg apiConfig) getFilePath(id uuid.UUID, ext string) string {
	return filepath.Join(fmt.Sprintf("%v/%v.%v", cfg.assetsRoot, id.String(), ext))

}

func (cfg apiConfig) getFileURL(filePath string) string {
	return fmt.Sprintf("http://localhost:%v/%v", cfg.port, filePath)

}
