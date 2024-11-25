package handlers

import (
	"compute-gauge/config"
	"compute-gauge/memory"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	models, err := config.LoadModelConfigs()
	if err != nil {
		log.Printf("Error loading models: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles(filepath.Join("templates", "index.html"))
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	dataTypes := make([]string, 0, len(config.DataTypeSizes))
	for dtype := range config.DataTypeSizes {
		dataTypes = append(dataTypes, dtype)
	}
	data := memory.PageData{
		Models:    models,
		DataTypes: dataTypes,
	}
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func HandleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req memory.MemoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, fmt.Sprintf("Invalid request format: %v", err), http.StatusBadRequest)
		return
	}
	result, err := memory.CalculateMemoryRequirements(&req)
	if err != nil {
		log.Printf("Error calculating memory: %v", err)
		http.Error(w, fmt.Sprintf("Error calculating memory: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// HandleDocs serves the documentation markdown file
func HandleDocs(w http.ResponseWriter, r *http.Request) {
	docPath := filepath.Join("docs", "documentation.md")
	docContent, err := os.ReadFile(docPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading documentation: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/markdown")
	w.Write(docContent)
}

// HandleStatic serves static files
func HandleStatic(w http.ResponseWriter, r *http.Request, prefix string) {
	fs := http.FileServer(http.Dir("."))
	http.StripPrefix(prefix, fs).ServeHTTP(w, r)
}
