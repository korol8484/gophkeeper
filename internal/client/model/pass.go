package model

import (
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/korol8484/gophkeeper/pkg/model"
)

type Password struct {
	id       uuid.UUID
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

func (p *Password) GetType() Type {
	return TypePassword
}

func (p *Password) View() string {
	return p.password
}

func (p *Password) GetId() uuid.UUID {
	return p.id
}

func (p *Password) GetTitle() string {
	return p.title
}

func (p *Password) ToModel() *model.SecretCreateRequest {
	return &model.SecretCreateRequest{
		MetaData: map[string]interface{}{
			"login": p.login,
			"title": p.title,
			typeKey: p.GetType(),
		},
		Content: []byte(p.password),
	}
}

func (p *Password) load(data *model.Secret) BaseI {
	var title string
	if t, ok := data.MetaData["title"]; ok {
		title = t.(string)
	}

	return &Password{
		id:       data.ID,
		title:    title,
		password: string(data.Content),
	}
}
