package handlers_test

import (
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/handlers"
	"github.com/ansedo/url-shortener/internal/services/shortener"
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

	r := chi.NewRouter()
	r.Post("/", handlers.EncodeURL(shortener.New()))
	cfg := config.New()

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
				body:       cfg.SiteAddress + "/0",
			},
		},
		{
			name: "Encode: wrong url",
			url:  "/",
			body: "://this url does not exist:8080",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       cfg.RequestNotAllowedError,
			},
		},
		{
			name: "Encode: empty url",
			url:  "/",
			body: "",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       cfg.RequestNotAllowedError,
			},
		},
		{
			name: "Encode: https://google.com",
			url:  "/",
			body: "https://google.com",
			want: want{
				statusCode: http.StatusCreated,
				body:       cfg.SiteAddress + "/1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := httptest.NewServer(r)

			c := s.Client()
			c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}

			resp, err := c.Post(
				s.URL+tt.url,
				"text/plain; charset=utf8",
				strings.NewReader(tt.body),
			)
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.body, strings.TrimSpace(string(body)))
		})
	}
}
