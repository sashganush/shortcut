package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	handlers "github.com/sashganush/shortcut/internal/handlers"
	storage "github.com/sashganush/shortcut/internal/storage"
	"log"
	"net/http"
)

func NewRouter() chi.Router {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Route("/", func(r chi.Router) {
		r.Get("/ping", handlers.Ping)
		r.Post("/", handlers.PostRequestHandler)

		r.Route("/{ID}", func(r chi.Router) {
			r.Get("/", handlers.GetRequestHandler)
		})
	})

	return r
}

func main() {

    r := NewRouter()
	log.Fatal(http.ListenAndServe(storage.DefaultHostName+storage.DefaultHostPort, r), nil)
}
