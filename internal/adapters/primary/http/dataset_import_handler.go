package http

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"hydraforge/internal/config"
	"hydraforge/internal/domain"
)

// HandleDatasetImport handles POST /api/v1/training/datasets/import (multipart zip upload).
func (h *TrainingHandler) HandleDatasetImport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"POST required"}`, http.StatusMethodNotAllowed)
		return
	}

	// 1GB max multipart memory
	if err := r.ParseMultipartForm(1024 << 20); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to parse form: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		files = r.MultipartForm.File["file"]
	}
	if len(files) == 0 {
		http.Error(w, `{"error":"no zip files provided in multipart request"}`, http.StatusBadRequest)
		return
	}

	registered := make([]*domain.Dataset, 0, len(files))
	datasetsDir := config.GetDatasetsDir()
	_ = os.MkdirAll(datasetsDir, 0755)

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		tmpZip, err := os.CreateTemp("", "upload-*.zip")
		if err != nil {
			continue
		}
		_, err = io.Copy(tmpZip, file)
		tmpZip.Close()
		if err != nil {
			os.Remove(tmpZip.Name())
			continue
		}

		cleanBase := strings.TrimSuffix(fileHeader.Filename, filepath.Ext(fileHeader.Filename))
		cleanBase = strings.ToLower(strings.ReplaceAll(cleanBase, " ", "_"))
		targetDir := filepath.Clean(filepath.Join(datasetsDir, cleanBase))
		_ = os.MkdirAll(targetDir, 0755)

		// Unzip with strict Zip Slip protection
		rZip, err := zip.OpenReader(tmpZip.Name())
		if err == nil {
			for _, f := range rZip.File {
				if strings.Contains(f.Name, "..") {
					continue
				}
				fpath := filepath.Clean(filepath.Join(targetDir, f.Name))
				if !strings.HasPrefix(fpath, targetDir+string(filepath.Separator)) && fpath != targetDir {
					continue
				}

				if f.FileInfo().IsDir() {
					os.MkdirAll(fpath, os.ModePerm)
					continue
				}
				os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
				outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
				if err != nil {
					continue
				}
				rc, err := f.Open()
				if err == nil {
					_, _ = io.Copy(outFile, rc)
					rc.Close()
				}
				outFile.Close()
			}
			rZip.Close()
		}
		os.Remove(tmpZip.Name())

		_, _ = h.useCase.RescanDatasets(r.Context())
		if ds, err := h.useCase.GetDataset(r.Context(), cleanBase); err == nil {
			registered = append(registered, ds)
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(registered)
}
