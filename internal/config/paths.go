package config

import (
	"os"
	"os/exec"
	"path/filepath"
)

// GetHomeDir returns the current user's home directory.
func GetHomeDir() string {
	h, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return h
}

// GetDatasetsDir returns the primary datasets storage directory.
func GetDatasetsDir() string {
	if env := os.Getenv("DATASETS_DIR"); env != "" {
		return filepath.Clean(env)
	}
	home := GetHomeDir()
	candidates := []string{
		filepath.Join(home, "datasets"),
		"datasets",
	}
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && fi.IsDir() {
			return c
		}
	}
	return filepath.Join(home, "datasets")
}

// GetDatasetsSearchDirs returns all potential dataset locations on the system.
func GetDatasetsSearchDirs() []string {
	home := GetHomeDir()
	dirs := []string{
		"datasets",
		filepath.Join(home, "datasets"),
		filepath.Join(home, "Documents/HydraForge/datasets"),
	}
	if env := os.Getenv("DATASETS_DIR"); env != "" {
		dirs = append([]string{filepath.Clean(env)}, dirs...)
	}
	return dirs
}

// GetRunsSearchDirs returns potential directories where training runs are stored.
func GetRunsSearchDirs() []string {
	home := GetHomeDir()
	dirs := []string{
		"runs/train",
		filepath.Join(home, "runs/train"),
		filepath.Join(home, "Documents/HydraForge/runs/train"),
	}
	if env := os.Getenv("RUNS_DIR"); env != "" {
		dirs = append([]string{filepath.Clean(env)}, dirs...)
	}
	return dirs
}

// GetPythonBin discovers the active Python binary for YOLO execution.
func GetPythonBin() string {
	if env := os.Getenv("PYTHON_BIN"); env != "" {
		if _, err := os.Stat(env); err == nil {
			return env
		}
	}
	if conda := os.Getenv("CONDA_PREFIX"); conda != "" {
		p := filepath.Join(conda, "bin", "python")
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	home := GetHomeDir()
	candidates := []string{
		filepath.Join(home, "miniconda3/envs/analytics-env/bin/python"),
		filepath.Join(home, "miniconda3/envs/hydraforge/bin/python"),
		filepath.Join(home, "anaconda3/envs/analytics-env/bin/python"),
		filepath.Join(home, ".conda/envs/analytics-env/bin/python"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	if p, err := exec.LookPath("python3"); err == nil {
		return p
	}
	return "python3"
}

// GetHydraStreamSamplesDir returns the location of HydraStream sample images.
func GetHydraStreamSamplesDir() string {
	if env := os.Getenv("HYDRASTREAM_SAMPLES_DIR"); env != "" {
		return filepath.Clean(env)
	}
	home := GetHomeDir()
	candidates := []string{
		"../HydraStream/samples",
		filepath.Join(home, "Documents/HydraStream/samples"),
		"samples",
	}
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && fi.IsDir() {
			return c
		}
	}
	return filepath.Join(home, "Documents/HydraStream/samples")
}
