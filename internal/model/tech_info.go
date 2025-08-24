package model

type TechData struct {
	CreateUser string `db:"create_user"`
	Created    string `db:"created"`
	UpdateUser string `db:"update_user"`
	Updated    string `db:"updated"`
}
