package main

import "net/http"

func main() {
	handler := http.NewServeMux()

	handler.Handle("/", http.FileServer(http.Dir(".")))
	handler.Handle("/assets", http.FileServer(http.Dir("./assets")))

	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
