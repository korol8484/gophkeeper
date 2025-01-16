package model

import (
	"github.com/korol8484/gophkeeper/pkg/model"
)

type Text struct {
	title string
	text  string
}

func NewText(title, text string) *Text {
	return &Text{
		title: title,
		text:  text,
	}
}

func (p *Text) GetType() string {
	return "text"
}

func (p *Text) View() string {
	return ""
}

func (p *Text) ToModel() *model.SecretCreateRequest {
	return &model.SecretCreateRequest{
		MetaData: map[string]interface{}{
			"title": p.title,
			"type":  p.GetType(),
		},
		Content: []byte(p.text),
	}
}
