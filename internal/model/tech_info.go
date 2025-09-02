package model

type TechData struct {
	CreateUser string `db:"create_user" json:"create_user,omitempty"`
	Created    string `db:"created" json:"created,omitempty"`
	UpdateUser string `db:"update_user" json:"update_user,omitempty"`
	Updated    string `db:"updated" json:"updated,omitempty"`
}
