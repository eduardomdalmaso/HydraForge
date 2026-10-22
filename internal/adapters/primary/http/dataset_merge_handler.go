package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"hydraforge/internal/config"
	"hydraforge/internal/domain"
)

// HandleDatasetMerge handles POST /api/v1/training/datasets/merge.
func (h *TrainingHandler) HandleDatasetMerge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"POST required"}`, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		TargetName string                       `json:"target_name"`
		DatasetIDs []string                     `json:"dataset_ids"`
		Mappings   map[string]map[string]string `json:"mappings"`
		Classes    []string                     `json:"classes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid JSON request"}`, http.StatusBadRequest)
		return
	}

	if len(req.DatasetIDs) == 0 {
		http.Error(w, `{"error":"at least 1 dataset required"}`, http.StatusBadRequest)
		return
	}

	cleanTarget := strings.ToLower(strings.ReplaceAll(req.TargetName, " ", "_"))
	if cleanTarget == "" {
		cleanTarget = fmt.Sprintf("merged_%d", time.Now().Unix())
	}
	outDir := filepath.Join(config.GetDatasetsDir(), cleanTarget)
	_ = os.RemoveAll(outDir)
	_ = os.MkdirAll(filepath.Join(outDir, "train", "images"), 0755)
	_ = os.MkdirAll(filepath.Join(outDir, "train", "labels"), 0755)
	_ = os.MkdirAll(filepath.Join(outDir, "valid", "images"), 0755)
	_ = os.MkdirAll(filepath.Join(outDir, "valid", "labels"), 0755)

	if len(req.Classes) == 0 {
		req.Classes = []string{"car", "motorcycle", "truck", "bus"}
	}

	classIdxMap := make(map[string]string)
	for i, c := range req.Classes {
		idxStr := strconv.Itoa(i)
		classIdxMap[c] = idxStr
		classIdxMap[strings.ToLower(c)] = idxStr
		classIdxMap[canonicalCategory(c)] = idxStr
		classIdxMap[idxStr] = idxStr
	}

	trainCount, valCount := 0, 0

	for _, dsID := range req.DatasetIDs {
		var srcDir string
		for _, base := range config.GetDatasetsSearchDirs() {
			cand := filepath.Join(base, dsID)
			if fi, err := os.Stat(cand); err == nil && fi.IsDir() {
				srcDir = cand
				break
			}
		}
		if srcDir == "" {
			srcDir = filepath.Join(config.GetDatasetsDir(), dsID)
		}
		prefix := dsID + "_"
		dsMap := req.Mappings[dsID]
		if dsMap == nil {
			dsMap = make(map[string]string)
		}
		srcDataset, _ := h.useCase.GetDataset(r.Context(), dsID)

		if srcDataset == nil || len(srcDataset.Classes) == 0 {
			yamlClasses := parseDatasetYamlClasses(filepath.Join(srcDir, "data.yaml"))
			if len(yamlClasses) > 0 {
				srcDataset = &domain.Dataset{DatasetID: dsID, Name: dsID, Classes: yamlClasses, NumClasses: len(yamlClasses)}
			}
		}

		toScan := buildDatasetSplitScans(srcDir)

		for _, sc := range toScan {
			entries, err := os.ReadDir(sc.imgDir)
			if err != nil {
				continue
			}

			fileIdx := 0
			for _, e := range entries {
				if e.IsDir() {
					continue
				}
				imgName := e.Name()
				ext := strings.ToLower(filepath.Ext(imgName))
				if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
					continue
				}
				fileIdx++

				targetSplit := sc.split
				if targetSplit == "flat" {
					if fileIdx%10 == 0 {
						targetSplit = "valid"
					} else {
						targetSplit = "train"
					}
				}

				base := strings.TrimSuffix(imgName, filepath.Ext(imgName))
				dstImg := filepath.Join(outDir, targetSplit, "images", prefix+imgName)
				srcImg := filepath.Join(sc.imgDir, imgName)
				_ = os.Remove(dstImg)
				_ = os.Symlink(srcImg, dstImg)

				srcLbl := filepath.Join(sc.lblDir, base+".txt")
				dstLbl := filepath.Join(outDir, targetSplit, "labels", prefix+base+".txt")
				rewriteMergedLabel(srcLbl, dstLbl, dsMap, srcDataset, dsID, classIdxMap, len(req.Classes))

				if targetSplit == "train" {
					trainCount++
				} else {
					valCount++
				}
			}
		}
	}

	yamlPath := writeUnifiedDataYaml(outDir, req.Classes)

	ds := &domain.Dataset{
		DatasetID:   cleanTarget,
		Name:        req.TargetName,
		Task:        domain.TaskDetect,
		YAMLPath:    yamlPath,
		Classes:     req.Classes,
		NumClasses:  len(req.Classes),
		TrainImages: trainCount,
		ValImages:   valCount,
		CreatedAt:   time.Now(),
	}

	saved, err := h.useCase.RegisterDataset(r.Context(), ds)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(saved)
}
