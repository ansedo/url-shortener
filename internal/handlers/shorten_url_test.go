package handlers_test

import (
	"context"
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

func TestShortenURL(t *testing.T) {
	type want struct {
		statusCode int
		body       string
		location   string
	}

	svc := shortener.New(context.Background())
	r := chi.NewRouter()
	r.Post("/", handlers.ShortenURL(svc))

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
				body:       svc.BaseURL + "/VRb8948o",
			},
		},
		{
			name: "Encode: wrong url",
			url:  "/",
			body: "://this url does not exist:8080",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       handlers.ErrRequestNotAllowed.Error(),
			},
		},
		{
			name: "Encode: empty url",
			url:  "/",
			body: "",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       handlers.ErrRequestNotAllowed.Error(),
			},
		},
		{
			name: "Encode: https://google.com",
			url:  "/",
			body: "https://google.com",
			want: want{
				statusCode: http.StatusCreated,
				body:       svc.BaseURL + "/65lvAYxL",
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

			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.body, strings.TrimSpace(string(body)))
		})
	}
}
