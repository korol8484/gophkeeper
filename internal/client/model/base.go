package model

//type base struct {
//	id      uuid.UUID
//	version int
//	time    time.Time
//}

type BaseI interface {
	GetType() string
	View() string
}
