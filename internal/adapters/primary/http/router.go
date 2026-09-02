package http

import "net/http"

// RegisterRoutes registers all API endpoints on the mux.
func RegisterRoutes(mux *http.ServeMux, h *TrainingHandler, bh *BenchmarkHandler) {
	mux.HandleFunc("/api/v1/training/jobs", h.HandleJobs)
	mux.HandleFunc("/api/v1/training/jobs/", h.HandleJobByID)
	mux.HandleFunc("/api/v1/training/datasets/import", h.HandleDatasetImport)
	mux.HandleFunc("/api/v1/training/datasets/merge", h.HandleDatasetMerge)
	mux.HandleFunc("/api/v1/training/datasets/audit/", h.HandleDatasetAudit)
	mux.HandleFunc("/api/v1/training/datasets/image", h.HandleDatasetImage)
	mux.HandleFunc("/api/v1/training/datasets/sample", h.HandleDatasetSample)
	mux.HandleFunc("/api/v1/training/datasets/rescan", h.HandleDatasetRescan)
	mux.HandleFunc("/api/v1/training/datasets/register-path", h.HandleDatasetRegisterPath)
	mux.HandleFunc("/api/v1/training/datasets", h.HandleDatasets)
	mux.HandleFunc("/api/v1/training/datasets/", h.HandleDatasetByID)
	mux.HandleFunc("/api/v1/training/telemetry", h.HandleTelemetry)
	mux.HandleFunc("/api/v1/training/export", h.HandleExport)
	mux.HandleFunc("/api/v1/training/models", h.HandleModels)

	// Benchmarks API
	mux.HandleFunc("/api/v1/benchmarks", bh.HandleBenchmarks)
	mux.HandleFunc("/api/v1/benchmarks/formats", bh.HandleBenchmarkFormats)
	mux.HandleFunc("/api/v1/benchmarks/", bh.HandleBenchmarkByID)


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
