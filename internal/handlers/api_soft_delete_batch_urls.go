package handlers

import (
	"encoding/json"
	"github.com/ansedo/url-shortener/internal/helpers"
	"github.com/ansedo/url-shortener/internal/models"
	"github.com/ansedo/url-shortener/internal/services/shortener"
	"net/http"
)

func APISoftDeleteBatchURLs(s *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Add("X-Content-Type-Options", "nosniff")

		var batch []string
		if err := json.NewDecoder(r.Body).Decode(&batch); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ShortenResponse{Error: err.Error()})
			return
		}

		uid := helpers.GetUIDFromCtx(r.Context())
		shortenList := make([]models.ShortenLink, len(batch))
		for _, url := range batch {
			shortenList = append(shortenList, models.ShortenLink{
				UID:        uid,
				ShortURLID: url,
			})
		}

		w.WriteHeader(http.StatusAccepted)
		s.Storage.AsyncSoftDeleteBatch(r.Context(), shortenList)
	}
}
