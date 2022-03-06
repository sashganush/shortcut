package main

import (
	"bytes"
	"encoding/json"
	"github.com/sashganush/shortcut/internal/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string) (*http.Response, string) {
	ioBody := bytes.NewReader([]byte(body))
	req, err := http.NewRequest(method, ts.URL+path, ioBody)
	require.NoError(t, err)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	resp.Body.Close()
	require.NoError(t, err)

	return resp, string(respBody)
}


func TestRouter(t *testing.T) {
	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/ping","")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "pong", body)

	resp, _ = testRequest(t, ts, "GET", "/123","")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, _ = testRequest(t, ts, "GET", "/","")
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	resp, body = testRequest(t, ts, "POST", "/", "http://www.ya.ru/1")
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	u, _ := url.Parse(body)
	uri := u.RequestURI()

	resp, _ = testRequest(t, ts, "GET", uri, "")
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)

	resp, body = testRequest(t, ts, "POST", "/api/shorten", "{\"url\":\"http://www.ya.ru/2\"}")
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t,"application/json", resp.Header.Get("Content-Type"))

	var responseJson handlers.ResponseJson
	err := json.Unmarshal([]byte(body), &responseJson)
	require.NoError(t, err)

	u, _ = url.Parse(responseJson.Result)
	uri = u.RequestURI()

	resp, _ = testRequest(t, ts, "GET", uri, "")
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)

}
