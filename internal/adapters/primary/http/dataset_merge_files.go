package http

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type splitScan struct {
	imgDir, lblDir, split string
}

func buildDatasetSplitScans(srcDir string) []splitScan {
	var toScan []splitScan
	hasStandardSplits := false
	for _, s := range []string{"train", "valid", "val", "test"} {
		iDir := filepath.Join(srcDir, s, "images")
		lDir := filepath.Join(srcDir, s, "labels")
		if _, err := os.Stat(iDir); err == nil {
			hasStandardSplits = true
			target := "train"
			if s == "valid" || s == "val" || s == "test" {
				target = "valid"
			}
			toScan = append(toScan, splitScan{imgDir: iDir, lblDir: lDir, split: target})
		}
	}

	if !hasStandardSplits {
		iDir, lDir := srcDir, filepath.Join(srcDir, "labels")
		if _, err := os.Stat(filepath.Join(srcDir, "images")); err == nil {
			iDir = filepath.Join(srcDir, "images")
		}
		if _, err := os.Stat(lDir); err != nil {
			lDir = iDir
		}
		toScan = append(toScan, splitScan{imgDir: iDir, lblDir: lDir, split: "flat"})
	}
	return toScan
}

func writeUnifiedDataYaml(outDir string, classes []string) string {
	yamlPath := filepath.Join(outDir, "data.yaml")
	var namesYaml strings.Builder
	for i, c := range classes {
		namesYaml.WriteString(fmt.Sprintf("  %d: %s\n", i, c))
	}
	yamlContent := fmt.Sprintf("# Ultralytics YOLO Unified Dataset\npath: %s\ntrain: train/images\nval: valid/images\ntest: valid/images\n\nnc: %d\nnames:\n%s", outDir, len(classes), namesYaml.String())
	_ = os.WriteFile(yamlPath, []byte(yamlContent), 0644)
	return yamlPath
}
