package model

import (
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/korol8484/gophkeeper/pkg/model"
)

type Binary struct {
	id    uuid.UUID
	title string
}

func NewBinary(title string) *Binary {
	return &Binary{
		title: title,
	}
}

func (p *Binary) GetType() Type {
	return TypeBinary
}

func (p *Binary) View() string {
	return ""
}

func (p *Binary) GetId() uuid.UUID {
	return p.id
}

func (p *Binary) GetTitle() string {
	return p.title
}

func (p *Binary) ToModel() *model.SecretCreateRequest {
	return &model.SecretCreateRequest{
		MetaData: map[string]interface{}{
			"title": p.title,
			typeKey: p.GetType(),
		},
		Content: []byte(""),
	}
}

func (p *Binary) load(data *model.Secret) BaseI {
	var title string
	if t, ok := data.MetaData["title"]; ok {
		title = t.(string)
	}

	return &Binary{
		id:    data.ID,
		title: title,
	}
}
