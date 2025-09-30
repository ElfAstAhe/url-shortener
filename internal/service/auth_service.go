package service

import _auth "github.com/ElfAstAhe/url-shortener/internal/service/auth"

type AuthService interface {
	Authenticate(user, password string) (bool, error)
	Authorize(user string) (_auth.Roles, error)
}
