package http

import (
	"encoding/json"
	"net/http"

	"hydraforge/internal/adapters/primary/http/middleware"
)

// RegisterRoutes registers all API endpoints on the mux.
func RegisterRoutes(mux *http.ServeMux, h *TrainingHandler, bh *BenchmarkHandler, mh *MediaHandler) {
	// Helper to protect endpoints with Auth and optional RBAC middlewares
	authEndpoint := func(pattern string, fn http.HandlerFunc, roles ...string) {
		var handler http.Handler = http.HandlerFunc(fn)
		if len(roles) > 0 {
			handler = middleware.RequireRole(roles...)(handler)
		}
		handler = middleware.AuthMiddleware(handler)
		mux.Handle(pattern, handler)
	}

	// 1. Public Endpoints (Info, Swagger, Health Probes)
	mux.HandleFunc("/api/v1/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"app_name":     "HydraForge AI Training & Inference",
			"version":      "1.0.0",
			"engine_mode":  "TensorRT / PyTorch CUDA 13.3",
			"gpu_detected": true,
			"gpu_model":    "NVIDIA GeForce RTX 5090 (32GB VRAM)",
			"features":     map[string]bool{"tensorrt": true, "sahi": true, "amp_fp16": true},
		})
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("READY"))
	})
	mux.HandleFunc("/swagger/", ServeSwaggerUI)
	mux.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(OpenAPI3Spec))
	})

	// 2. Training Jobs API (Auth Required; Create/Delete requires admin/operator)
	authEndpoint("/api/v1/training/jobs", h.HandleJobs)
	authEndpoint("/api/v1/training/jobs/", h.HandleJobByID)
	authEndpoint("/api/v1/training/telemetry", h.HandleTelemetry)
	authEndpoint("/api/v1/training/export", h.HandleExport, "admin", "superadmin")
	authEndpoint("/api/v1/training/models", h.HandleModels)

	// 3. Datasets API
	authEndpoint("/api/v1/training/datasets", h.HandleDatasets)
	authEndpoint("/api/v1/training/datasets/", h.HandleDatasetByID)
	authEndpoint("/api/v1/training/datasets/import", h.HandleDatasetImport, "admin", "operator", "superadmin")
	authEndpoint("/api/v1/training/datasets/merge", h.HandleDatasetMerge, "admin", "operator", "superadmin")
	authEndpoint("/api/v1/training/datasets/audit/", h.HandleDatasetAudit)
	authEndpoint("/api/v1/training/datasets/image", h.HandleDatasetImage)
	authEndpoint("/api/v1/training/datasets/sample", h.HandleDatasetSample)
	authEndpoint("/api/v1/training/datasets/rescan", h.HandleDatasetRescan, "admin", "operator", "superadmin")
	authEndpoint("/api/v1/training/datasets/register-path", h.HandleDatasetRegisterPath, "admin", "superadmin")

	// 4. Benchmarks API
	authEndpoint("/api/v1/benchmarks", bh.HandleBenchmarks)
	authEndpoint("/api/v1/benchmarks/formats", bh.HandleBenchmarkFormats)
	authEndpoint("/api/v1/benchmarks/", bh.HandleBenchmarkByID)

	// 5. Media Sources API (Video Loops & Folders)
	if mh != nil {
		authEndpoint("/api/v1/media/sources", mh.HandleSources)
		authEndpoint("/api/v1/media/folders", mh.HandleFolders, "admin", "operator", "superadmin")
		authEndpoint("/api/v1/media/folders/", mh.HandleFolders, "admin", "operator", "superadmin")
		authEndpoint("/api/v1/media/upload", mh.HandleUpload, "admin", "operator", "superadmin")
		authEndpoint("/api/v1/media/files/", mh.HandleFiles, "admin", "operator", "superadmin")
		mux.HandleFunc("/api/v1/media/stream/", mh.HandleStream)
	}

	// 6. HydraStream Integration Proxy (Port 8080)
	authEndpoint("/api/v1/hydrastream/", HydraStreamProxy("http://localhost:8080"))

	// 7. Real-Time Optical YOLO Inference (NVIDIA RTX 5090)
	authEndpoint("/api/v1/inference/predict", HandleInferencePredict)
	authEndpoint("/api/v1/inference/live", HandleInferenceLiveStream)
	authEndpoint("/api/v1/inference/frame", HandleWebcamFrameUpload)
}

// WithCORS wraps an http.Handler with universal Cross-Origin Resource Sharing headers.
func WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key, X-Requested-With, Range, Accept, Origin")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range, Content-Type, X-Total-Count")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
