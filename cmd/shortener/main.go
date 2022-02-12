package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const SHORTURLLEN = 8

var allRedirects = map[string]string{
	"4444": "https://google.com/abcd",
	"1111": "https://ya.ru/1",
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

func PostRequestHandler(w http.ResponseWriter, r *http.Request) (s int) {

	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return http.StatusInternalServerError
	}

	n := RandStringRunes(SHORTURLLEN)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "https://tinyurl.com/%s", n)
	return http.StatusCreated
}

func GetRequestHandler(w http.ResponseWriter, r *http.Request) (s int) {

	if s, exists := allRedirects[strings.TrimPrefix(r.URL.Path, "/")]; exists {
		http.Redirect(w, r, s, http.StatusTemporaryRedirect)
		return http.StatusTemporaryRedirect
	}
	http.Error(w, "Unknown redirect", http.StatusBadRequest)
	return http.StatusBadRequest
}

func MainHandler(w http.ResponseWriter, r *http.Request) {

	var s int

	switch r.Method {
	case "POST":
		s = PostRequestHandler(w, r)
	default:
		s = GetRequestHandler(w, r)
	}

	fmt.Println(time.Now().Format(time.RFC3339), s, r.URL.Path)
}

func main() {

	http.HandleFunc("/", MainHandler)
	http.ListenAndServe("localhost:8080", nil)
}
