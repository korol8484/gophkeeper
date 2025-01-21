package model

type Type string

const (
	TypePassword Type = "password"
	TypeText     Type = "text"
	TypeCard     Type = "card"
	TypeBinary   Type = "binary"
)

const typeKey = "type"

type BaseI interface {
	GetType() Type
	View() string
}
