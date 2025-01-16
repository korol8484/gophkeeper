package model

import (
	"github.com/korol8484/gophkeeper/pkg/model"
	"strings"
)

type Card struct {
	title  string
	number string
	year   string
	month  string
	cvv    string
}

func NewCard(title, number, year, month, cvv string) *Card {
	return &Card{
		title:  title,
		number: number,
		year:   year,
		month:  month,
		cvv:    cvv,
	}
}

func (p *Card) GetType() string {
	return "card"
}

func (p *Card) View() string {
	return ""
}

func (p *Card) ToModel() *model.SecretCreateRequest {
	return &model.SecretCreateRequest{
		MetaData: map[string]interface{}{
			"title": p.title,
			"type":  p.GetType(),
		},
		Content: []byte(strings.Join([]string{p.number, p.year, p.month, p.cvv}, "||")),
	}
}
