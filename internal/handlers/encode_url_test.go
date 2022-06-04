package handlers_test

import (
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/handlers"
	"github.com/ansedo/url-shortener/internal/services/shortener"
	"github.com/ansedo/url-shortener/internal/storage/memory"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEncodeURL(t *testing.T) {
	type want struct {
		statusCode int
		body       string
		location   string
	}

	testRouter := chi.NewRouter()
	testShortener := shortener.NewShortener(memory.NewStorage())
	testRouter.Post("/", handlers.EncodeURL(testShortener))

	tests := []struct {
		name string
		url  string
		body string
		want want
	}{
		{
			name: "Encode: https://ya.ru",
			url:  "/",
			body: "https://ya.ru",
			want: want{
				statusCode: http.StatusCreated,
				body:       config.SiteAddress + "/0",
			},
		},
		{
			name: "Encode: wrong url",
			url:  "/",
			body: "://this url does not exist:8080",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       config.RequestNotAllowedError,
			},
		},
		{
			name: "Encode: empty url",
			url:  "/",
			body: "",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       config.RequestNotAllowedError,
			},
		},
		{
			name: "Encode: https://google.com",
			url:  "/",
			body: "https://google.com",
			want: want{
				statusCode: http.StatusCreated,
				body:       config.SiteAddress + "/1",
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

			response, err := testClient.Post(
				testServer.URL+tt.url,
				"text/plain; charset=utf8",
				strings.NewReader(tt.body),
			)
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, response.StatusCode)

			body, err := io.ReadAll(response.Body)
			defer response.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.body, strings.TrimSpace(string(body)))
		})
	}
}
