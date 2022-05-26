package shortener

import (
	"github.com/ansedo/url-shortener/internal/app/config"
	"github.com/ansedo/url-shortener/internal/app/storage"
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
	storage storage.Storage
}

func NewShortener(storage storage.Storage) *Shortener {
	return &Shortener{
		storage: storage,
	}
}

func (s *Shortener) GenerateId() (string, error) {
	// Strong and tiny approach with extra short ids (at least for first 10 values) and no collisions!
	return strconv.Itoa(s.storage.Count()), nil
}

func (s *Shortener) EncodeUrl(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uri, err := url.ParseRequestURI(string(body))
	if err != nil {
		http.Error(w, RequestNotAllowedError, http.StatusBadRequest)
		return
	}

	id, err := s.GenerateId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.storage.Set(id, uri.String())
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
}

func (s *Shortener) DecodeUrl(w http.ResponseWriter, r *http.Request) {
	uri, err := s.storage.Get(r.URL.Path[1:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, uri, http.StatusTemporaryRedirect)
}

func (s *Shortener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.DecodeUrl(w, r)
	case http.MethodPost:
		s.EncodeUrl(w, r)
	default:
		http.Error(w, MethodNotAllowedError, http.StatusMethodNotAllowed)
	}
}
