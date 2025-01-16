package model

import (
	"github.com/korol8484/gophkeeper/pkg/model"
)

type Binary struct {
	title string
}

func NewBinary(title string) *Binary {
	return &Binary{
		title: title,
	}
}

func (p *Binary) GetType() string {
	return "binary"
}

func (p *Binary) View() string {
	return ""
}

func (p *Binary) ToModel() *model.SecretCreateRequest {
	return &model.SecretCreateRequest{
		MetaData: map[string]interface{}{
			"title": p.title,
			"type":  p.GetType(),
		},
		Content: []byte(""),
	}
}
