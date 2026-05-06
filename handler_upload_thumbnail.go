package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

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
		respondWithError(w, http.StatusBadRequest, "Couldn't get image", err)
		return
	}

	mimeType := header.Header.Get("Content-Type")

	mediaType, err := validateImageMediaType(mimeType)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid image format", err)
		return
	}

	video, err := cfg.db.GetVideo(videoID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting video", err)
		return
	}

	if video.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized request", err)
		return
	}

	encoded := getRandomPath()

	imageFilePath := getAssetPath(encoded, mediaType)

	assetFilepath := cfg.getAssetFilepath(imageFilePath)
	fileDest, err := os.Create(assetFilepath)

	if err != nil {
		fmt.Print(err)
		respondWithError(w, http.StatusInternalServerError, "Error creating image file", err)
		return
	}

	_, err = io.Copy(fileDest, file)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error saving image to FS", err)
		return
	}

	/* imgBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	imageUrl := fmt.Sprintf("data:%s;base64,%s", mediaType, imgBase64)

	*/
	thumbnailURL := cfg.getAssetURL(imageFilePath)

	video.ThumbnailURL = &thumbnailURL

	err = cfg.db.UpdateVideo(video)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt update video thumbnail", err)
		return
	}

	respondWithJSON(w, http.StatusOK, video)
}
