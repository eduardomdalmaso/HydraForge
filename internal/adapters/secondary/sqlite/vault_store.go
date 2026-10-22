package sqlite

import (
	"database/sql"
	"fmt"
)

// VaultStore bundles Vault dataset and frame repositories.
type VaultStore struct {
	DatasetRepo *VaultDatasetRepository
	FrameRepo   *VaultFrameRepository
}

// NewVaultStore creates and initializes Vault dataset and frame repositories.
func NewVaultStore(db *sql.DB) (*VaultStore, error) {
	if err := initVaultSchema(db); err != nil {
		return nil, fmt.Errorf("failed to initialize vault schema: %w", err)
	}
	return &VaultStore{
		DatasetRepo: &VaultDatasetRepository{db: db},
		FrameRepo:   &VaultFrameRepository{db: db},
	}, nil
}

func initVaultSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS vault_datasets (
		dataset_id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT NOT NULL DEFAULT '',
		task TEXT NOT NULL DEFAULT 'detect',
		classes_json TEXT NOT NULL DEFAULT '[]',
		train_images INTEGER NOT NULL DEFAULT 0,
		val_images INTEGER NOT NULL DEFAULT 0,
		test_images INTEGER NOT NULL DEFAULT 0,
		total_frames INTEGER NOT NULL DEFAULT 0,
		size_bytes INTEGER NOT NULL DEFAULT 0,
		yaml_path TEXT NOT NULL DEFAULT '',
		checksum_sha TEXT NOT NULL DEFAULT '',
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS frames (
		frame_id TEXT PRIMARY KEY,
		dataset_id TEXT NOT NULL,
		camera_id TEXT NOT NULL,
		file_path TEXT NOT NULL,
		file_name TEXT NOT NULL,
		phash INTEGER NOT NULL,
		width INTEGER NOT NULL,
		height INTEGER NOT NULL,
		size_bytes INTEGER NOT NULL,
		status TEXT NOT NULL DEFAULT 'RAW',
		split TEXT NOT NULL DEFAULT 'none',
		bboxes_json TEXT NOT NULL DEFAULT '[]',
		polygons_json TEXT NOT NULL DEFAULT '[]',
		source_event TEXT NOT NULL DEFAULT '',
		captured_at DATETIME NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_frames_dataset ON frames(dataset_id);
	CREATE INDEX IF NOT EXISTS idx_frames_status ON frames(status);
	CREATE INDEX IF NOT EXISTS idx_frames_split ON frames(split);
	CREATE INDEX IF NOT EXISTS idx_frames_camera ON frames(camera_id);
	CREATE INDEX IF NOT EXISTS idx_frames_phash ON frames(phash);
	`
	_, err := db.Exec(schema)
	return err
}
