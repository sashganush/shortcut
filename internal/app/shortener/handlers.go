package internal

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
