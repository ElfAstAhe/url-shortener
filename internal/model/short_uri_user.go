package model

type ShortURIUser struct {
	ID         string `db:"id" json:"id"`
	ShortURIID string `db:"short_uri_id" json:"short_uri_id"`
	UserID     string `db:"user_id" json:"user_id"`
}

func NewShortURIUser(shortURIID string, userID string) (*ShortURIUser, error) {
	// ToDo: validation
	// ..

	return &ShortURIUser{
		ShortURIID: shortURIID,
		UserID:     userID,
	}, nil
}
