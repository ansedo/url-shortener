package handlers

import (
	"github.com/ansedo/url-shortener/internal/services/shortener"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/url"
)

func GetOriginalURL(s *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		if id == "" {
			http.Error(w, ErrRequestNotAllowed.Error(), http.StatusBadRequest)
			return
		}

		storageURL, err := s.Storage.GetByShortURL(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		validURL, err := url.ParseRequestURI(storageURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, validURL.String(), http.StatusTemporaryRedirect)
	}
}
