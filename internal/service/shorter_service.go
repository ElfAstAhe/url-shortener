package service

// ShorterService app service
type ShorterService interface {
	// GetUrl return full URL
	GetUrl(key string) (string, error)

	// Store URL and return short key
	Store(url string) (string, error)
}
