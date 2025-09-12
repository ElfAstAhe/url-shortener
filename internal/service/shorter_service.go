package service

type CorrelationUrls map[string]string

type CorrelationShorts map[string]string

// ShorterService app service
type ShorterService interface {
	// GetURL return full URL
	GetURL(key string) (string, error)

	// Store URL and return short key
	Store(url string) (string, error)

	// BatchStore URLs and return correlation shorts
	BatchStore(source CorrelationUrls) (CorrelationShorts, error)
}
