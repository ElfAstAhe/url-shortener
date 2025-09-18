package dto

// UserShorten is user shorten url
type UserShorten struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
