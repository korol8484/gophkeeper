package list

import (
	"context"
	"github.com/korol8484/gophkeeper/internal/server/domain"
	"github.com/korol8484/gophkeeper/pkg/model"
)

type SecretServiceList struct {
	repository repositoryI
}

type repositoryI interface {
	GetAllByUserID(ctx context.Context, userID domain.UserID) ([]*model.Secret, error)
}

func NewSecretServiceList(repository repositoryI) *SecretServiceList {
	return &SecretServiceList{repository: repository}
}

func (s *SecretServiceList) GetAllByUserID(ctx context.Context, userID domain.UserID) ([]*model.Secret, error) {
	return s.repository.GetAllByUserID(ctx, userID)
}
