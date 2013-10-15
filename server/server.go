package main

import (
	"fmt"
	"github.com/jim/monk"
	"log"
	"net/http"
)

func main() {
	cache := &monk.LocalCache{}
	http.Handle("/assets/", http.StripPrefix("/assets/", logRequest(http.FileServer(cache))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type statusResponseWriter struct {
	status int
	http.ResponseWriter
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &statusResponseWriter{-1, w}
		h.ServeHTTP(writer, r)
		fmt.Printf("GET %s %v\n", r.URL.Path, writer.status)
	})
}
