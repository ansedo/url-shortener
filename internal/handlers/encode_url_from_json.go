package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/services/shortener"
	"io"
	"net/http"
	"net/url"
)

type request struct {
	URL string `json:"url"`
}

type response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func EncodeURLFromJSON(s *shortener.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Add("X-Content-Type-Options", "nosniff")

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			writeBadRequestErrorJSON(w, err.Error())
			return
		}

		req := &request{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			writeBadRequestErrorJSON(w, err.Error())
			return
		}

		uri, err := url.ParseRequestURI(req.URL)
		if err != nil {
			writeBadRequestErrorJSON(w, config.New().RequestNotAllowedError)
			return
		}

		id, err := s.GenerateID()
		if err != nil {
			writeBadRequestErrorJSON(w, err.Error())
			return
		}

		err = s.Storage.Set(id, uri.String())
		if err != nil {
			writeBadRequestErrorJSON(w, err.Error())
			return
		}

		resp, err := json.Marshal(&response{Result: config.New().BaseURL + "/" + id})
		if err != nil {
			writeBadRequestErrorJSON(w, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err = fmt.Fprint(w, string(resp))
		if err != nil {
			writeBadRequestErrorJSON(w, err.Error())
			return
		}
	}
}

func writeBadRequestErrorJSON(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response{Error: err})
}
