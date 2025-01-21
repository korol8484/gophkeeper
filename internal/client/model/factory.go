package model

import "github.com/korol8484/gophkeeper/pkg/model"

func LoadModels(data []model.Secret) []BaseI {
	res := make([]BaseI, 0, len(data))
	for _, v := range data {
		mt := LoadModel(&v)
		if mt == nil {
			continue
		}

		res = append(res, mt)
	}

	return res
}

func LoadModel(data *model.Secret) BaseI {
	mt, ok := data.MetaData[typeKey]
	if !ok {
		return nil
	}

	switch Type(mt.(string)) {
	case TypeText:
		return new(Text).load(data)
	case TypePassword:
		return new(Password).load(data)
	case TypeCard:
		return new(Card).load(data)
	case TypeBinary:
		return new(Binary).load(data)
	}

	return nil
}
