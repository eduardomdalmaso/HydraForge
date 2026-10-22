package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"hydraforge/internal/config"
	"hydraforge/internal/domain"
)

// HandleDatasets handles GET and POST /api/v1/training/datasets.
func (h *TrainingHandler) HandleDatasets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		datasets, err := h.useCase.ListDatasets(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(datasets)

	case http.MethodPost:
		var d domain.Dataset
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			http.Error(w, `{"error":"invalid JSON dataset payload"}`, http.StatusBadRequest)
			return
		}
		saved, err := h.useCase.RegisterDataset(r.Context(), &d)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(saved)

	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// HandleDatasetByID handles GET and DELETE /api/v1/training/datasets/{id}.
func (h *TrainingHandler) HandleDatasetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/training/datasets/")
	if id == "rescan" {
		h.HandleDatasetRescan(w, r)
		return
	}
	if id == "register-path" {
		h.HandleDatasetRegisterPath(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		d, err := h.useCase.GetDataset(r.Context(), id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(d)

	case http.MethodDelete:
		deleteFiles := r.URL.Query().Get("delete_files") == "true"
		if err := h.useCase.DeleteDataset(r.Context(), id, deleteFiles); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     "deleted",
			"dataset_id": id,
		})

	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// HandleDatasetRescan handles POST /api/v1/training/datasets/rescan.
func (h *TrainingHandler) HandleDatasetRescan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	datasets, err := h.useCase.RescanDatasets(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to rescan datasets: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(datasets)
}

// HandleDatasetRegisterPath handles POST /api/v1/training/datasets/register-path.
func (h *TrainingHandler) HandleDatasetRegisterPath(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		Path      string `json:"path"`
		DatasetID string `json:"dataset_id"`
		Name      string `json:"name"`
		Task      string `json:"task"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid JSON body"}`, http.StatusBadRequest)
		return
	}
	if req.Path == "" {
		http.Error(w, `{"error":"directory path is required"}`, http.StatusBadRequest)
		return
	}

	cleanReqPath := filepath.Clean(req.Path)
	normalizedPath := filepath.ToSlash(cleanReqPath)
	disallowedPrefixes := []string{"/etc", "/root", "/proc", "/sys", "/dev", "/boot", "/var/run", "/usr", "/bin", "/sbin", "C:/Windows", "C:/Program Files"}
	for _, p := range disallowedPrefixes {
		if strings.HasPrefix(normalizedPath, p) || strings.HasPrefix(cleanReqPath, p) || strings.EqualFold(normalizedPath, p) {
			http.Error(w, `{"error":"forbidden system directory path"}`, http.StatusForbidden)
			return
		}
	}

	fi, err := os.Stat(cleanReqPath)
	if err != nil || !fi.IsDir() {
		http.Error(w, fmt.Sprintf(`{"error":"directory not found: %s"}`, cleanReqPath), http.StatusBadRequest)
		return
	}

	id := req.DatasetID
	if id == "" {
		id = filepath.Base(cleanReqPath)
	}
	id = strings.ToLower(strings.ReplaceAll(id, " ", "_"))

	targetLink := filepath.Join(config.GetDatasetsDir(), id)
	if _, err := os.Stat(targetLink); err != nil && targetLink != cleanReqPath {
		_ = os.Symlink(cleanReqPath, targetLink)
	}

	datasets, err := h.useCase.RescanDatasets(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	var found *domain.Dataset
	for _, d := range datasets {
		if d.DatasetID == id {
			found = d
			break
		}
	}
	if found == nil && len(datasets) > 0 {
		found = datasets[0]
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(found)
}
