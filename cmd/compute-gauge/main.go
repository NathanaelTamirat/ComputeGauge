package main

import (
	"compute-gauge/pkg/handlers"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		path = strings.TrimSuffix(path, "/")
		if strings.HasPrefix(path, "/static/") {
			handlers.HandleStatic(w, r, "/static/")
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
	})
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
