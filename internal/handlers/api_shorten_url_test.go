package handlers_test

import (
	"context"
	"fmt"
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

func TestAPIShortenURL(t *testing.T) {
	type want struct {
		statusCode int
		body       string
	}

	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "Encode from JSON: https://ya.ru",
			body: `{"url":"https://ya.ru"}`,
			want: want{
				statusCode: http.StatusCreated,
				body:       fmt.Sprintf(`{"result":"%s/VRb8948o"}`, config.Get().BaseURL),
			},
		},
		{
			name: "Encode from JSON: wrong url",
			body: `{"url":"--///this url does not_exist%8080"}`,
			want: want{
				statusCode: http.StatusBadRequest,
				body:       fmt.Sprintf(`{"error":"%s"}`, handlers.ErrRequestNotAllowed.Error()),
			},
		},
		{
			name: "Encode from JSON: empty url",
			body: `{"url": ""}`,
			want: want{
				statusCode: http.StatusBadRequest,
				body:       fmt.Sprintf(`{"error":"%s"}`, handlers.ErrRequestNotAllowed.Error()),
			},
		},
		{
			name: "Encode from JSON: https://google.com",
			body: `{"url":"https://google.com"}`,
			want: want{
				statusCode: http.StatusCreated,
				body:       fmt.Sprintf(`{"result":"%s/65lvAYxL"}`, config.Get().BaseURL),
			},
		},
	}

	r := chi.NewRouter()
	r.Post("/", handlers.APIShortenURL(shortener.New(context.Background())))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := httptest.NewServer(r)
			c := s.Client()
			c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}

			resp, err := c.Post(s.URL, "application/json; charset=utf-8", strings.NewReader(tt.body))
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			require.NoError(t, err)
			assert.JSONEq(t, tt.want.body, strings.TrimSpace(string(body)))
		})
	}
}
