package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	_srv "github.com/ElfAstAhe/url-shortener/internal/service"
	"github.com/golang-jwt/jwt/v4"
)

type AppClaims struct {
	jwt.RegisteredClaims
	Roles  []string `json:"roles,omitempty"`
	User   string   `json:"user,omitempty"`
	UserID string   `json:"user_id,omitempty"`
}

const secretKey = "secret-key"
const timeExpirationDuration = time.Minute * 30

const TestUser = "test_user"
const TestUserID = "test_user_id"

const Cookie = "Authorization"

var TestRoles _srv.Roles = []string{
	"test_role1",
	"test_role2",
}

func NewTokenString(userID string, user string, roles ...string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newAppClaims(userID, user, roles...))

	res, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return res, nil
}

func newAppClaims(userID string, user string, roles ...string) *AppClaims {
	return &AppClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(timeExpirationDuration)),
		},
		UserID: userID,
		User:   user,
		Roles:  roles,
	}
}

func GetUserID(r *http.Request) (string, error) {
	if r == nil {
		return "", nil
	}

	token, err := getAuthToken(r)
	if err != nil {
		return "", err
	}
	if token == nil {
		return "", nil
	}

	if !token.Valid {
		return "", errors.New("invalid jwt token")
	}

	return token.Claims.(*AppClaims).UserID, nil
}

func getAuthToken(r *http.Request) (*jwt.Token, error) {
	cookie, _ := r.Cookie(Cookie)
	if cookie == nil {
		return nil, nil
	}
	if err := cookie.Valid(); err != nil {
		return nil, err
	}

	claims := &AppClaims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}
