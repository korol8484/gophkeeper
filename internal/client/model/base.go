package model

import "time"

type ModelType string

type base struct {
	version int
	time    time.Time
}
