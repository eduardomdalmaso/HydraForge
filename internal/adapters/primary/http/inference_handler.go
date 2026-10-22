package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"hydraforge/internal/config"
)

// InferenceRequest represents payload for YOLO model inference.
type InferenceRequest struct {
	Model       string  `json:"model"`
	Source      string  `json:"source"`
	ImageBase64 string  `json:"image_base64,omitempty"`
	Conf        float64 `json:"conf"`
	IoU         float64 `json:"iou"`
	Device      string  `json:"device"`
	Classes     string  `json:"classes,omitempty"`
	Sahi        bool    `json:"sahi,omitempty"`
	NmsFree     bool    `json:"nms_free,omitempty"`
	Track       bool    `json:"track,omitempty"`
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
		Model:   "yolo26n",
		Source:  "cam_entrance_01",
		Conf:    0.25,
		IoU:     0.45,
		Device:  "0",
		Track:   true,
		Sahi:    false,
		NmsFree: false,
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
		if q.Get("sahi") == "true" {
			req.Sahi = true
		}
		if q.Get("nms_free") == "true" {
			req.NmsFree = true
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

	cleanDevice := strings.TrimSpace(req.Device)
	if cleanDevice != "cpu" {
		if _, err := strconv.Atoi(cleanDevice); err != nil {
			cleanDevice = "0"
		}
	}
	req.Device = cleanDevice

	cleanModel := filepath.Base(req.Model)
	if cleanModel == "" || strings.Contains(cleanModel, "..") {
		cleanModel = "yolo26n.pt"
	}
	req.Model = cleanModel

	modelFile := resolveModelWeights(req.Model)
	pythonBin := config.GetPythonBin()

	scriptPath := "worker_python/predict.py"
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	args := []string{
		scriptPath,
		"--weights", modelFile,
		"--source", req.Source,
		"--conf", fmt.Sprintf("%.2f", req.Conf),
		"--iou", fmt.Sprintf("%.2f", req.IoU),
		"--device", req.Device,
	}
	if req.Classes != "" {
		args = append(args, "--classes", req.Classes)
	}
	if req.Track && !req.Sahi {
		args = append(args, "--track")
	}
	if req.Sahi {
		args = append(args, "--sahi")
	}
	if req.NmsFree {
		args = append(args, "--nms-free")
	}

	log.Printf("[HydraForge Inference] Executing: %s %v", pythonBin, args)
	cmd := exec.CommandContext(ctx, pythonBin, args...)
	cmd.Env = append(os.Environ(), "PYTHONUNBUFFERED=1")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  fmt.Sprintf("inference failed: %v | stdout: %s | stderr: %s", err, stdout.String(), stderr.String()),
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
		filepath.Join("storage/models", model),
		filepath.Join("storage/models", cleanID+".pt"),
		filepath.Join("storage/engines", cleanID+".engine"),
		filepath.Join("runs/train", cleanID, "weights/best.pt"),
		filepath.Join("weights", model),
		filepath.Join("weights", cleanID+".pt"),
		filepath.Join("weights", cleanID+".engine"),
	}
	for _, base := range config.GetRunsSearchDirs() {
		candidates = append(candidates, filepath.Join(base, cleanID, "weights/best.pt"))
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return filepath.Join("weights", model)
}
