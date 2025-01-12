package domain

import "time"

type Secret struct {
	ID        int64                  `json:"id"`
	UserID    UserID                 `json:"user_id"`
	MetaData  map[string]interface{} `json:"meta_data"`
	Version   int                    `json:"version"`
	Content   []byte                 `json:"content"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}
