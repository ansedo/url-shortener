package shortener_test

import (
	"github.com/ansedo/url-shortener/internal/app/config"
	"github.com/ansedo/url-shortener/internal/app/shortener"
	"github.com/ansedo/url-shortener/internal/app/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortener(t *testing.T) {
	type want struct {
		statusCode int
		body       string
		location   string
	}

	tests := []struct {
		name   string
		url    string
		body   string
		method string
		want   want
	}{
		{
			name:   "Encode: https://ya.ru",
			url:    "/",
			body:   "https://ya.ru",
			method: http.MethodPost,
			want: want{
				statusCode: http.StatusCreated,
				body:       config.SiteAddress + "/0",
			},
		},
		{
			name:   "Encode: wrong url",
			url:    "/",
			body:   "://this url does not exist:8080",
			method: http.MethodPost,
			want: want{
				statusCode: http.StatusBadRequest,
				body:       shortener.RequestNotAllowedError,
			},
		},
		{
			name:   "Decode: https://ya.ru",
			url:    "/0",
			method: http.MethodGet,
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "https://ya.ru",
			},
		},
		{
			name:   "Encode: empty url",
			url:    "/",
			method: http.MethodPost,
			want: want{
				statusCode: http.StatusBadRequest,
				body:       shortener.RequestNotAllowedError,
			},
		},
		{
			name:   "Encode: https://google.com",
			url:    "/",
			body:   "https://google.com",
			method: http.MethodPost,
			want: want{
				statusCode: http.StatusCreated,
				body:       config.SiteAddress + "/1",
			},
		},
		{
			name:   "Decode: empty key",
			url:    "/",
			method: http.MethodGet,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:   "Decode: wrong key",
			url:    "/42",
			method: http.MethodGet,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:   "Decode: https://google.com",
			url:    "/1",
			method: http.MethodGet,
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "https://google.com",
			},
		},
	}

	testShortener := shortener.NewShortener(memory.NewStorage())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			writer := httptest.NewRecorder()
			http.HandlerFunc(testShortener.ServeHTTP).ServeHTTP(writer, request)
			response := writer.Result()
			assert.Equal(t, tt.want.statusCode, response.StatusCode)

			body, err := io.ReadAll(response.Body)
			defer response.Body.Close()
			require.NoError(t, err)

			if tt.want.body != "" {
				assert.Equal(t, tt.want.body, strings.TrimSpace(string(body)))
			}

			if tt.want.location != "" {
				assert.Equal(t, tt.want.location, response.Header.Get("Location"))
			}
		})
	}
}
