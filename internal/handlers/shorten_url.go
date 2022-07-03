package handlers

import (
	"errors"
	"fmt"
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/services/shortener"
	"github.com/ansedo/url-shortener/internal/storages"
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
		originalURL := string(body)

		if _, err = url.ParseRequestURI(originalURL); err != nil {
			http.Error(w, ErrRequestNotAllowed.Error(), http.StatusBadRequest)
			return
		}

		shortURLID, err := s.GenerateID(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = s.Storage.Add(r.Context(), shortURLID, originalURL); err != nil {
			if errors.Is(err, storages.ErrRowAlreadyExists) {
				if existsShortURLID, err := s.Storage.GetByOriginalURL(r.Context(), originalURL); err == nil {
					w.WriteHeader(http.StatusConflict)
					fmt.Fprintf(w, config.Get().BaseURL+"/"+existsShortURLID)
					return
				}
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, config.Get().BaseURL+"/"+shortURLID)
	}
}
