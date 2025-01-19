package get

import (
	"context"
	"github.com/google/uuid"
	"github.com/korol8484/gophkeeper/internal/server/domain"
	"github.com/korol8484/gophkeeper/pkg/model"
	"testing"
)

type repoGet struct {
}

func (r *repoGet) Get(ctx context.Context, userID domain.UserID, ID uuid.UUID) (*model.Secret, error) {
	return nil, nil
}

func TestSecretServiceGet_Get(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repoGet{}
			s := NewSecretServiceList(repo)
			s.Get(nil, 1, uuid.New())
		})
	}
}
