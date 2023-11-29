package models

import "github.com/google/uuid"

type User struct {
	Uuid         uuid.UUID `db:"id"`
	Login        string    `db:"login"`
	PasswordHash []byte    `db:"password_hash"`
}
