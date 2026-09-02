package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"hydraforge/internal/domain"
	"hydraforge/internal/ports"
)

// BenchmarkHandler handles HTTP requests for Ultralytics YOLO model benchmarking.
type BenchmarkHandler struct {
	useCase ports.TrainingUseCase
}

// NewBenchmarkHandler creates a new BenchmarkHandler instance.
func NewBenchmarkHandler(useCase ports.TrainingUseCase) *BenchmarkHandler {
	return &BenchmarkHandler{useCase: useCase}
}

// HandleBenchmarks handles GET /api/v1/benchmarks and POST /api/v1/benchmarks.
func (h *BenchmarkHandler) HandleBenchmarks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		status := r.URL.Query().Get("status")
		benchmarks, err := h.useCase.ListBenchmarkJobs(r.Context(), status)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(benchmarks)

	case http.MethodPost:
		var job domain.BenchmarkJob
		if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
			http.Error(w, `{"error":"invalid JSON request payload"}`, http.StatusBadRequest)
			return
		}
		created, err := h.useCase.CreateBenchmarkJob(r.Context(), &job)
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

// HandleBenchmarkByID handles GET and DELETE /api/v1/benchmarks/{id}.
func (h *BenchmarkHandler) HandleBenchmarkByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/benchmarks/")
	if id == "" {
		http.Error(w, `{"error":"benchmark ID required"}`, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		job, err := h.useCase.GetBenchmarkJob(r.Context(), id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(job)

	case http.MethodDelete:
		if err := h.useCase.StopBenchmarkJob(r.Context(), id); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"stopped"}`))

	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// HandleBenchmarkFormats handles GET /api/v1/benchmarks/formats.
func (h *BenchmarkHandler) HandleBenchmarkFormats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	formats, err := h.useCase.GetSupportedBenchmarkFormats(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(formats)
}
