package service

// ShorterService app service
type ShorterService interface {
	// GetURL return full URL
	GetURL(key string) (string, error)

	// Store URL and return short key
	Store(url string) (string, error)
}
