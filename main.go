package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.metricsMiddleware(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.Handle("/app/assets", http.StripPrefix("/app/assets", http.FileServer(http.Dir("./assets"))))

	// API
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.validateChirp)

	// ADMIN
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("Server running on :8080")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
