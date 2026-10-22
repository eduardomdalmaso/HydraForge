package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"hydraforge/internal/application"
	"hydraforge/internal/domain"
	"hydraforge/internal/ports"
)

// VaultHandler exposes REST endpoints for curated datasets and active learning ingestion.
type VaultHandler struct {
	datasetService *application.VaultDatasetService
	ingestService  *application.VaultIngestService
	fileStore      ports.IFileStore
}

// NewVaultHandler creates a new VaultHandler.
func NewVaultHandler(
	datasetService *application.VaultDatasetService,
	ingestService *application.VaultIngestService,
	fileStore ports.IFileStore,
) *VaultHandler {
	return &VaultHandler{
		datasetService: datasetService,
		ingestService:  ingestService,
		fileStore:      fileStore,
	}
}

// HandleDatasets handles GET (list) and POST (create) on /api/v1/vault/datasets.
func (h *VaultHandler) HandleDatasets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		datasets, err := h.datasetService.ListDatasets(r.Context())
		if err != nil {
			writeVaultError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if datasets == nil {
			datasets = []*domain.VaultDataset{}
		}
		writeVaultJSON(w, http.StatusOK, datasets)

	case http.MethodPost:
		var d domain.VaultDataset
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			writeVaultError(w, http.StatusBadRequest, "invalid json payload")
			return
		}
		if err := h.datasetService.CreateDataset(r.Context(), &d); err != nil {
			writeVaultError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeVaultJSON(w, http.StatusCreated, d)

	default:
		writeVaultError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// HandleDatasetItem handles GET, DELETE and Export on /api/v1/vault/datasets/{id}.
func (h *VaultHandler) HandleDatasetItem(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/vault/datasets/")
	if !strings.HasPrefix(r.URL.Path, "/api/v1/vault/datasets/") {
		path = strings.TrimPrefix(r.URL.Path, "/api/v1/datasets/")
	}
	parts := strings.Split(path, "/")
	datasetID := parts[0]

	if datasetID == "" {
		writeVaultError(w, http.StatusBadRequest, "dataset_id required")
		return
	}

	if len(parts) >= 2 && parts[1] == "export" && r.Method == http.MethodPost {
		res, err := h.datasetService.ExportDataset(r.Context(), datasetID, 0.2)
		if err != nil {
			writeVaultError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeVaultJSON(w, http.StatusOK, res)
		return
	}

	if len(parts) >= 2 && parts[1] == "frames" && r.Method == http.MethodGet {
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		if limit <= 0 || limit > 100 {
			limit = 50
		}
		frames, total, err := h.datasetService.ListFrames(r.Context(), ports.FrameFilter{
			DatasetID: datasetID,
			Limit:     limit,
			Offset:    offset,
		})
		if err != nil {
			writeVaultError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeVaultJSON(w, http.StatusOK, map[string]any{
			"dataset_id": datasetID,
			"total":      total,
			"limit":      limit,
			"offset":     offset,
			"frames":     frames,
		})
		return
	}

	if r.Method == http.MethodGet {
		d, err := h.datasetService.GetDataset(r.Context(), datasetID)
		if errors.Is(err, domain.ErrDatasetNotFound) {
			writeVaultError(w, http.StatusNotFound, "dataset not found")
			return
		}
		if err != nil {
			writeVaultError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeVaultJSON(w, http.StatusOK, d)
		return
	}

	writeVaultError(w, http.StatusMethodNotAllowed, "method not allowed")
}

func writeVaultJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeVaultError(w http.ResponseWriter, status int, detail string) {
	writeVaultJSON(w, status, map[string]any{
		"error":  strings.TrimSpace(detail),
		"status": status,
	})
}
