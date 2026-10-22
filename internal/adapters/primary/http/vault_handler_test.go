package http_test

import (
	"bytes"
	"context"
	"database/sql"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	primaryHttp "hydraforge/internal/adapters/primary/http"
	yoloExporter "hydraforge/internal/adapters/secondary/exporter/yolo"
	fsAdapter "hydraforge/internal/adapters/secondary/filestore"
	phashAdapter "hydraforge/internal/adapters/secondary/phash"
	sqliteAdapter "hydraforge/internal/adapters/secondary/sqlite"
	"hydraforge/internal/application"
)

func createTestImage(t *testing.T, w, h int, c color.Color) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		t.Fatalf("failed to encode test jpeg: %v", err)
	}
	return buf.Bytes()
}

func setupTestVaultApp(t *testing.T) (http.Handler, func()) {
	tmpDir, err := os.MkdirTemp("", "forge_vault_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	dbPath := filepath.Join(tmpDir, "test_forge_vault.db")
	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)")
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}

	fileStore, err := fsAdapter.NewFileStore(filepath.Join(tmpDir, "datasets"))
	if err != nil {
		t.Fatalf("failed to create file store: %v", err)
	}

	vaultStore, err := sqliteAdapter.NewVaultStore(db)
	if err != nil {
		t.Fatalf("failed to initialize vault store: %v", err)
	}

	dedup := phashAdapter.NewDeduplicator()
	exporter := yoloExporter.NewExporter(filepath.Join(tmpDir, "datasets", "curated"))

	datasetService := application.NewVaultDatasetService(vaultStore.DatasetRepo, vaultStore.FrameRepo, exporter)
	ingestService := application.NewVaultIngestService(vaultStore.FrameRepo, vaultStore.DatasetRepo, fileStore, dedup)

	vaultHandler := primaryHttp.NewVaultHandler(datasetService, ingestService, fileStore)

	mux := http.NewServeMux()
	primaryHttp.RegisterVaultRoutes(mux, vaultHandler)

	cleanup := func() {
		db.Close()
		_ = os.RemoveAll(tmpDir)
	}

	return mux, cleanup
}

func TestHTTP_VaultIngestAndExportFlow(t *testing.T) {
	handler, cleanup := setupTestVaultApp(t)
	defer cleanup()

	// 1. Create dataset via API
	createJSON := `{"dataset_id":"ds_ppe_safety","name":"PPE Safety","task":"detect","classes":[{"id":0,"name":"helmet"},{"id":1,"name":"vest"}]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/vault/datasets", bytes.NewBufferString(createJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 on dataset creation, got %d: %s", w.Code, w.Body.String())
	}

	// 2. Ingest Image with Bounding Boxes
	imgBytes := createTestImage(t, 640, 480, color.RGBA{R: 200, G: 50, B: 50, A: 255})
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	_ = writer.WriteField("dataset_id", "ds_ppe_safety")
	_ = writer.WriteField("camera_id", "cam_entrance_01")
	_ = writer.WriteField("bboxes", `[{"class_id":0,"class_name":"helmet","x_center":0.5,"y_center":0.5,"width":0.2,"height":0.2}]`)
	part, _ := writer.CreateFormFile("image", "frame01.jpg")
	_, _ = io.Copy(part, bytes.NewReader(imgBytes))
	_ = writer.Close()

	req = httptest.NewRequest(http.MethodPost, "/api/v1/vault/inbox/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 on ingest upload, got %d: %s", w.Code, w.Body.String())
	}

	// 3. Test Deduplication (Uploading identical image should be conflict)
	var bodyDup bytes.Buffer
	writerDup := multipart.NewWriter(&bodyDup)
	_ = writerDup.WriteField("dataset_id", "ds_ppe_safety")
	_ = writerDup.WriteField("camera_id", "cam_entrance_01")
	partDup, _ := writerDup.CreateFormFile("image", "frame01_dup.jpg")
	_, _ = io.Copy(partDup, bytes.NewReader(imgBytes))
	_ = writerDup.Close()

	req = httptest.NewRequest(http.MethodPost, "/api/v1/vault/inbox/upload", &bodyDup)
	req.Header.Set("Content-Type", writerDup.FormDataContentType())
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409 Conflict on perceptual duplicate, got %d", w.Code)
	}

	// 4. Export Dataset
	req = httptest.NewRequest(http.MethodPost, "/api/v1/vault/datasets/ds_ppe_safety/export", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 on export, got %d: %s", w.Code, w.Body.String())
	}
}

func TestFileStore_PathTraversalRejection(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fs_vault_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := fsAdapter.NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	traversalPaths := []string{
		"../outside.jpg",
		"../../etc/passwd",
		"inbox/../../secret.txt",
		"/etc/shadow",
		"\\windows\\system32",
	}

	for _, p := range traversalPaths {
		_, _, err := store.Save(context.Background(), p, bytes.NewReader([]byte("test")))
		if err == nil {
			t.Errorf("expected path traversal error for %s, got nil", p)
		}
	}
}
