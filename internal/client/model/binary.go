package model

import (
	"github.com/korol8484/gophkeeper/pkg/model"
	"time"
)

type Binary struct {
	base
	title string
}

func NewBinary(title string) *Binary {
	return &Binary{
		base: base{
			version: 1,
			time:    time.Now(),
		},
		title: title,
	}
}

func (p *Binary) ToModel() *model.SecretCreateRequest {
	return &model.SecretCreateRequest{
		MetaData: map[string]interface{}{
			"title": p.title,
			"type":  "binary",
		},
		Content: []byte(""),
	}
}
