package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const SHORTURLLEN = 10

var allRedirects = map[string]string{
	"4444": "https://google.com/abcd",
	"1111": "https://ya.ru/1",
	"kyQFSSqy": "http://t9xiost0mawj62.ru",
	"ytOkkmjo": "http://y4qcfp8ur.yandex",

}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func PostRequestHandler(w http.ResponseWriter, r *http.Request)  {

	oldURL, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newURL := RandStringRunes(SHORTURLLEN)
	allRedirects[newURL] = string(oldURL[:])

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "http://localhost:8080/%s", newURL)
}

func GetRequestHandler(w http.ResponseWriter, r *http.Request)  {

	if s, exists := allRedirects[strings.TrimPrefix(r.URL.Path, "/")]; exists {

		w.Header().Set("Location", s)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	http.Error(w, "Unknown redirect", http.StatusBadRequest)
}

type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		// handle returned error here.
		w.WriteHeader(503)
		w.Write([]byte("bad"))
	}
}

func customHandler(w http.ResponseWriter, r *http.Request) error {
	GetRequestHandler(w,r)
	return nil
}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Method("GET", "/*", Handler(customHandler))

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		PostRequestHandler(w,r)
	})

	http.ListenAndServe(":8080", r)
}
