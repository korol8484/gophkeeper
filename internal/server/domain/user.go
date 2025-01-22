package domain

import "errors"

var (
	ErrIssetUser = errors.New("логин уже занят")
)

type UserID int64

type User struct {
	ID           UserID `json:"id"`
	Login        string `json:"login"`
	PasswordHash string `json:"-"`
}
