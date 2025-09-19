package auth

import (
	"context"
	"errors"
)

type Roles []string

type UserInfo struct {
	Admin  bool   `json:"admin"`
	UserID string `json:"user_id,omitempty"`
	User   string `json:"user,omitempty"`
	Roles  Roles  `json:"roles,omitempty"`
}

func NewUserInfo(admin bool, userID string, user string, roles Roles) *UserInfo {
	return &UserInfo{
		Admin:  admin,
		UserID: userID,
		User:   user,
		Roles:  roles,
	}
}

func UserInfoFromContext(ctx context.Context) (*UserInfo, error) {
	res, ok := ctx.Value("UserInfo").(*UserInfo)
	if !ok {
		return nil, errors.New("user info not found")
	}

	return res, nil
}
