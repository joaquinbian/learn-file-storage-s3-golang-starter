package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// TODO: implement the upload here

	const maxMemory = 10 << 20
	r.ParseMultipartForm(maxMemory)

	file, header, err := r.FormFile("thumbnail")

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing image", err)
		return
	}

	mediaType := header.Header.Get("Content-Type")

	data, err := io.ReadAll(file)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error reading image data", err)
		return

	}

	video, err := cfg.db.GetVideo(videoID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting video", err)
		return
	}

	if video.UserID.String() != userID.String() {
		respondWithError(w, http.StatusUnauthorized, "invalid request", errors.New("user does not own the video"))
		return
	}

	thumbnail := thumbnail{
		data:      data,
		mediaType: mediaType,
	}

	encodedThumbnail := base64.StdEncoding.EncodeToString(data)

	//videoThumbnails[videoID] = thumbnail

	//thumbnailURL := fmt.Sprintf("http://localhost:%v/api/thumbnails/%v", cfg.port, videoID)

	thumbnailURL := fmt.Sprintf("data:%v;base64,%v", thumbnail.mediaType, encodedThumbnail)

	video.ThumbnailURL = &thumbnailURL

	err = cfg.db.UpdateVideo(video)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating video", err)
		return
	}

	respondWithJSON(w, http.StatusOK, video)

}
