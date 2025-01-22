package model

import (
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"time"
)

type Secret struct {
	ID        uuid.UUID              `json:"id"`
	MetaData  map[string]interface{} `json:"meta_data"`
	Version   int                    `json:"version"`
	Content   []byte                 `json:"content"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type SecretCreateRequest struct {
	MetaData map[string]interface{} `json:"meta_data"`
	Content  []byte                 `json:"content"`
}
