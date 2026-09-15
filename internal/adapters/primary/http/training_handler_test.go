package http_test

import (
	"archive/zip"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	primaryHttp "hydraforge/internal/adapters/primary/http"
	"hydraforge/internal/adapters/secondary/gpu"
	"hydraforge/internal/adapters/secondary/sqlite"
	"hydraforge/internal/adapters/secondary/storage"
	"hydraforge/internal/adapters/secondary/worker"
	"hydraforge/internal/application"
)

func init() {
	os.Setenv("JWT_SECRET", "super-secret-key-that-is-at-least-32-chars-long-for-testing!")
	os.Setenv("SERVICE_API_KEY", "test-hydraforge-api-key")
}

func generateTestToken(tenantID, role string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payloadData, _ := json.Marshal(map[string]interface{}{
		"user_id":   "test-user-forge",
		"tenant_id": tenantID,
		"role":      role,
		"exp":       time.Now().Add(1 * time.Hour).Unix(),
	})
	payload := base64.RawURLEncoding.EncodeToString(payloadData)

	mac := hmac.New(sha256.New, []byte("super-secret-key-that-is-at-least-32-chars-long-for-testing!"))
	mac.Write([]byte(header + "." + payload))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return fmt.Sprintf("%s.%s.%s", header, payload, sig)
}

func setupTestMux(t *testing.T) *http.ServeMux {
	tmpDir := t.TempDir()
	tmpDB := tmpDir + "/test.db"

	sqlStore, err := sqlite.NewSQLiteStore(tmpDB)
	if err != nil {
		t.Fatalf("failed to init test sqlite store: %v", err)
	}

	mediaStore, err := storage.NewLocalMediaStore(tmpDir + "/media")
	if err != nil {
		t.Fatalf("failed to init test media store: %v", err)
	}

	gpuDetector := gpu.NewDetector()
	pyWorker := worker.NewPythonWorker()
	svc := application.NewTrainingService(sqlStore, sqlStore, sqlStore, sqlStore, pyWorker, gpuDetector)
	mediaSvc := application.NewMediaService(mediaStore)

	handler := primaryHttp.NewTrainingHandler(svc)
	bmkHandler := primaryHttp.NewBenchmarkHandler(svc)
	mediaHandler := primaryHttp.NewMediaHandler(mediaSvc)
	mux := http.NewServeMux()
	primaryHttp.RegisterRoutes(mux, handler, bmkHandler, mediaHandler)

	return mux
}

func TestUnauthenticatedRequestBlocked(t *testing.T) {
	mux := setupTestMux(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/training/jobs", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 Unauthorized, got %d", w.Code)
	}
}

func TestAuthenticatedJobsList(t *testing.T) {
	mux := setupTestMux(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/training/jobs", nil)
	req.Header.Set("Authorization", "Bearer "+generateTestToken("tenant_alpha", "viewer"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", w.Code)
	}
}

func TestZipSlipBlockedOnImport(t *testing.T) {
	mux := setupTestMux(t)

	// Create a zip with malicious traversal path
	var zipBuf bytes.Buffer
	zw := zip.NewWriter(&zipBuf)
	f, err := zw.Create("../../evil.txt")
	if err != nil {
		t.Fatalf("failed to create zip entry: %v", err)
	}
	f.Write([]byte("malicious payload"))
	zw.Close()

	// Form body
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	fw, err := mw.CreateFormFile("files", "malicious_dataset.zip")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	fw.Write(zipBuf.Bytes())
	mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/training/datasets/import", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+generateTestToken("tenant_alpha", "admin"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if _, err := os.Stat("../../evil.txt"); err == nil {
		os.Remove("../../evil.txt")
		t.Fatalf("SECURITY VIOLATION: Zip Slip file was written to host filesystem!")
	}
}

func TestRegisterPathSensitiveDirectoryBlocked(t *testing.T) {
	mux := setupTestMux(t)

	body := `{"path":"/etc","dataset_id":"system_etc"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/training/datasets/register-path", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+generateTestToken("tenant_alpha", "admin"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden on registering /etc, got %d", w.Code)
	}
}

func TestMediaFolderCreateAndUpload(t *testing.T) {
	mux := setupTestMux(t)

	// 1. Create folder "carros"
	body := `{"name":"carros"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/media/folders", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+generateTestToken("tenant_alpha", "admin"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on media folder create, got %d", w.Code)
	}

	// 2. Upload video file to folder "carros"
	var uploadBody bytes.Buffer
	mw := multipart.NewWriter(&uploadBody)
	_ = mw.WriteField("folder", "carros")
	fw, _ := mw.CreateFormFile("file", "test_video.mp4")
	fw.Write([]byte("fake-mp4-video-content-stream"))
	mw.Close()

	uploadReq := httptest.NewRequest(http.MethodPost, "/api/v1/media/upload", &uploadBody)
	uploadReq.Header.Set("Content-Type", mw.FormDataContentType())
	uploadReq.Header.Set("Authorization", "Bearer "+generateTestToken("tenant_alpha", "admin"))
	wUpload := httptest.NewRecorder()
	mux.ServeHTTP(wUpload, uploadReq)

	if wUpload.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on media upload, got %d: %s", wUpload.Code, wUpload.Body.String())
	}

	// 3. List media sources
	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/media/sources", nil)
	listReq.Header.Set("Authorization", "Bearer "+generateTestToken("tenant_alpha", "viewer"))
	wList := httptest.NewRecorder()
	mux.ServeHTTP(wList, listReq)

	if wList.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on media list, got %d", wList.Code)
	}

	if !bytes.Contains(wList.Body.Bytes(), []byte("test_video.mp4")) {
		t.Fatalf("expected list response to contain test_video.mp4, got %s", wList.Body.String())
	}
}
