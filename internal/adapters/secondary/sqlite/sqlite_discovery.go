package sqlite

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"hydraforge/internal/config"
)

func (s *SQLiteStore) RescanDatasets() {
	for _, dir := range config.GetDatasetsSearchDirs() {
		if _, err := os.Stat(dir); err == nil {
			s.autoDiscoverDatasets(dir)
		}
	}
	s.purgeStaleDatasets()
	for _, dir := range config.GetRunsSearchDirs() {
		if _, err := os.Stat(dir); err == nil {
			s.autoDiscoverTrainingRuns(dir)
		}
	}
}

func (s *SQLiteStore) autoDiscoverDatasets(datasetsDir string) {
	entries, err := os.ReadDir(datasetsDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		id := entry.Name()
		fullPath := filepath.Join(datasetsDir, id)
		fi, err := os.Stat(fullPath)
		if err != nil || !fi.IsDir() {
			continue
		}

		yamlPath := filepath.Join(fullPath, "data.yaml")
		if _, err := os.Stat(yamlPath); err != nil {
			imgCount := countImages(fullPath)
			if imgCount > 0 {
				clsName := strings.ReplaceAll(id, "_", " ")
				clsName = strings.ReplaceAll(clsName, "-", " ")
				content := fmt.Sprintf("path: %s\ntrain: .\nval: .\nnc: 1\nnames:\n  0: %s\n", fullPath, clsName)
				_ = os.WriteFile(yamlPath, []byte(content), 0644)
			}
		}

		if _, err := os.Stat(yamlPath); err == nil {
			sanitizeYamlPaths(yamlPath, fullPath)
			trainCount := countDir(filepath.Join(fullPath, "images", "train"))
			if trainCount == 0 {
				trainCount = countDir(filepath.Join(fullPath, "train", "images"))
			}
			valCount := countDir(filepath.Join(fullPath, "images", "val"))
			if valCount == 0 {
				valCount = countDir(filepath.Join(fullPath, "val", "images"))
			}
			if valCount == 0 {
				valCount = countDir(filepath.Join(fullPath, "images", "valid"))
			}
			if valCount == 0 {
				valCount = countDir(filepath.Join(fullPath, "valid", "images"))
			}
			testCount := countDir(filepath.Join(fullPath, "images", "test"))
			if testCount == 0 {
				testCount = countDir(filepath.Join(fullPath, "test", "images"))
			}
			if trainCount == 0 && valCount == 0 {
				trainCount = countImages(fullPath)
			}
			sizeBytes := calcDirSize(fullPath)
			name := strings.ReplaceAll(id, "_", " ")
			name = strings.ReplaceAll(name, "-", " ")
			name = strings.ToUpper(name[:1]) + name[1:]
			classes := parseYamlClasses(yamlPath)
			classesJSON, _ := json.Marshal(classes)

			query := `INSERT INTO registered_datasets 
				(dataset_id, name, task, yaml_path, classes_json, num_classes, train_images, val_images, test_images, size_bytes, created_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
				ON CONFLICT(dataset_id) DO UPDATE SET 
					name=excluded.name, yaml_path=excluded.yaml_path,
					classes_json=excluded.classes_json, num_classes=excluded.num_classes,
					train_images=excluded.train_images, val_images=excluded.val_images, test_images=excluded.test_images,
					size_bytes=excluded.size_bytes`
			_, _ = s.db.Exec(query, id, name, "detect", yamlPath, string(classesJSON), len(classes), trainCount, valCount, testCount, sizeBytes, time.Now())
		}
	}
}

func (s *SQLiteStore) purgeStaleDatasets() {
	rows, err := s.db.Query("SELECT dataset_id, yaml_path FROM registered_datasets")
	if err != nil {
		return
	}
	defer rows.Close()

	var toDelete []string
	for rows.Next() {
		var dsID, yPath string
		if err := rows.Scan(&dsID, &yPath); err == nil {
			if yPath != "" {
				if _, err := os.Stat(yPath); os.IsNotExist(err) {
					if _, errDir := os.Stat(filepath.Dir(yPath)); os.IsNotExist(errDir) {
						toDelete = append(toDelete, dsID)
					}
				}
			}
		}
	}
	if err := rows.Err(); err != nil {
		return
	}
	_ = rows.Close()

	for _, staleID := range toDelete {
		_, _ = s.db.Exec("DELETE FROM registered_datasets WHERE dataset_id = ?", staleID)
	}
}
