package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ansedo/url-shortener/internal/models"
	"github.com/ansedo/url-shortener/internal/services/shortener"
	"net/http"
)

func APIGetURLsByUID(s *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Add("X-Content-Type-Options", "nosniff")

		entities, err := s.Storage.GetByUID(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ShortenResponse{Error: err.Error()})
			return
		}

		if len(entities) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		resp, err := json.Marshal(&entities)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ShortenResponse{Error: err.Error()})
		}

		_, err = fmt.Fprint(w, string(resp))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ShortenResponse{Error: err.Error()})
			return
		}
	}
}
