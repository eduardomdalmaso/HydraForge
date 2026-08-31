package http

import "net/http"

// RegisterRoutes registers all API endpoints on the mux.
func RegisterRoutes(mux *http.ServeMux, h *TrainingHandler) {
	mux.HandleFunc("/api/v1/training/jobs", h.HandleJobs)
	mux.HandleFunc("/api/v1/training/jobs/", h.HandleJobByID)
	mux.HandleFunc("/api/v1/training/datasets", h.HandleDatasets)
	mux.HandleFunc("/api/v1/training/datasets/", h.HandleDatasetByID)
	mux.HandleFunc("/api/v1/training/telemetry", h.HandleTelemetry)
	mux.HandleFunc("/api/v1/training/export", h.HandleExport)

	// Swagger Interactive Docs
	mux.HandleFunc("/swagger/", ServeSwaggerUI)
	mux.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(OpenAPI3Spec))
	})

	// Health probes
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("READY"))
	})
}
