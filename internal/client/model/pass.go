package model

import (
	"github.com/korol8484/gophkeeper/pkg/model"
)

type Password struct {
	title    string
	login    string
	password string
}

func NewPassword(title, login, password string) *Password {
	return &Password{
		title:    title,
		login:    login,
		password: password,
	}
}

func (p *Password) GetType() string {
	return "password"
}

func (p *Password) View() string {
	return ""
}

func (p *Password) ToModel() *model.SecretCreateRequest {
	return &model.SecretCreateRequest{
		MetaData: map[string]interface{}{
			"login": p.login,
			"title": p.title,
			"type":  p.GetType(),
		},
		Content: []byte(p.password),
	}
}
