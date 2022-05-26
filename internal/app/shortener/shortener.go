package shortener

import (
	"errors"
	"github.com/ansedo/url-shortener/internal/app/config"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	MethodNotAllowedError  = "this http method is not allowed"
	RequestNotAllowedError = "this request is not allowed"
)

type Shortener struct {
	db map[string]string
}

func NewShortener() *Shortener {
	return &Shortener{
		db: make(map[string]string),
	}
}

func (s *Shortener) GetURLbyID(id string) (string, error) {
	URL, ok := s.db[id]
	if !ok {
		return "", errors.New(RequestNotAllowedError)
	}
	return URL, nil
}

func (s *Shortener) GetIDbyURL(URL string) (string, error) {
	validURL, err := url.ParseRequestURI(URL)
	if err != nil {
		return "", err
	}
	for id, dbURL := range s.db {
		if dbURL == validURL.String() {
			return id, nil
		}
	}
	newID := strconv.Itoa(len(s.db))
	s.db[newID] = validURL.String()
	return newID, nil
}

func (s *Shortener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		redirectURL, err := s.GetURLbyID(r.URL.Path[1:])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)

	case http.MethodPost:
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := s.GetIDbyURL(string(bytes))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(config.SiteAddress + "/" + id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	default:
		http.Error(w, MethodNotAllowedError, http.StatusMethodNotAllowed)
	}
}
