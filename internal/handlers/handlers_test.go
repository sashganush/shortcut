package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler(t *testing.T) {
	type want struct {
		codeGet  int
		response string
	}
	tests := []struct {
		name    string
		want    want
		request string
	}{
		{
			name: "positive test #1",
			want: want{
				codeGet:  200,
				response: `pong`,
			},
			request: "/ping",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(Ping)
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.codeGet, result.StatusCode)

			GetBody, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			require.NoError(t, err)
			assert.Equal(t, string(GetBody), tt.want.response)
		})
	}
}
