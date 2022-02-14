package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	internal "github.com/sashganush/shortcut/internal/app/shortener"
	"net/http"
)

func main() {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		internal.PostRequestHandler(w,r)
	})

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		internal.GetRequestHandler(w,r)
	})

	http.ListenAndServe(":8080", r)
}
