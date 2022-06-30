package handlers

import (
	"github.com/ansedo/url-shortener/internal/services/shortener"
	"net/http"
)

func PingDB(s *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.Storage.Ping(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
