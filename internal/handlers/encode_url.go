package handlers

import (
	"github.com/ansedo/url-shortener/internal/services/shortener"
	"io"
	"net/http"
	"net/url"
)

func EncodeURL(s *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		uri, err := url.ParseRequestURI(string(body))
		if err != nil {
			http.Error(w, s.Config.RequestNotAllowedError, http.StatusBadRequest)
			return
		}

		id, err := s.GenerateID()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.Storage.Set(id, uri.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(s.Config.SiteAddress + "/" + id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
