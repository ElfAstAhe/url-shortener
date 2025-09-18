package service

type Roles []string

type AuthService interface {
	Authenticate(user, password string) (bool, error)
	Authorize(user string) (Roles, error)
}
