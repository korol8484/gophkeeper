package get

import (
	"context"
	"github.com/google/uuid"
	"github.com/korol8484/gophkeeper/internal/server/domain"
	"github.com/korol8484/gophkeeper/pkg/model"
)

type SecretServiceGet struct {
	repository repositoryI
}

type repositoryI interface {
	Get(ctx context.Context, userID domain.UserID, ID uuid.UUID) (*model.Secret, error)
}

func NewSecretServiceList(repository repositoryI) *SecretServiceGet {
	return &SecretServiceGet{repository: repository}
}

func (s *SecretServiceGet) Get(ctx context.Context, userID domain.UserID, ID uuid.UUID) (*model.Secret, error) {
	return s.repository.Get(ctx, userID, ID)
}
