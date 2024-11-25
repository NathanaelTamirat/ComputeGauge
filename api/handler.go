package api

import (
	"compute-gauge/handlers"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
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
