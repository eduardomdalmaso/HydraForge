package domain

import "time"

// MediaFile represents a single video file within a media category folder.
type MediaFile struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Folder    string    `json:"folder"`
	Path      string    `json:"path"`
	SizeBytes int64     `json:"size_bytes"`
	ModTime   time.Time `json:"mod_time"`
	StreamURL string    `json:"stream_url"`
}

// MediaFolder represents a folder containing grouped media files.
type MediaFolder struct {
	Name       string      `json:"name"`
	Label      string      `json:"label"`
	FileCount  int         `json:"file_count"`
	TotalBytes int64       `json:"total_bytes"`
	Files      []MediaFile `json:"files"`
}

// MediaSourcesResponse represents the full hierarchy of categorized media.
type MediaSourcesResponse struct {
	TotalFolders int           `json:"total_folders"`
	TotalFiles   int           `json:"total_files"`
	Folders      []MediaFolder `json:"folders"`
}
