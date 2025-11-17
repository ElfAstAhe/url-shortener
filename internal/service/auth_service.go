package service

import "github.com/ElfAstAhe/url-shortener/internal/service/auth"

type AuthService interface {
	Authenticate(user, password string) (bool, error)
	Authorize(user string) (auth.Roles, error)
}
