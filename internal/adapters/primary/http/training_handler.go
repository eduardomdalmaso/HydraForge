package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"hydraforge/internal/domain"
	"hydraforge/internal/ports"
)

// TrainingHandler handles HTTP requests for training management.
type TrainingHandler struct {
	useCase ports.TrainingUseCase
}

// NewTrainingHandler creates a new handler instance.
func NewTrainingHandler(useCase ports.TrainingUseCase) *TrainingHandler {
	return &TrainingHandler{useCase: useCase}
}

// HandleJobs handles GET /api/v1/training/jobs and POST /api/v1/training/jobs.
func (h *TrainingHandler) HandleJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		status := r.URL.Query().Get("status")
		jobs, err := h.useCase.ListTrainingJobs(r.Context(), status)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(jobs)

	case http.MethodPost:
		var job domain.TrainingJob
		if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
			http.Error(w, `{"error":"invalid JSON request payload"}`, http.StatusBadRequest)
			return
		}
		created, err := h.useCase.CreateTrainingJob(r.Context(), &job)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(created)

	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// HandleJobByID handles GET and DELETE /api/v1/training/jobs/{id}.
func (h *TrainingHandler) HandleJobByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/training/jobs/")
	if id == "" {
		http.Error(w, `{"error":"job ID required"}`, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		job, err := h.useCase.GetTrainingJob(r.Context(), id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(job)

	case http.MethodDelete:
		if err := h.useCase.StopTrainingJob(r.Context(), id); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"stopped"}`))

	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

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

// HandleDatasetByID handles GET /api/v1/training/datasets/{id}.
func (h *TrainingHandler) HandleDatasetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/training/datasets/")
	d, err := h.useCase.GetDataset(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(d)
}

// HandleTelemetry handles GET /api/v1/training/telemetry.
func (h *TrainingHandler) HandleTelemetry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	telemetry, err := h.useCase.GetCockpitTelemetry(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(telemetry)
}

// HandleExport handles POST /api/v1/training/export.
func (h *TrainingHandler) HandleExport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"POST required"}`, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CheckpointID string              `json:"checkpoint_id"`
		Format       domain.ExportFormat `json:"format"`
		Precision    string              `json:"precision"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid JSON request"}`, http.StatusBadRequest)
		return
	}

	exported, err := h.useCase.ExportModel(r.Context(), req.CheckpointID, req.Format, req.Precision)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(exported)
}
