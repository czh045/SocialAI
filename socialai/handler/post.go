package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"socialai/model"
	"socialai/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var mediaTypes = map[string]string{
	".jpeg": "image",
	".jpg":  "image",
	".gif":  "image",
	".png":  "image",
	".webp": "image",
	".mov":  "video",
	".mp4":  "video",
	".avi":  "video",
	".flv":  "video",
	".wmv":  "video",
}

const maxMediaUploadSize = 20 * 1024 * 1024

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := currentUsername(r)
	if !ok {
		http.Error(w, "Missing authenticated user", http.StatusUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxMediaUploadSize)
	file, header, err := r.FormFile("media_file")
	if err != nil {
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			http.Error(w, "Media file must be 20 MB or smaller", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "Media file is not available", http.StatusBadRequest)
		return
	}
	defer file.Close()

	extension := strings.ToLower(filepath.Ext(header.Filename))
	mediaType, ok := mediaTypes[extension]
	if !ok {
		http.Error(w, "Only image and video files are supported", http.StatusBadRequest)
		return
	}

	post := model.Post{
		ID:      uuid.NewString(),
		User:    username,
		Message: r.FormValue("message"),
		Type:    mediaType,
	}

	if err := service.SavePost(&post, file, extension); err != nil {
		fmt.Printf("failed to save post: %v\n", err)
		http.Error(w, "Failed to save post to backend", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, post)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	keywords := r.URL.Query().Get("keywords")

	var (
		posts []model.Post
		err   error
	)
	if user != "" {
		posts, err = service.SearchPostsByUser(user)
	} else {
		posts, err = service.SearchPostsByKeywords(keywords)
	}
	if err != nil {
		fmt.Printf("failed to read posts: %v\n", err)
		http.Error(w, "Failed to read post from backend", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, posts)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := currentUsername(r)
	if !ok {
		http.Error(w, "Missing authenticated user", http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]
	if err := service.DeletePost(id, username); err != nil {
		fmt.Printf("failed to delete post: %v\n", err)
		http.Error(w, "Failed to delete post from backend", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
