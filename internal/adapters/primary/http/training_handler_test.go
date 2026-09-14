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
	tmpDB := fmt.Sprintf("/tmp/test_hydraforge_%d.db", time.Now().UnixNano())
	t.Cleanup(func() { os.Remove(tmpDB) })

	sqlStore, err := sqlite.NewSQLiteStore(tmpDB)
	if err != nil {
		t.Fatalf("failed to init test sqlite store: %v", err)
	}

	gpuDetector := gpu.NewDetector()
	pyWorker := worker.NewPythonWorker()
	svc := application.NewTrainingService(sqlStore, sqlStore, sqlStore, sqlStore, pyWorker, gpuDetector)

	handler := primaryHttp.NewTrainingHandler(svc)
	bmkHandler := primaryHttp.NewBenchmarkHandler(svc)
	mux := http.NewServeMux()
	primaryHttp.RegisterRoutes(mux, handler, bmkHandler)

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

	if _, err := os.Stat("/home/hades/evil.txt"); err == nil {
		os.Remove("/home/hades/evil.txt")
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
