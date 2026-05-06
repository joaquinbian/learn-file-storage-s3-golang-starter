package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadVideo(w http.ResponseWriter, r *http.Request) {

	const MAX_MEMORY = 1 << 30
	r.Body = http.MaxBytesReader(w, r.Body, MAX_MEMORY)

	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	video, err := cfg.db.GetVideo(videoID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting video", err)
		return
	}

	if video.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "Forbidden request", err)
		return
	}

	file, fileData, err := r.FormFile("video")

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting video from form", err)
		return
	}

	defer file.Close()

	mimeType := fileData.Header.Get("Content-Type")
	mediaType, err := validateVideoMediaType(mimeType)

	ext := mediaTypeToExt(mediaType)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "video key must be a video type", err)
		return

	}

	tmpFile, err := os.CreateTemp("", "tubely-video.mp4")

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating temp file", err)
		return
	}

	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, file)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error copying temp file", err)
		return
	}
	_, err = tmpFile.Seek(0, io.SeekStart)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error seeking temp file", err)
		return
	}
	key := getRandomPath() + ext

	obj := s3.PutObjectInput{
		Bucket:      &cfg.s3Bucket,
		Key:         &key,
		Body:        tmpFile,
		ContentType: &mimeType,
	}

	_, err = cfg.s3Client.PutObject(r.Context(), &obj)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error saving video to S3", err)
		return
	}

	videoURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", cfg.s3Bucket, cfg.s3Region, key)

	video.VideoURL = &videoURL

	err = cfg.db.UpdateVideo(video)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating video to DB", err)
		return
	}
	respondWithJSON(w, http.StatusOK, video)

}
