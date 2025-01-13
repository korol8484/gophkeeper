package model

import "time"

type Secret struct {
	ID        int64                  `json:"id"`
	MetaData  map[string]interface{} `json:"meta_data"`
	Version   int                    `json:"version"`
	Content   []byte                 `json:"content"`
	UpdatedAt time.Time              `json:"updated_at"`
}
