package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

const ShortURLLen = 10
const DefaultShema = "http://"
const DefaultHostName = "localhost"
const DefaultHostPort = ":8080"

var allRedirects = map[string]string{}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

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

func Ping(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func PostRequestHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	oldURL, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newURL := RandStringRunes(ShortURLLen)
	allRedirects[newURL] = string(oldURL)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s%s%s/%s", DefaultShema, DefaultHostName, DefaultHostPort, newURL)
}

func GetRequestHandler(w http.ResponseWriter, r *http.Request) {

	if s, exists := allRedirects[chi.URLParam(r, "ID")]; exists {

		w.Header().Set("Location", s)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	http.Error(w, "Unknown redirect", http.StatusBadRequest)
}

type RequestJson struct {
	ID  int    `json:"-"`
	Url string `json:"url"`
}

type ResponseJson struct {
	ID  int       `json:"-"`
	Result string `json:"result"`
}

func PostRequestApiHandler(w http.ResponseWriter, r *http.Request) {

	var tmpRequest RequestJson
	var tmpResponse ResponseJson

	defer r.Body.Close()
	oldURL, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(oldURL, &tmpRequest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newURL := RandStringRunes(ShortURLLen)
	allRedirects[newURL] = tmpRequest.Url

	tmpResponse.Result = "http://"+r.Host+"/"+newURL

	ret, err := json.Marshal(tmpResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(ret)
}