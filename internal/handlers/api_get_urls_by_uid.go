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

		shortenList, err := s.Storage.GetByUID(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ShortenResponse{Error: err.Error()})
			return
		}

		if len(shortenList) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		for i := range shortenList {
			shortenList[i].ShortURL = s.BaseURL + "/" + shortenList[i].ShortURLID
			shortenList[i].ShortURLID = ""
		}

		resp, err := json.Marshal(&shortenList)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ShortenResponse{Error: err.Error()})
			return
		}

		if _, err = fmt.Fprint(w, string(resp)); err != nil {
			json.NewEncoder(w).Encode(models.ShortenResponse{Error: err.Error()})
			return
		}
	}
}
