package model

import (
	"github.com/google/uuid"
	"time"
)

type base struct {
	id      uuid.UUID
	version int
	time    time.Time
}
