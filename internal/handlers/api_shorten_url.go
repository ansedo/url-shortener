package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/models"
	"github.com/ansedo/url-shortener/internal/services/shortener"
	"github.com/ansedo/url-shortener/internal/storages"
	"io"
	"net/http"
	"net/url"
)

func APIShortenURL(s *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Add("X-Content-Type-Options", "nosniff")

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ShortenResponse{Error: err.Error()})
			return
		}

		req := &models.ShortenRequest{}
		if err = json.Unmarshal(body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ShortenResponse{Error: err.Error()})
			return
		}

		if _, err = url.ParseRequestURI(req.URL); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ShortenResponse{Error: ErrRequestNotAllowed.Error()})
			return
		}

		shortURLID, err := s.GenerateID(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ShortenResponse{Error: err.Error()})
			return
		}

		if err = s.Storage.Add(r.Context(), shortURLID, req.URL); err != nil {
			if errors.Is(err, storages.ErrRowAlreadyExists) {
				if existsShortURLID, err := s.Storage.GetByOriginalURL(r.Context(), req.URL); err == nil {
					w.WriteHeader(http.StatusConflict)
					json.NewEncoder(w).Encode(models.ShortenResponse{Result: config.Get().BaseURL + "/" + existsShortURLID})
					return
				}
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ShortenResponse{Error: err.Error()})
			return
		}

		resp, err := json.Marshal(&models.ShortenResponse{Result: config.Get().BaseURL + "/" + shortURLID})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ShortenResponse{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusCreated)
		if _, err = fmt.Fprint(w, string(resp)); err != nil {
			json.NewEncoder(w).Encode(models.ShortenResponse{Error: err.Error()})
			return
		}
	}
}
