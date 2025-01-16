package model

import (
	"github.com/korol8484/gophkeeper/pkg/model"
	"time"
)

type Password struct {
	base
	title    string
	login    string
	password string
}

func NewPassword(title, login, password string) *Password {
	return &Password{
		base: base{
			version: 1,
			time:    time.Now(),
		},
		title:    title,
		login:    login,
		password: password,
	}
}

//func (p *Password) Format(r service.Render) error {
//	return r.Render(p.)
//
//
//	_, err := w.Write([]byte(fmt.Sprintf("login: %s\n", p.login)))
//	if err != nil {
//		return err
//	}
//
//	_, err = w.Write([]byte(fmt.Sprintf("password: %s\n", p.password)))
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func (p *Password) ToModel() *model.SecretCreateRequest {
	return &model.SecretCreateRequest{
		MetaData: map[string]interface{}{
			"login": p.login,
			"title": p.title,
			"type":  "password",
		},
		Content: []byte(p.password),
	}
}

//func (p *Password) Title() string {
//	return p.login
//}
