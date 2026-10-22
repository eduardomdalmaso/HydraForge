package http

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"hydraforge/internal/config"
)

// HandleDatasetImage serves raw images from datasets directory safely.
func (h *TrainingHandler) HandleDatasetImage(w http.ResponseWriter, r *http.Request) {
	relPath := r.URL.Query().Get("path")
	if relPath == "" {
		http.Error(w, "missing path", http.StatusBadRequest)
		return
	}
	clean := filepath.Clean(relPath)
	var fullPath string
	for _, base := range config.GetDatasetsSearchDirs() {
		absBase, err := filepath.Abs(base)
		if err != nil {
			continue
		}
		cand := filepath.Clean(filepath.Join(absBase, clean))
		rel, err := filepath.Rel(absBase, cand)
		if err == nil && !strings.HasPrefix(rel, "..") {
			if _, err := os.Stat(cand); err == nil {
				fullPath = cand
				break
			}
		}
	}
	if fullPath == "" {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	if strings.HasSuffix(fullPath, ".png") {
		w.Header().Set("Content-Type", "image/png")
	} else {
		w.Header().Set("Content-Type", "image/jpeg")
	}
	http.ServeFile(w, r, fullPath)
}
