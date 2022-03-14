package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sashganush/shortcut/internal/config"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const ShortURLLen = 10
const DefaultSchema = "http://"

var AllRedirects = map[string]string{}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var v = os.Getenv("BASE_URL")

func GetBaseUri() string {
	return "/"+os.Getenv("BASE_URI")
}


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

	s := RandStringRunes(ShortURLLen)
	SaveToStorage(string(oldURL),s)
	newURL := GetBaseUri()+s
	AllRedirects[newURL] = string(oldURL)


	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s%s%s" ,DefaultSchema , r.Host, newURL)
}

func GetRequestHandler(w http.ResponseWriter, r *http.Request) {

	if s, exists := AllRedirects[GetBaseUri()+chi.URLParam(r, "ID")]; exists {
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
	ID     int    `json:"-"`
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

	s := RandStringRunes(ShortURLLen)
	newURL := GetBaseUri()+s
	SaveToStorage(tmpRequest.Url,s)
	AllRedirects[newURL] = tmpRequest.Url

	tmpResponse.Result =  DefaultSchema+r.Host+newURL

	ret, err := json.Marshal(tmpResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(ret)
}

// TODO: move to storage internal

type Redirects struct {
	ID       uint    `json:"id"`
	ShortUrl string  `json:"short_url"`
	LongUrl  string  `json:"long_url"`
}
type producer struct {
	file    *os.File
	encoder *json.Encoder
}
func NewProducer(fileName string) (*producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}
func (p *producer) WriteRedirect(redirect *Redirects) error {
	return p.encoder.Encode(&redirect)
}
func (p *producer) Close() error {
	return p.file.Close()
}


func SaveToStorage(l string, s string) error {

	var redirects = []*Redirects{
		{
			ShortUrl: s,
			LongUrl:  l,
		},
	}

	producer, err := NewProducer(config.Cfg.FileStoragePath)
	if err != nil {
		fmt.Println("Error open file:", err.Error())
		return err
	}
	defer producer.Close()

	for _, redirect := range redirects {
		if err := producer.WriteRedirect(redirect); err != nil {
			fmt.Println("Error write file:", err.Error())
			return err
		}
	}

	return nil
}

type consumer struct {
	file    *os.File
	decoder *json.Decoder
}
func NewConsumer(fileName string) (*consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}
func (c *consumer) ReadRedirect() (*Redirects, error) {
	redirect := &Redirects{}
	if err := c.decoder.Decode(&redirect); err != nil {
		return nil, err
	}
	return redirect, nil
}
func (c *consumer) Close() error {
	return c.file.Close()
}

func LoadAllToStorage() error {

	consumer, err := NewConsumer(config.Cfg.FileStoragePath)
	fmt.Println("Load from file:",config.Cfg.FileStoragePath)
	if err != nil {
		fmt.Println("Error open file:",err.Error())
		log.Fatal(err)
	}
	defer consumer.Close()

	i := 0
	for true {
		readedRedirect, err := consumer.ReadRedirect()
		if err != nil {
			fmt.Println("Read from storage",i,"redirects")
			break
		}
		i++
//		fmt.Println(readedRedirect,readedRedirect.ShortUrl,readedRedirect.LongUrl)
		AllRedirects[GetBaseUri()+readedRedirect.ShortUrl] = readedRedirect.LongUrl
//		fmt.Println(AllRedirects[GetBaseUri()+readedRedirect.ShortUrl])
	}
  return nil
}
