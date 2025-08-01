package swagger

import (
	"embed"
	"net/http"
)

//go:embed swagger-ui.html swagger.yaml
var swaggerFiles embed.FS

// SwaggerUIController serves the embedded Swagger UI HTML
func SwaggerUIController(w http.ResponseWriter, r *http.Request) {
	data, err := swaggerFiles.ReadFile("swagger-ui.html")
	if err != nil {
		http.Error(w, "Swagger UI not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

// SwaggerYAMLController serves the embedded Swagger YAML file
func SwaggerYAMLController(w http.ResponseWriter, r *http.Request) {
	data, err := swaggerFiles.ReadFile("swagger.yaml")
	if err != nil {
		http.Error(w, "Swagger YAML not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/x-yaml; charset=utf-8")
	w.Write(data)
}
