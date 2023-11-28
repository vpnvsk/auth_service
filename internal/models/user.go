package models

type User struct {
	Uuid         string `db:"id"`
	Login        string `db:"login"`
	PasswordHash string `db:"passwordHash"`
}
