package model

import (
	"github.com/korol8484/gophkeeper/pkg/model"
	"time"
)

type Password struct {
	base
	login    string
	password string
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

func (p *Password) ToModel() *model.Secret {
	return &model.Secret{
		MetaData: map[string]interface{}{
			"login":    p.login,
			"password": p.password,
		},
		Version:   0,
		Content:   nil,
		UpdatedAt: time.Time{},
	}
}

//func (p *Password) Title() string {
//	return p.login
//}
