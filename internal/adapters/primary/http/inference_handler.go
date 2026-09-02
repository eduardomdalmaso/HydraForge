package http

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// InferenceRequest represents payload for YOLO model inference.
type InferenceRequest struct {
	Model       string  `json:"model"`
	Source      string  `json:"source"`
	ImageBase64 string  `json:"image_base64,omitempty"`
	Conf        float64 `json:"conf"`
	IoU         float64 `json:"iou"`
	Device      string  `json:"device"`
}

// HandleInferencePredict executes real YOLO model inference on NVIDIA GPU via Python Worker.
func HandleInferencePredict(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	req := InferenceRequest{
		Model:  "yolo26n",
		Source: "/home/hades/datasets/frota_urbana_4classes/train/images/rf100_adit_mp4-100_jpg.rf.1ceee906e811b590f48e8e4decda7380.jpg",
		Conf:   0.25,
		IoU:    0.45,
		Device: "0",
	}

	if r.Method == http.MethodPost {
		_ = json.NewDecoder(r.Body).Decode(&req)
	} else {
		q := r.URL.Query()
		if m := q.Get("model"); m != "" {
			req.Model = m
		}
		if s := q.Get("source"); s != "" {
			req.Source = s
		}
		if c, err := strconv.ParseFloat(q.Get("conf"), 64); err == nil && c > 0 {
			req.Conf = c
		}
		if i, err := strconv.ParseFloat(q.Get("iou"), 64); err == nil && i > 0 {
			req.IoU = i
		}
	}

	if req.ImageBase64 != "" {
		b64Data := req.ImageBase64
		if idx := strings.Index(b64Data, ","); idx != -1 {
			b64Data = b64Data[idx+1:]
		}
		if imgBytes, err := base64.StdEncoding.DecodeString(b64Data); err == nil && len(imgBytes) > 0 {
			shmPath := "/dev/shm/hydra_webcam_frame.jpg"
			if err := os.WriteFile(shmPath, imgBytes, 0644); err == nil {
				req.Source = shmPath
			}
		}
	}

	// Resolve model weight path
	modelFile := resolveModelWeights(req.Model)

	pythonBin := "/home/hades/miniconda3/envs/analytics-env/bin/python"
	if _, err := os.Stat(pythonBin); err != nil {
		pythonBin = "python3"
	}

	scriptPath := "worker_python/predict.py"
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, pythonBin, scriptPath,
		"--weights", modelFile,
		"--source", req.Source,
		"--conf", fmt.Sprintf("%.2f", req.Conf),
		"--iou", fmt.Sprintf("%.2f", req.IoU),
		"--device", req.Device,
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  fmt.Sprintf("inference failed: %v | stderr: %s", err, stderr.String()),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(stdout.Bytes())
}

func resolveModelWeights(model string) string {
	if filepath.IsAbs(model) {
		return model
	}
	cleanID := strings.TrimSuffix(model, ".pt")
	candidates := []string{
		filepath.Join("runs/train", cleanID, "weights/best.pt"),
		filepath.Join("/home/hades/Documents/HydraForge/runs/train", cleanID, "weights/best.pt"),
		filepath.Join("/home/hades/runs/train", cleanID, "weights/best.pt"),
		filepath.Join("weights", model),
		filepath.Join("weights", cleanID+".pt"),
		filepath.Join("weights", cleanID+".engine"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return filepath.Join("weights", model)
}

// HandleInferenceLiveStream provides continuous Server-Sent Events (SSE) of real-time detections.
func HandleInferenceLiveStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	q := r.URL.Query()
	model := q.Get("model")
	if model == "" {
		model = "yolo26n"
	}
	source := q.Get("source")
	if source == "" {
		source = "cam_entrance_01"
	}
	conf := q.Get("conf")
	if conf == "" {
		conf = "0.25"
	}

	modelFile := resolveModelWeights(model)
	pythonBin := "/home/hades/miniconda3/envs/analytics-env/bin/python"
	if _, err := os.Stat(pythonBin); err != nil {
		pythonBin = "python3"
	}

	sourcePath := source
	if source == "webcam" {
		sourcePath = "/dev/shm/hydra_webcam_frame.jpg"
	} else if !filepath.IsAbs(source) {
		shmSample := filepath.Join("/home/hades/Documents/HydraStream/samples", fmt.Sprintf("%s.jpg", source))
		if _, err := os.Stat(shmSample); err == nil {
			sourcePath = shmSample
		}
	}

	cmd := exec.CommandContext(r.Context(), pythonBin, "worker_python/live_engine.py",
		"--weights", modelFile,
		"--source", sourcePath,
		"--conf", conf,
		"--device", "0",
	)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, "Failed to start analytics pipe", http.StatusInternalServerError)
		return
	}

	if err := cmd.Start(); err != nil {
		http.Error(w, "Failed to start analytics worker", http.StatusInternalServerError)
		return
	}
	defer func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	}()

	scanner := bufio.NewScanner(stdoutPipe)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && strings.HasPrefix(line, "{") {
			fmt.Fprintf(w, "data: %s\n\n", line)
			flusher.Flush()
		}
	}
}

// HandleWebcamFrameUpload receives raw webcam frames from the browser and writes to /dev/shm.
func HandleWebcamFrameUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req struct {
		ImageBase64 string `json:"image_base64"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err == nil && req.ImageBase64 != "" {
		b64Data := req.ImageBase64
		if idx := strings.Index(b64Data, ","); idx != -1 {
			b64Data = b64Data[idx+1:]
		}
		if imgBytes, err := base64.StdEncoding.DecodeString(b64Data); err == nil && len(imgBytes) > 0 {
			_ = os.WriteFile("/dev/shm/hydra_webcam_frame.jpg", imgBytes, 0644)
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
