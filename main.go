package main

import (
	"fmt"
	"net/http"
)

func main() {
	handler := http.NewServeMux()

	handler.Handle("/app", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	handler.Handle("/app/assets", http.StripPrefix("/app/assets", http.FileServer(http.Dir("./assets"))))
	handler.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	fmt.Println("Server running on :8080")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
