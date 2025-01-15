package crypt

import (
	"github.com/liangrog/vault"
)

type Service struct{}

// NewService - factory
func NewService() *Service {
	return &Service{}
}

// Encrypt - crypt data by key
func (s *Service) Encrypt(data []byte, key string) ([]byte, error) {
	secret, err := vault.Encrypt(data, key)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

// Decrypt - encrypt data by key
func (s *Service) Decrypt(data []byte, key string) ([]byte, error) {
	secret, err := vault.Decrypt(key, data)
	if err != nil {
		return nil, err
	}

	return secret, nil
}
