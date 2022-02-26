package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	handlers "github.com/sashganush/shortcut/internal/handlers"
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
		r.Get("/{ID}", handlers.GetRequestHandler)
	})

	return r
}

func main() {

    r := NewRouter()
	log.Fatal(http.ListenAndServe(handlers.DefaultHostName+handlers.DefaultHostPort, r), nil)
}
