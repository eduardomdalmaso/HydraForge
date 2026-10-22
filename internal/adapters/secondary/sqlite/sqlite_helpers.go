package sqlite

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func countImages(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() {
			ext := strings.ToLower(filepath.Ext(e.Name()))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" || ext == ".bmp" {
				count++
			}
		}
	}
	return count
}

func sanitizeYamlPaths(yamlPath, datasetDir string) {
	content, err := os.ReadFile(yamlPath)
	if err != nil {
		return
	}
	s := string(content)
	modified := false
	if strings.Contains(s, "../train/images") {
		s = strings.ReplaceAll(s, "../train/images", "train/images")
		modified = true
	}
	if strings.Contains(s, "../valid/images") {
		s = strings.ReplaceAll(s, "../valid/images", "valid/images")
		modified = true
	}
	if strings.Contains(s, "../test/images") {
		s = strings.ReplaceAll(s, "../test/images", "test/images")
		modified = true
	}
	if !strings.Contains(s, "path:") {
		s = fmt.Sprintf("path: %s\n%s", datasetDir, s)
		modified = true
	}
	if modified {
		_ = os.WriteFile(yamlPath, []byte(s), 0644)
	}
}

func countDir(dirPath string) int {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() {
			count++
		}
	}
	return count
}

func calcDirSize(path string) int64 {
	var size int64
	_ = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err == nil && info != nil && !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}

func parseYamlClasses(yamlPath string) []string {
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return []string{"object"}
	}
	lines := strings.Split(string(data), "\n")
	var classes []string
	inNames := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "names:") {
			rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "names:"))
			if strings.HasPrefix(rest, "[") && strings.HasSuffix(rest, "]") {
				inner := strings.Trim(rest, "[]")
				parts := strings.Split(inner, ",")
				for _, p := range parts {
					c := strings.TrimSpace(p)
					c = strings.Trim(c, "'\"")
					if c != "" {
						classes = append(classes, c)
					}
				}
				if len(classes) > 0 {
					return classes
				}
			}
			inNames = true
			continue
		}
		if inNames {
			if strings.HasPrefix(trimmed, "-") {
				classes = append(classes, strings.TrimSpace(strings.TrimPrefix(trimmed, "-")))
			} else if strings.Contains(trimmed, ":") && !strings.HasPrefix(trimmed, "#") {
				parts := strings.SplitN(trimmed, ":", 2)
				if len(parts) == 2 {
					val := strings.TrimSpace(parts[1])
					val = strings.Trim(val, "'\"")
					if val != "" {
						classes = append(classes, val)
					}
				}
			} else if len(trimmed) > 0 && !strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
				break
			}
		}
	}
	if len(classes) == 0 {
		return []string{"object"}
	}
	return classes
}
