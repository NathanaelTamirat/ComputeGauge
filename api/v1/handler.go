package api

import (
	"compute-gauge/pkg/handlers"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	path = strings.TrimSuffix(path, "/")
	if strings.HasPrefix(path, "/static/") {
		handlers.HandleStatic(w, r, "/")
		return
	}
	switch path {
	case "", "/index", "/index.html":
		handlers.HandleIndex(w, r)
	case "/api/calculate":
		handlers.HandleCalculate(w, r)
	case "/documentation":
		handlers.HandleDocs(w, r)
	default:
		http.NotFound(w, r)
	}
}
