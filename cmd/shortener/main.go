package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sashganush/shortcut/internal/config"
	handlers "github.com/sashganush/shortcut/internal/handlers"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"net/url"
	"os"
)

func NewRouter() chi.Router {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	handlers.LoadAllToStorage()

	r.Route("/", func(r chi.Router) {
		r.Get("/ping", handlers.Ping)
		r.Post("/", handlers.PostRequestHandler)
		r.Get(handlers.GetBaseUri()+"{ID}", handlers.GetRequestHandler)
		r.Post("/api/shorten", handlers.PostRequestApiHandler)
	})

	return r
}

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "shortener",
		Short: "Shortener Applications",
		Long:  `Shortenter service create a simple redirect service for long urls.`,
	}
)

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return nil
}

var flagServerAddress string
var flagServerPort string
var flagBaseUrl string
var flagFileStoragePath string

func init() {
	rootCmd.PersistentFlags().StringVarP(&flagServerAddress, "serveraddress","a", "", "service address the same from SERVER_ADDRESS")
	rootCmd.PersistentFlags().StringVarP(&flagServerPort, "serverport","p", "", "service port the same from SERVER_PORT")
	rootCmd.PersistentFlags().StringVarP(&flagBaseUrl, "baseurl","b", "", "base url the same like BASE_URL env")
	rootCmd.PersistentFlags().StringVarP(&flagFileStoragePath, "filestoragepath","f", "", "file storage path the same like FILE_STORAGE_PATH env")
}

func main() {

	if err := env.Parse(&config.Cfg); err != nil {
		fmt.Println("failed:", err)
		panic(err)
	}

	_, err := url.Parse(config.Cfg.BaseUrl)
	if err != nil {
		panic(err)
	}

	Execute()

	if flagServerAddress != "" {
		config.Cfg.ServerAddress = flagServerAddress
	}

	if flagBaseUrl != "" {
		config.Cfg.BaseUrl = flagBaseUrl
	}

	if flagFileStoragePath != "" {
		config.Cfg.FileStoragePath = flagFileStoragePath
	}

	fmt.Printf("Flag Server Address at %s\n", flagServerAddress)
	fmt.Printf("Flag base_url at %s\n", flagBaseUrl)
	fmt.Printf("Flag filestoragepath at %s\n", flagFileStoragePath)

	fmt.Printf("Cfg Server Address at %s\n", config.Cfg.ServerAddress)
	fmt.Printf("Cfg base_url at %s\n", config.Cfg.BaseUrl)
	fmt.Printf("Cfg filestoragepath at %s\n", config.Cfg.FileStoragePath)

	r := NewRouter()
	log.Fatal(http.ListenAndServe(config.Cfg.ServerAddress, r), nil)
}
