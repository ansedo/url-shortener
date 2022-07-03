package models

type (
	ShortenRequest struct {
		URL string `json:"url"`
	}
	ShortenResponse struct {
		Result string `json:"result,omitempty"`
		Error  string `json:"error,omitempty"`
	}
	ShortenList struct {
		ShortURLID    string `json:"short_url_id,omitempty"`
		ShortURL      string `json:"short_url,omitempty"`
		OriginalURL   string `json:"original_url,omitempty"`
		CorrelationID string `json:"correlation_id,omitempty"`
	}
)
