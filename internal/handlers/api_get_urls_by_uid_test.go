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

func TestAPIGetURLsByUIDFetchURLs(t *testing.T) {
	urls := []string{`https://ya.ru`, `https://google.com`}
	contentType := "text/plain; charset=utf8"

	svc := shortener.New()
	r := chi.NewRouter()
	r.Post("/", handlers.ShortenURL(svc))
	r.Get("/api/user/urls", handlers.APIGetURLsByUID(svc))
	s := httptest.NewServer(r)
	c := s.Client()

	t.Run("shorten first url", func(t *testing.T) {
		resp, err := c.Post(s.URL, contentType, strings.NewReader(urls[0]))
		require.NoError(t, err)

		resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("fetch url", func(t *testing.T) {
		resp, err := c.Get(s.URL + "/api/user/urls")
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.NoError(t, err)
		assert.Contains(t, string(body), urls[0])
		assert.Contains(t, string(body), config.Get().BaseURL)
	})

	t.Run("shorten second url", func(t *testing.T) {
		resp, err := c.Post(s.URL, contentType, strings.NewReader(urls[1]))
		require.NoError(t, err)

		resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("fetch urls", func(t *testing.T) {
		resp, err := c.Get(s.URL + "/api/user/urls")
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.NoError(t, err)
		assert.Contains(t, string(body), urls[0])
		assert.Contains(t, string(body), urls[1])
		assert.Contains(t, string(body), config.Get().BaseURL)
	})
}
func TestAPIGetURLsByUIDFetchEmptyURLs(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/api/user/urls", handlers.APIGetURLsByUID(shortener.New()))
	s := httptest.NewServer(r)
	c := s.Client()

	t.Run("fetch empty urls", func(t *testing.T) {
		resp, err := c.Get(s.URL + "/api/user/urls")
		require.NoError(t, err)

		_, err = io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}
