package http

import (
	"fmt"
	"net/http"
)

// OpenAPI3Spec is the OpenAPI 3.0 specification for HydraForge.
const OpenAPI3Spec = `{
  "openapi": "3.0.0",
  "info": {
    "title": "HydraForge - YOLO AI Training Studio REST API",
    "version": "1.0.0",
    "description": "High-Performance Control Plane API for YOLOv8, YOLO11, and YOLO26 Training, Telemetry, and TensorRT Export."
  },
  "paths": {
    "/api/v1/training/jobs": {
      "get": {
        "summary": "List training jobs",
        "responses": { "200": { "description": "Array of training jobs" } }
      },
      "post": {
        "summary": "Launch a new YOLO training run",
        "responses": { "201": { "description": "Training job initialized" } }
      }
    },
    "/api/v1/training/datasets": {
      "get": { "summary": "List registered datasets", "responses": { "200": { "description": "Datasets array" } } },
      "post": { "summary": "Register a new dataset", "responses": { "201": { "description": "Dataset saved" } } }
    },
    "/api/v1/training/telemetry": {
      "get": { "summary": "Get live GPU and training metrics telemetry", "responses": { "200": { "description": "Telemetry snapshot" } } }
    },
    "/api/v1/training/export": {
      "post": { "summary": "Export model to TensorRT or ONNX", "responses": { "200": { "description": "Exported artifact" } } }
    }
  }
}`

// ServeSwaggerUI renders the dark-themed Swagger UI.
func ServeSwaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>HydraForge API Docs - Swagger UI</title>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.18.3/swagger-ui.css" />
  <style>body { margin: 0; background: #07080c; color: #fff; } .swagger-ui { filter: invert(88%%) hue-rotate(180deg); }</style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.18.3/swagger-ui-bundle.js"></script>
  <script>SwaggerUIBundle({ url: '/swagger/doc.json', dom_id: '#swagger-ui' });</script>
</body>
</html>`)
	w.Write([]byte(html))
}
