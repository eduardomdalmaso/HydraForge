package http

import (
	"bufio"
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

	// Validate device and model to prevent command/argument injection
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

	// Resolve model weight path
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
	iou := q.Get("iou")
	if iou == "" {
		iou = "0.45"
	}

	modelFile := resolveModelWeights(model)
	pythonBin := config.GetPythonBin()

	sourcePath := source
	if source == "webcam" {
		sourcePath = filepath.Join(os.TempDir(), "hydra_webcam_frame.jpg")
	} else if strings.HasPrefix(source, "rtsp://") || strings.HasPrefix(source, "rtmp://") || strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") || strings.HasPrefix(source, "cam") {
		sourcePath = source
	} else if !filepath.IsAbs(source) {
		shmSample := filepath.Join(config.GetHydraStreamSamplesDir(), fmt.Sprintf("%s.jpg", source))
		if _, err := os.Stat(shmSample); err == nil {
			sourcePath = shmSample
		}
	}

	cmdArgs := []string{
		"worker_python/live_engine.py",
		"--weights", modelFile,
		"--source", sourcePath,
		"--conf", conf,
		"--iou", iou,
		"--device", "0",
	}
	if classes := q.Get("classes"); classes != "" {
		cmdArgs = append(cmdArgs, "--classes", classes)
	}
	if q.Get("sahi") == "true" {
		cmdArgs = append(cmdArgs, "--sahi")
	}
	if q.Get("nms_free") == "true" {
		cmdArgs = append(cmdArgs, "--nms-free")
	}

	cmd := exec.CommandContext(r.Context(), pythonBin, cmdArgs...)
	cmd.Env = append(os.Environ(), "PYTHONUNBUFFERED=1")

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
	if err := scanner.Err(); err != nil {
		log.Printf("[HydraForge] Scanner error in analytics stream: %v", err)
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

// HandleGetLiveDetections returns real-time detections from /dev/shm.
func HandleGetLiveDetections(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	camID := r.URL.Query().Get("camera_id")
	if camID == "" {
		camID = "cam_01"
	}

	baseID := filepath.Base(camID)
	candidates := []string{
		fmt.Sprintf("/dev/shm/face_detections_%s.json", baseID),
		fmt.Sprintf("/dev/shm/detections_%s.json", baseID),
		filepath.Join(os.TempDir(), fmt.Sprintf("face_detections_%s.json", baseID)),
		filepath.Join(os.TempDir(), fmt.Sprintf("detections_%s.json", baseID)),
	}

	for _, p := range candidates {
		if data, err := os.ReadFile(p); err == nil && len(data) > 0 {
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"waiting","detections":[],"count":0}`))
}

// HandleGetInferenceFrame serves the latest captured JPEG frame for a given camera/source.
func HandleGetInferenceFrame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Content-Type", "image/jpeg")

	src := r.URL.Query().Get("source")
	if src == "" {
		src = "latest"
	}
	cleanID := filepath.Base(src)
	cleanID = strings.TrimSuffix(cleanID, ".jpg")

	candidates := []string{
		filepath.Join("storage/captures", cleanID+".jpg"),
		filepath.Join("storage/captures", "live_"+cleanID+".jpg"),
		filepath.Join(os.TempDir(), fmt.Sprintf("live_infer_%s.jpg", cleanID)),
		filepath.Join("storage/captures", "latest.jpg"),
		filepath.Join(os.TempDir(), "rtsp_capture_frame.jpg"),
		filepath.Join(config.GetHydraStreamSamplesDir(), fmt.Sprintf("%s.jpg", cleanID)),
	}

	for _, p := range candidates {
		if data, err := os.ReadFile(p); err == nil && len(data) > 0 {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(data)
			return
		}
	}

	http.Error(w, "frame not available", http.StatusNotFound)
}
