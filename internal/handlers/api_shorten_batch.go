package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ansedo/url-shortener/internal/models"
	"github.com/ansedo/url-shortener/internal/services/shortener"
	"io"
	"net/http"
	"net/url"
)

func APIShortenBatch(s *shortener.Shortener) http.HandlerFunc {
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

		var shortenList []models.ShortenLink
		if err = json.Unmarshal(body, &shortenList); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ShortenResponse{Error: err.Error()})
			return
		}

		for i := range shortenList {
			if _, err = url.ParseRequestURI(shortenList[i].OriginalURL); err != nil {
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

			shortenList[i].ShortURLID = shortURLID
			shortenList[i].ShortURL = s.BaseURL + "/" + shortenList[i].ShortURLID
		}

		if err = s.Storage.AddBatch(r.Context(), shortenList); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ShortenResponse{Error: err.Error()})
			return
		}

		for i := range shortenList {
			shortenList[i].ShortURLID = ""
			shortenList[i].OriginalURL = ""
		}

		resp, err := json.Marshal(&shortenList)
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
