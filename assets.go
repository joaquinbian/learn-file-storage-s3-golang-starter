package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"mime"
	"os"
	"os/exec"
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

func getVideoAspectRatio(filepath string) (string, error) {
	var fileStream FileStream
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filepath)

	var res = &bytes.Buffer{}
	cmd.Stdout = res
	err := cmd.Run()

	if err != nil {
		return "", errors.New("error executing command")
	}

	if err := json.Unmarshal(res.Bytes(), &fileStream); err != nil {
		return "", errors.New("error unmarshalling file stream")
	}

	ratio := getAspectRatio(fileStream.Streams[0].Width, fileStream.Streams[0].Height)

	return translateAspectRatio(ratio), nil
}

// AspectRatio calcula el aspect ratio (ancho/alto) de una imagen dados su width y height.
// Devuelve 0 si height es 0 para evitar división por cero.
func getAspectRatio(width, height int) string {
	if height == 0 {
		return "other"
	}
	ratio := float64(width / height)
	const tolerance = 0.02
	if math.Abs(ratio-(16/9)) < tolerance {
		return "16:9"
	} else if math.Abs(ratio-(9/16)) < tolerance {
		return "9:16"
	} else {
		return "other"
	}
}

func translateAspectRatio(ratio string) string {
	if strings.Compare(ratio, "16:9") == 0 {
		return "landscape"
	} else if strings.Compare(ratio, "9:16") == 0 {
		return "portrait"
	} else {
		return "other"
	}
}

func processVideoForFastStart(filepath string) (string, error) {

	output := fmt.Sprintf("%s.%s", filepath, "processing")

	cmd := exec.Command("ffmpeg", "-i", filepath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", output)

	err := cmd.Run()

	if err != nil {
		fmt.Print(err)
		return "", err
	}

	return output, nil
}
