package add

import (
	"context"
	"github.com/google/uuid"
	"github.com/korol8484/gophkeeper/internal/server/domain"
	"time"
)

type SecretService struct {
	repository repositoryI
}

type repositoryI interface {
	Add(ctx context.Context, metaData map[string]interface{}, context []byte, userId domain.UserID, version int, added time.Time) (*uuid.UUID, error)
}

func NewSecretService(repository repositoryI) *SecretService {
	return &SecretService{repository: repository}
}

func (s *SecretService) Add(ctx context.Context, metaData map[string]interface{}, context []byte, userId domain.UserID) (*uuid.UUID, error) {
	return s.repository.Add(ctx, metaData, context, userId, 1, time.Now())
}
