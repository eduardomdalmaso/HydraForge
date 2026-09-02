package http

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// HydraStreamProxy forwards requests from /api/v1/hydrastream/ to HydraStream Control Plane.
func HydraStreamProxy(targetURL string) http.HandlerFunc {
	target, err := url.Parse(targetURL)
	if err != nil {
		target, _ = url.Parse("http://localhost:8080")
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director

	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// Strip the /api/v1/hydrastream prefix
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/api/v1/hydrastream")
		if !strings.HasPrefix(req.URL.Path, "/") {
			req.URL.Path = "/" + req.URL.Path
		}
		req.Host = target.Host
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		proxy.ServeHTTP(w, r)
	}
}
