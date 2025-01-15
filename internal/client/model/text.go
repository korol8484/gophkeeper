package model

import (
	"github.com/korol8484/gophkeeper/pkg/model"
	"time"
)

type Text struct {
	base
	title string
	text  string
}

func NewText(title, text string) *Text {
	return &Text{
		base: base{
			version: 1,
			time:    time.Now(),
		},
		title: title,
		text:  text,
	}
}

func (p *Text) ToModel() *model.Secret {
	return &model.Secret{
		MetaData: map[string]interface{}{
			"title": p.title,
			"type":  "text",
		},
		Version:   p.version,
		Content:   []byte(p.text),
		UpdatedAt: p.time,
	}
}
