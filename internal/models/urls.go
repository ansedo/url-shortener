package models

type (
	ShortenRequest struct {
		URL string `json:"url"`
	}
	ShortenResponse struct {
		Result string `json:"result,omitempty"`
		Error  string `json:"error,omitempty"`
	}
	ShortenListResponse struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}
)
