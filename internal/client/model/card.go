package model

import (
	"fmt"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/korol8484/gophkeeper/pkg/model"
	"strings"
)

type Card struct {
	id     uuid.UUID
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

func (p *Card) GetType() Type {
	return TypeCard
}

func (p *Card) View() string {
	return fmt.Sprintf("Card: %s\nDate: %s\\%s\nCvv:%s", p.number, p.year, p.month, p.cvv)
}

func (p *Card) GetId() uuid.UUID {
	return p.id
}

func (p *Card) GetTitle() string {
	return p.title
}

func (p *Card) ToModel() *model.SecretCreateRequest {
	return &model.SecretCreateRequest{
		MetaData: map[string]interface{}{
			"title": p.title,
			typeKey: p.GetType(),
		},
		Content: []byte(strings.Join([]string{p.number, p.year, p.month, p.cvv}, "||")),
	}
}

func (p *Card) load(data *model.Secret) BaseI {
	var title string
	if t, ok := data.MetaData["title"]; ok {
		title = t.(string)
	}

	var number, year, month, cvv string

	for i, v := range strings.Split(string(data.Content), "||") {
		switch i {
		case 0:
			number = v
		case 1:
			year = v
		case 2:
			month = v
		case 3:
			cvv = v
		}
	}

	return &Card{
		id:     data.ID,
		title:  title,
		number: number,
		year:   year,
		month:  month,
		cvv:    cvv,
	}
}
