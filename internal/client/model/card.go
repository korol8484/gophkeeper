package model

import (
	"github.com/korol8484/gophkeeper/pkg/model"
	"strings"
	"time"
)

type Card struct {
	base
	title  string
	number string
	year   string
	month  string
	cvv    string
}

func NewCard(title, number, year, month, cvv string) *Card {
	return &Card{
		base: base{
			version: 1,
			time:    time.Now(),
		},
		title:  title,
		number: number,
		year:   year,
		month:  month,
		cvv:    cvv,
	}
}

func (p *Card) ToModel() *model.Secret {
	return &model.Secret{
		MetaData: map[string]interface{}{
			"title": p.title,
			"type":  "card",
		},
		Version:   p.version,
		Content:   []byte(strings.Join([]string{p.number, p.year, p.month, p.cvv}, "||")),
		UpdatedAt: p.time,
	}
}
