package handlers

import (
	"github.com/ansedo/url-shortener/internal/services/shortener"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func DecodeURL(s *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uri, err := s.Storage.Get(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, uri, http.StatusTemporaryRedirect)
	}
}
