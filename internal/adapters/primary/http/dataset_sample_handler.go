package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"hydraforge/internal/config"
)


// HandleDatasetSample returns sample images with bounding boxes for a given class.
func (h *TrainingHandler) HandleDatasetSample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dsID := r.URL.Query().Get("id")
	targetClass := r.URL.Query().Get("class")
	if dsID == "" {
		http.Error(w, `{"error":"id parameter required"}`, http.StatusBadRequest)
		return
	}

	var dsDir string
	for _, base := range config.GetDatasetsSearchDirs() {
		cand := filepath.Join(base, dsID)
		if fi, err := os.Stat(cand); err == nil && fi.IsDir() {
			dsDir = cand
			break
		}
	}
	if dsDir == "" {
		dsDir = filepath.Join("datasets", dsID)
	}

	classNamesMap := make(map[string]string)
	targetClassIdx := targetClass
	if ds, err := h.useCase.GetDataset(r.Context(), dsID); err == nil && ds != nil {
		for idx, c := range ds.Classes {
			sIdx := strconv.Itoa(idx)
			classNamesMap[sIdx] = c
			classNamesMap[c] = c
			if strings.EqualFold(c, targetClass) {
				targetClassIdx = sIdx
			}
		}
	}

	pairs := []candidatePair{
		{filepath.Join(dsDir, "images", "train"), filepath.Join(dsDir, "labels", "train"), filepath.Join("images", "train")},
		{filepath.Join(dsDir, "train", "images"), filepath.Join(dsDir, "train", "labels"), filepath.Join("train", "images")},
		{filepath.Join(dsDir, "images", "val"), filepath.Join(dsDir, "labels", "val"), filepath.Join("images", "val")},
		{filepath.Join(dsDir, "val", "images"), filepath.Join(dsDir, "val", "labels"), filepath.Join("val", "images")},
		{filepath.Join(dsDir, "images", "valid"), filepath.Join(dsDir, "labels", "valid"), filepath.Join("images", "valid")},
		{filepath.Join(dsDir, "valid", "images"), filepath.Join(dsDir, "valid", "labels"), filepath.Join("valid", "images")},
		{filepath.Join(dsDir, "images"), filepath.Join(dsDir, "labels"), "images"},
		{dsDir, dsDir, ""},
	}

	var allMatches []*SampleMatch
	seenImgs := make(map[string]bool)

	for _, pair := range pairs {
		entries, err := os.ReadDir(pair.labelsDir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".txt") {
				continue
			}
			data, err := os.ReadFile(filepath.Join(pair.labelsDir, e.Name()))
			if err != nil {
				continue
			}
			lines := strings.Split(string(data), "\n")
			hasTarget := false
			var imgAnnotations []BoxAnnotation
			var primaryTargetBBox []float64

			for _, line := range lines {
				parts := strings.Fields(line)
				if len(parts) >= 5 {
					cid := parts[0]
					x, e1 := strconv.ParseFloat(parts[1], 64)
					y, e2 := strconv.ParseFloat(parts[2], 64)
					w, e3 := strconv.ParseFloat(parts[3], 64)
					h, e4 := strconv.ParseFloat(parts[4], 64)
					if e1 == nil && e2 == nil && e3 == nil && e4 == nil {
						isTgt := (targetClass == "" || cid == targetClass || cid == targetClassIdx)
						if isTgt {
							hasTarget = true
							if len(primaryTargetBBox) == 0 {
								primaryTargetBBox = []float64{x, y, w, h}
							}
						}
						cName := classNamesMap[cid]
						if cName == "" {
							cName = fmt.Sprintf("class_%s", cid)
						}
						imgAnnotations = append(imgAnnotations, BoxAnnotation{
							ClassID:   cid,
							ClassName: cName,
							BBox:      []float64{x, y, w, h},
							IsTarget:  isTgt,
						})
					}
				}
			}

			if hasTarget {
				baseName := strings.TrimSuffix(e.Name(), ".txt")
				var imgFile string
				for _, ext := range []string{".jpg", ".png", ".jpeg", ".JPG", ".PNG"} {
					if _, err := os.Stat(filepath.Join(pair.imagesDir, baseName+ext)); err == nil {
						imgFile = baseName + ext
						break
					}
				}
				if imgFile != "" && !seenImgs[imgFile] {
					seenImgs[imgFile] = true
					relPath := filepath.ToSlash(filepath.Join(dsID, pair.subRel, imgFile))
					allMatches = append(allMatches, &SampleMatch{
						ImageURL:    fmt.Sprintf("/api/v1/training/datasets/image?path=%s", relPath),
						Filename:    imgFile,
						BBox:        primaryTargetBBox,
						ClassID:     targetClass,
						Annotations: imgAnnotations,
					})
				}
			}
		}
	}

	if len(allMatches) == 0 {
		for _, pair := range pairs {
			for _, ext := range []string{".jpg", ".png", ".jpeg", ".JPG", ".PNG"} {
				imgs, _ := filepath.Glob(filepath.Join(pair.imagesDir, "*"+ext))
				for _, picked := range imgs {
					imgBase := filepath.Base(picked)
					if seenImgs[imgBase] {
						continue
					}
					seenImgs[imgBase] = true
					relPath := filepath.ToSlash(filepath.Join(dsID, pair.subRel, imgBase))
					defaultBox := []float64{0.5, 0.5, 0.6, 0.6}
					allMatches = append(allMatches, &SampleMatch{
						ImageURL: fmt.Sprintf("/api/v1/training/datasets/image?path=%s", relPath),
						Filename: imgBase,
						BBox:     defaultBox,
						ClassID:  targetClass,
						Annotations: []BoxAnnotation{
							{ClassID: targetClass, ClassName: targetClass, BBox: defaultBox, IsTarget: true},
						},
					})
					if len(allMatches) >= 15 {
						break
					}
				}
				if len(allMatches) >= 15 {
					break
				}
			}
			if len(allMatches) > 0 {
				break
			}
		}
	}

	if len(allMatches) == 0 {
		http.Error(w, `{"error":"no images found in dataset"}`, http.StatusNotFound)
		return
	}

	if len(allMatches) > 20 {
		allMatches = allMatches[:20]
	}

	json.NewEncoder(w).Encode(allMatches)
}
