package util

type PrivateKey struct {
	ID        string `json:"id"`
	Algorithm string `json:"algorithm"`
	Key       string `json:"key"`
}

func ToPtr[T any](val T) *T {
	return &val
}
