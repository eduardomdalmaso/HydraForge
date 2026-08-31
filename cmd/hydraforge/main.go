package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	primaryHttp "hydraforge/internal/adapters/primary/http"
	"hydraforge/internal/adapters/secondary/gpu"
	"hydraforge/internal/adapters/secondary/memory"
	"hydraforge/internal/adapters/secondary/worker"
	"hydraforge/internal/application"
)

func main() {
	log.Println("🚀 [HydraForge] Initializing YOLO AI Training Studio (Hexagonal Architecture + DDD)...")

	// 1. Initialize Secondary Adapters (Driven)
	memStore := memory.NewMemoryStore()
	gpuDetector := gpu.NewDetector()
	pyWorker := worker.NewPythonWorker()

	// 2. Query hardware status
	gpuTelemetry, err := gpuDetector.QueryGPU(context.Background())
	if err == nil && gpuTelemetry != nil {
		log.Printf("⚡ [HydraForge] Active Hardware: %s (%.0f MB VRAM) | Temp: %.1f°C", gpuTelemetry.Model, gpuTelemetry.TotalVRAMMB, gpuTelemetry.TempCelsius)
	}

	// 3. Initialize Application Service (Use Case)
	trainingService := application.NewTrainingService(memStore, memStore, memStore, pyWorker, gpuDetector)

	// 4. Initialize Primary HTTP Adapter (Driving)
	handler := primaryHttp.NewTrainingHandler(trainingService)
	mux := http.NewServeMux()
	primaryHttp.RegisterRoutes(mux, handler)

	// Static SPA Web UI
	mux.Handle("/", http.FileServer(http.Dir("./web/dist")))

	server := &http.Server{
		Addr:         ":8081",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🌐 [HydraForge] Training Studio & REST Control Plane listening on http://localhost:8081")
		log.Printf("📖 [HydraForge] Swagger API Docs available at http://localhost:8081/swagger/")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 [HydraForge] Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
	log.Println("✅ [HydraForge] Stopped.")
}
