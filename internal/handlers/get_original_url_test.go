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
	"testing"
)

func TestGetOriginalURL(t *testing.T) {
	type want struct {
		statusCode int
		body       string
		location   string
	}

	ctx := context.Background()
	svc := shortener.New(ctx)
	data := map[string]string{"short-ya": "https://ya.ru", "short-google": "https://google.com"}
	for key, value := range data {
		err := svc.Storage.Add(ctx, key, value)
		require.NoError(t, err)
	}

	r := chi.NewRouter()
	r.Get("/{id}", handlers.GetOriginalURL(svc))

	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Decode: " + data["short-ya"],
			url:  "/short-ya",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   data["short-ya"],
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
			name: "Decode: " + data["short-google"],
			url:  "/short-google",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   data["short-google"],
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

			resp, err := c.Get(s.URL + tt.url)
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			_, err = io.ReadAll(resp.Body)
			defer resp.Body.Close()
			assert.NoError(t, err)

			if tt.want.location != "" {
				assert.Equal(t, tt.want.location, resp.Header.Get("Location"))
			}
		})
	}
}
