package service

type ShorterService interface {
	GetUrl(key string) (string, error)
	Store(url string) (string, error)
}
