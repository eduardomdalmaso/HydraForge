package http

import (
	"encoding/json"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"hydraforge/internal/application"
)

// MediaHandler exposes REST and streaming endpoints for video loop media assets.
type MediaHandler struct {
	mediaService *application.MediaService
}

// NewMediaHandler creates a new MediaHandler.
func NewMediaHandler(mediaService *application.MediaService) *MediaHandler {
	return &MediaHandler{mediaService: mediaService}
}

func writeJSONError(w http.ResponseWriter, status int, msg string, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
		"code":  code,
	})
}

// HandleSources routes GET /api/v1/media/sources.
func (h *MediaHandler) HandleSources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed", "METHOD_NOT_ALLOWED")
		return
	}

	resp, err := h.mediaService.ListSources(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error(), "MEDIA_SCAN_ERROR")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// HandleFolders routes POST /api/v1/media/folders and DELETE /api/v1/media/folders/{name}.
func (h *MediaHandler) HandleFolders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body", "INVALID_BODY")
			return
		}

		if err := h.mediaService.CreateFolder(r.Context(), req.Name); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error(), "FOLDER_CREATE_FAILED")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{"status": "created", "folder": req.Name})

	case http.MethodDelete:
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/media/folders/")
		folderName := strings.TrimSpace(path)
		if folderName == "" {
			writeJSONError(w, http.StatusBadRequest, "folder name required in path", "INVALID_FOLDER")
			return
		}

		if err := h.mediaService.DeleteFolder(r.Context(), folderName); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error(), "FOLDER_DELETE_FAILED")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"status": "deleted", "folder": folderName})

	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed", "METHOD_NOT_ALLOWED")
	}
}

// HandleUpload handles multipart/form-data upload to a category folder.
func (h *MediaHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed", "METHOD_NOT_ALLOWED")
		return
	}

	// Max 500MB upload
	if err := r.ParseMultipartForm(500 << 20); err != nil {
		writeJSONError(w, http.StatusBadRequest, "failed to parse multipart form: "+err.Error(), "UPLOAD_PARSE_ERROR")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "missing 'file' field", "MISSING_FILE")
		return
	}
	defer file.Close()

	folder := strings.TrimSpace(r.FormValue("folder"))
	if folder == "" {
		folder = "root"
	}

	mediaFile, err := h.mediaService.SaveFile(r.Context(), folder, header.Filename, file)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error(), "SAVE_FILE_ERROR")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mediaFile)
}

// HandleFiles handles deleting files: DELETE /api/v1/media/files/{folder}/{filename}.
func (h *MediaHandler) HandleFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed", "METHOD_NOT_ALLOWED")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/media/files/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		writeJSONError(w, http.StatusBadRequest, "path must be /api/v1/media/files/{folder}/{filename}", "INVALID_PATH")
		return
	}

	if err := h.mediaService.DeleteFile(r.Context(), parts[0], parts[1]); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error(), "DELETE_FILE_ERROR")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"status": "deleted", "folder": parts[0], "file": parts[1]})
}

// HandleStream serves media files with HTTP 206 Range Requests (Zero-Copy Linux sendfile).
func (h *MediaHandler) HandleStream(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/media/stream/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		http.Error(w, "invalid stream path", http.StatusBadRequest)
		return
	}

	folderName, filename := parts[0], parts[1]
	f, info, err := h.mediaService.OpenFile(r.Context(), folderName, filename)
	if err != nil {
		http.Error(w, "file not found: "+err.Error(), http.StatusNotFound)
		return
	}
	defer f.Close()

	ext := strings.ToLower(filepath.Ext(filename))
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		switch ext {
		case ".mp4":
			contentType = "video/mp4"
		case ".webm":
			contentType = "video/webm"
		case ".avi":
			contentType = "video/x-msvideo"
		case ".mkv":
			contentType = "video/x-matroska"
		default:
			contentType = "application/octet-stream"
		}
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	http.ServeContent(w, r, filename, info.ModTime(), f)
}
