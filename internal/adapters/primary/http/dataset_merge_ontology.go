package http

import (
	"os"
	"strconv"
	"strings"

	"hydraforge/internal/domain"
)

// canonicalCategory normalizes a label string into a standard visual taxonomy.
func canonicalCategory(s string) string {
	low := strings.ToLower(strings.TrimSpace(s))
	if low == "" || low == "ignore" || low == "ignorar" || low == "-1" || low == "drop" {
		return "ignore"
	}
	if strings.Contains(low, "cell") || strings.Contains(low, "phone") || strings.Contains(low, "celular") {
		return "cell-phone"
	}
	if strings.Contains(low, "truck") || strings.Contains(low, "caminh") || strings.Contains(low, "lorry") || strings.Contains(low, "pickup") {
		return "truck"
	}
	if strings.Contains(low, "bus") || strings.Contains(low, "onibus") || strings.Contains(low, "ônibus") {
		return "bus"
	}
	if strings.Contains(low, "moto") || strings.Contains(low, "bike") || strings.Contains(low, "cycle") || strings.Contains(low, "scooter") {
		return "motorcycle"
	}
	if strings.Contains(low, "person") || strings.Contains(low, "pedestrian") || strings.Contains(low, "pessoa") || strings.Contains(low, "people") || strings.Contains(low, "pedestre") {
		return "person"
	}
	if strings.Contains(low, "bicycle") || strings.Contains(low, "bicicleta") || strings.Contains(low, "ciclista") || strings.Contains(low, "cyclist") {
		return "bicycle"
	}
	if strings.Contains(low, "car") || strings.Contains(low, "van") || strings.Contains(low, "auto") || strings.Contains(low, "carro") || strings.Contains(low, "sedan") || strings.Contains(low, "suv") || strings.Contains(low, "taxi") || strings.Contains(low, "veiculo") || strings.Contains(low, "vehicle") {
		return "car"
	}
	return low
}

// parseDatasetYamlClasses reads class names from a dataset data.yaml.
func parseDatasetYamlClasses(yamlPath string) []string {
	yBytes, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil
	}
	var classes []string
	inNames := false
	for _, line := range strings.Split(string(yBytes), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "names:") {
			rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "names:"))
			if strings.HasPrefix(rest, "[") && strings.HasSuffix(rest, "]") {
				for _, p := range strings.Split(strings.Trim(rest, "[]"), ",") {
					c := strings.Trim(strings.TrimSpace(p), "'\"")
					if c != "" {
						classes = append(classes, c)
					}
				}
				if len(classes) > 0 {
					break
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
					v := strings.Trim(strings.TrimSpace(parts[1]), "'\"")
					if v != "" {
						classes = append(classes, v)
					}
				}
			} else if len(trimmed) > 0 && !strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
				break
			}
		}
	}
	return classes
}

// rewriteMergedLabel reads YOLO bounding boxes, remaps class IDs, and writes normalized labels.
func rewriteMergedLabel(
	srcLbl, dstLbl string,
	dsMap map[string]string,
	srcDataset *domain.Dataset,
	dsID string,
	classIdxMap map[string]string,
	reqClassesCount int,
) {
	data, err := os.ReadFile(srcLbl)
	if err != nil {
		return
	}
	var newLines []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 5 {
			oldCid := parts[0]
			targetClass := ""

			if mapped, ok := dsMap[oldCid]; ok && mapped != "" {
				targetClass = mapped
			} else if srcDataset != nil && len(srcDataset.Classes) > 0 {
				if idx, err := strconv.Atoi(oldCid); err == nil && idx >= 0 && idx < len(srcDataset.Classes) {
					origName := srcDataset.Classes[idx]
					if mapped, ok := dsMap[origName]; ok && mapped != "" {
						targetClass = mapped
					} else if mapped, ok := dsMap[strings.ToLower(origName)]; ok && mapped != "" {
						targetClass = mapped
					} else if mapped, ok := dsMap[strconv.Itoa(idx)]; ok && mapped != "" {
						targetClass = mapped
					} else {
						targetClass = origName
					}
				}
			}

			if targetClass == "" || targetClass == oldCid {
				dsLow := strings.ToLower(dsID)
				if strings.Contains(dsLow, "truck") {
					targetClass = "truck"
				} else if strings.Contains(dsLow, "moto") || strings.Contains(dsLow, "bike") {
					targetClass = "motorcycle"
				} else if strings.Contains(dsLow, "bus") {
					targetClass = "bus"
				} else if srcDataset != nil && len(srcDataset.Classes) > 0 {
					if idx, err := strconv.Atoi(oldCid); err == nil && idx >= 0 && idx < len(srcDataset.Classes) {
						targetClass = srcDataset.Classes[idx]
					}
				}
			}

			targetCanon := canonicalCategory(targetClass)
			if targetCanon == "ignore" || targetClass == "ignore" || targetClass == "ignorar" || targetClass == "-1" {
				continue
			}

			newCid := ""
			if idxStr, exists := classIdxMap[targetClass]; exists {
				newCid = idxStr
			} else if idxStr, exists := classIdxMap[strings.ToLower(targetClass)]; exists {
				newCid = idxStr
			} else if idxStr, exists := classIdxMap[targetCanon]; exists {
				newCid = idxStr
			} else if reqClassesCount == 1 {
				newCid = "0"
			}

			if newCid != "" {
				newLines = append(newLines, newCid+" "+strings.Join(parts[1:], " "))
			}
		}
	}
	_ = os.WriteFile(dstLbl, []byte(strings.Join(newLines, "\n")), 0644)
}
