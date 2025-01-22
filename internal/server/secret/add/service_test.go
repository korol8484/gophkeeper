package add

import (
	"context"
	"github.com/google/uuid"
	"github.com/korol8484/gophkeeper/internal/server/domain"
	"testing"
	"time"
)

type repoAdd struct {
}

func (r *repoAdd) Add(ctx context.Context, metaData map[string]interface{}, context []byte, userId domain.UserID, version int, added time.Time) (*uuid.UUID, error) {
	return nil, nil
}
func TestSecretService_Add(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repoAdd{}
			s := NewSecretService(repo)
			s.Add(nil, nil, nil, 1)
		})
	}
}
