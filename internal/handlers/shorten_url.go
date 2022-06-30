package handlers

import (
	"fmt"
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/services/shortener"
	"io"
	"net/http"
	"net/url"
)

func ShortenURL(s *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		uri, err := url.ParseRequestURI(string(body))
		if err != nil {
			http.Error(w, ErrRequestNotAllowed.Error(), http.StatusBadRequest)
			return
		}

		id, err := s.GenerateID(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.Storage.Add(r.Context(), id, uri.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err = fmt.Fprintf(w, config.Get().BaseURL+"/"+id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}