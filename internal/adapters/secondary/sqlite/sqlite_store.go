package sqlite

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"
)

// SQLiteStore provides persistent storage for HydraForge entities with primary keys.
type SQLiteStore struct {
	db *sql.DB
	mu sync.RWMutex
}

// NewSQLiteStore initializes a SQLite database connection and runs auto-migrations.
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)")
	if err != nil {
		return nil, err
	}

	store := &SQLiteStore{db: db}
	if err := store.migrate(); err != nil {
		return nil, err
	}

	store.RescanDatasets()
	return store, nil
}

// Close closes the underlying database connection.
func (s *SQLiteStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// DB returns the underlying sql.DB connection.
func (s *SQLiteStore) DB() *sql.DB {
	return s.db
}

func (s *SQLiteStore) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS registered_datasets (
		dataset_id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		task TEXT NOT NULL,
		yaml_path TEXT NOT NULL,
		classes_json TEXT NOT NULL,
		num_classes INTEGER NOT NULL,
		train_images INTEGER NOT NULL,
		val_images INTEGER NOT NULL,
		test_images INTEGER DEFAULT 0,
		size_bytes INTEGER DEFAULT 0,
		created_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS training_jobs (
		job_id TEXT PRIMARY KEY,
		model_architecture TEXT NOT NULL,
		task TEXT NOT NULL,
		dataset_id TEXT NOT NULL,
		dataset_path TEXT,
		hyperparameters_json TEXT NOT NULL,
		status TEXT NOT NULL,
		current_epoch INTEGER DEFAULT 0,
		total_epochs INTEGER DEFAULT 50,
		best_map50 REAL DEFAULT 0.0,
		total_energy_kwh REAL DEFAULT 0.0,
		avg_power_watts REAL DEFAULT 0.0,
		peak_vram_mb REAL DEFAULT 0.0,
		avg_fps REAL DEFAULT 0.0,
		duration_sec REAL DEFAULT 0.0,
		output_weights TEXT,
		error_message TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY(dataset_id) REFERENCES registered_datasets(dataset_id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS epoch_metrics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		job_id TEXT NOT NULL,
		epoch INTEGER NOT NULL,
		total_epochs INTEGER NOT NULL,
		box_loss REAL,
		cls_loss REAL,
		dfl_loss REAL,
		val_box_loss REAL,
		val_cls_loss REAL,
		map50 REAL,
		map50_95 REAL,
		precision REAL,
		recall REAL,
		learning_rate REAL,
		gpu_vram_mb REAL,
		power_watts REAL,
		temp_celsius REAL,
		gpu_util_pct REAL,
		fps REAL,
		epoch_duration_sec REAL,
		created_at DATETIME NOT NULL,
		UNIQUE(job_id, epoch),
		FOREIGN KEY(job_id) REFERENCES training_jobs(job_id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS model_checkpoints (
		checkpoint_id TEXT PRIMARY KEY,
		job_id TEXT NOT NULL,
		architecture TEXT NOT NULL,
		task TEXT NOT NULL,
		epoch INTEGER NOT NULL,
		map50 REAL,
		map50_95 REAL,
		weights_path TEXT NOT NULL,
		export_format TEXT,
		export_path TEXT,
		size_bytes INTEGER NOT NULL,
		precision TEXT,
		created_at DATETIME NOT NULL,
		FOREIGN KEY(job_id) REFERENCES training_jobs(job_id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS benchmark_jobs (
		job_id TEXT PRIMARY KEY,
		model TEXT NOT NULL,
		data TEXT NOT NULL DEFAULT '',
		model_architecture TEXT DEFAULT '',
		weights_path TEXT DEFAULT '',
		imgsz INTEGER NOT NULL,
		quantize INTEGER NOT NULL,
		device TEXT NOT NULL,
		target_formats_json TEXT NOT NULL,
		status TEXT NOT NULL,
		error_message TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS benchmark_results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		job_id TEXT NOT NULL,
		format TEXT NOT NULL,
		status TEXT NOT NULL,
		latency_ms REAL NOT NULL,
		fps REAL NOT NULL,
		size_mb REAL NOT NULL,
		map50_95 REAL NOT NULL,
		created_at DATETIME NOT NULL,
		FOREIGN KEY(job_id) REFERENCES benchmark_jobs(job_id) ON DELETE CASCADE
	);
	`
	_, err := s.db.Exec(schema)
	if err != nil {
		return err
	}
	_, _ = s.db.Exec(`ALTER TABLE training_jobs ADD COLUMN current_batch INTEGER DEFAULT 0;`)
	_, _ = s.db.Exec(`ALTER TABLE training_jobs ADD COLUMN total_batches INTEGER DEFAULT 0;`)
	_, _ = s.db.Exec(`ALTER TABLE training_jobs ADD COLUMN best_map50_95 REAL DEFAULT 0.0;`)
	_, _ = s.db.Exec(`ALTER TABLE training_jobs ADD COLUMN tenant_id TEXT DEFAULT 'default';`)
	_, _ = s.db.Exec(`ALTER TABLE registered_datasets ADD COLUMN tenant_id TEXT DEFAULT 'default';`)
	_, _ = s.db.Exec(`ALTER TABLE model_checkpoints ADD COLUMN tenant_id TEXT DEFAULT 'default';`)
	_, _ = s.db.Exec(`ALTER TABLE benchmark_jobs ADD COLUMN tenant_id TEXT DEFAULT 'default';`)
	_, _ = s.db.Exec(`ALTER TABLE benchmark_jobs ADD COLUMN model TEXT DEFAULT '';`)
	_, _ = s.db.Exec(`ALTER TABLE benchmark_jobs ADD COLUMN data TEXT DEFAULT '';`)
	return nil
}
