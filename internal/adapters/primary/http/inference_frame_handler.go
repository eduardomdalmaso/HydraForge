package http

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"hydraforge/internal/config"
)

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
