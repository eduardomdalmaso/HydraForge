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

	// Sub-resource: .../export
	if len(parts) >= 2 && parts[1] == "export" && r.Method == http.MethodPost {
		res, err := h.datasetService.ExportDataset(r.Context(), datasetID, 0.2)
		if err != nil {
			writeVaultError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeVaultJSON(w, http.StatusOK, res)
		return
	}

	// Sub-resource: .../frames
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

// HandleUpload receives multipart image uploads and executes active learning ingest.
func (h *VaultHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeVaultError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Max 32MB upload
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeVaultError(w, http.StatusBadRequest, "failed to parse multipart form: "+err.Error())
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		writeVaultError(w, http.StatusBadRequest, "missing 'image' form file")
		return
	}
	defer file.Close()

	datasetID := strings.TrimSpace(r.FormValue("dataset_id"))
	if datasetID == "" {
		writeVaultError(w, http.StatusBadRequest, "missing 'dataset_id'")
		return
	}

	cameraID := strings.TrimSpace(r.FormValue("camera_id"))
	if cameraID == "" {
		cameraID = "cam_manual_upload"
	}

	var bboxes []domain.BBox
	bboxesJSON := r.FormValue("bboxes")
	if bboxesJSON != "" {
		var rawList []map[string]any
		if err := json.Unmarshal([]byte(bboxesJSON), &rawList); err == nil {
			for _, item := range rawList {
				bbox := parseFlexibleVaultBBox(item)
				if err := bbox.Validate(); err == nil {
					bboxes = append(bboxes, bbox)
				}
			}
		}
	}

	frame, err := h.ingestService.IngestFrame(r.Context(), application.IngestCommand{
		DatasetID:   datasetID,
		CameraID:    cameraID,
		FileName:    header.Filename,
		ImageReader: file,
		BBoxes:      bboxes,
		SourceEvent: r.FormValue("source_event"),
	})
	if errors.Is(err, domain.ErrDuplicateFrame) {
		writeVaultJSON(w, http.StatusConflict, map[string]any{
			"error":    err.Error(),
			"frame_id": frame.FrameID,
			"status":   "duplicate_skipped",
		})
		return
	}
	if err != nil {
		writeVaultError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeVaultJSON(w, http.StatusCreated, frame)
}

func parseFlexibleVaultBBox(m map[string]any) domain.BBox {
	var bbox domain.BBox

	if id, ok := m["class_id"].(float64); ok {
		bbox.ClassID = int(id)
	}
	if name, ok := m["class_name"].(string); ok {
		bbox.ClassName = strings.ToLower(strings.TrimSpace(name))
	} else if lbl, ok := m["label"].(string); ok {
		bbox.ClassName = strings.ToLower(strings.TrimSpace(lbl))
	}

	if conf, ok := m["confidence"].(float64); ok {
		bbox.Confidence = conf
	}

	if xc, ok := m["x_center"].(float64); ok {
		bbox.XCenter = xc
		if yc, ok := m["y_center"].(float64); ok {
			bbox.YCenter = yc
		}
		if w, ok := m["width"].(float64); ok {
			bbox.Width = w
		}
		if h, ok := m["height"].(float64); ok {
			bbox.Height = h
		}
	} else if rawBox, ok := m["box"].([]any); ok && len(rawBox) >= 4 {
		x, _ := rawBox[0].(float64)
		y, _ := rawBox[1].(float64)
		w, _ := rawBox[2].(float64)
		h, _ := rawBox[3].(float64)

		if x > 1.0 || y > 1.0 || w > 1.0 || h > 1.0 {
			x /= 100.0
			y /= 100.0
			w /= 100.0
			h /= 100.0
		}
		bbox.XCenter = x + (w / 2.0)
		bbox.YCenter = y + (h / 2.0)
		bbox.Width = w
		bbox.Height = h
	}

	if bbox.Width <= 0.0 {
		bbox.Width = 0.1
	}
	if bbox.Height <= 0.0 {
		bbox.Height = 0.1
	}
	if bbox.XCenter <= 0.0 {
		bbox.XCenter = bbox.Width / 2.0
	}
	if bbox.YCenter <= 0.0 {
		bbox.YCenter = bbox.Height / 2.0
	}

	return bbox
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
