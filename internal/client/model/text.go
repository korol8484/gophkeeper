package model

import (
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/korol8484/gophkeeper/pkg/model"
)

type Text struct {
	id    uuid.UUID
	title string
	text  string
}

func NewText(title, text string) *Text {
	return &Text{
		title: title,
		text:  text,
	}
}

func (p *Text) GetType() Type {
	return TypeText
}

func (p *Text) View() string {
	return p.text
}

func (p *Text) GetId() uuid.UUID {
	return p.id
}

func (p *Text) GetTitle() string {
	return p.title
}

func (p *Text) ToModel() *model.SecretCreateRequest {
	return &model.SecretCreateRequest{
		MetaData: map[string]interface{}{
			"title": p.title,
			typeKey: p.GetType(),
		},
		Content: []byte(p.text),
	}
}

func (p *Text) load(data *model.Secret) BaseI {
	var title string
	if t, ok := data.MetaData["title"]; ok {
		title = t.(string)
	}

	return &Text{
		id:    data.ID,
		title: title,
		text:  string(data.Content),
	}
}
