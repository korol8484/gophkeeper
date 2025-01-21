package list

import (
	"context"
	"github.com/korol8484/gophkeeper/internal/server/domain"
	"github.com/korol8484/gophkeeper/pkg/model"
	"testing"
)

type repoList struct {
}

func (r *repoList) GetAllByUserID(ctx context.Context, userID domain.UserID) ([]*model.Secret, error) {
	return nil, nil
}
func TestSecretServiceList_GetAllByUserID(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repoList{}
			s := NewSecretServiceList(repo)
			s.GetAllByUserID(nil, 1)
		})
	}
}
