package handlers_test

import (
	"github.com/ansedo/url-shortener/internal/app/handlers"
	"github.com/ansedo/url-shortener/internal/app/shortener"
	"github.com/ansedo/url-shortener/internal/app/storage/memory"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecodeURL(t *testing.T) {
	type want struct {
		statusCode int
		body       string
		location   string
	}

	testRouter := chi.NewRouter()
	testStorage := memory.NewStorage()
	testData := map[string]string{"short-ya": "https://ya.ru", "short-google": "https://google.com"}
	for key, value := range testData {
		err := testStorage.Set(key, value)
		require.NoError(t, err)
	}
	testShortener := shortener.NewShortener(testStorage)
	testRouter.Get("/{id}", handlers.DecodeURL(testShortener))

	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Decode: " + testData["short-ya"],
			url:  "/short-ya",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   testData["short-ya"],
			},
		},
		{
			name: "Decode: empty key",
			url:  "/",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "Decode: wrong key",
			url:  "/42",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Decode: " + testData["short-google"],
			url:  "/short-google",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   testData["short-google"],
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(testRouter)

			testClient := testServer.Client()
			testClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}

			response, err := testClient.Get(testServer.URL + tt.url)
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, response.StatusCode)

			_, err = io.ReadAll(response.Body)
			defer response.Body.Close()
			assert.NoError(t, err)

			if tt.want.location != "" {
				assert.Equal(t, tt.want.location, response.Header.Get("Location"))
			}
		})
	}
}
