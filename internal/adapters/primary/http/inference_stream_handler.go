package http

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"hydraforge/internal/config"
)

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

