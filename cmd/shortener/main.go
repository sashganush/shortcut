package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sashganush/shortcut/internal/config"
	handlers "github.com/sashganush/shortcut/internal/handlers"
	"log"
	"net/http"
	"net/url"
)

var Cfg config.Config

func NewRouter() chi.Router {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Route("/", func(r chi.Router) {
		r.Get("/ping", handlers.Ping)
		r.Post("/", handlers.PostRequestHandler)
		r.Get(handlers.GetBaseUri()+"{ID}", handlers.GetRequestHandler)
		r.Post("/api/shorten", handlers.PostRequestApiHandler)
	})

	return r
}


func main() {

	if err := env.Parse(&Cfg); err != nil {
		fmt.Println("failed:", err)
	}

	fmt.Printf("Start at %s%s\n", Cfg.ServerAddress,Cfg.ServerPort)
	fmt.Printf("Start with base url %s\n", Cfg.BaseUrl)

	u, err := url.Parse(Cfg.BaseUrl)
	if err != nil {
		panic(err)
	}

	r := NewRouter()
	log.Fatal(http.ListenAndServe(u.Host, r), nil)
}
