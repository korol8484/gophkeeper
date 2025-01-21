package secret

import (
	"context"
	"fmt"
	"github.com/korol8484/gophkeeper/internal/server/domain"
	"github.com/korol8484/gophkeeper/pkg/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockSecretServiceList struct {
	flag string
}

func (s *mockSecretServiceList) setFlag(flag string) {
	s.flag = flag
}
func (s *mockSecretServiceList) GetAllByUserID(ctx context.Context, userID domain.UserID) ([]*model.Secret, error) {
	if s.flag == "success get" {
		list := []*model.Secret{}
		return list, nil
	}
	return nil, fmt.Errorf("some error")
}

func TestListHandler_secretGetAllHandler(t *testing.T) {
	tests := []struct {
		name   string
		userId domain.UserID
		status int
		flag   string
	}{
		{
			name:   "user not exists",
			userId: -1,
			flag:   "",
			status: http.StatusUnauthorized,
		},
		{
			name:   "service error",
			userId: 1,
			flag:   "",
			status: http.StatusInternalServerError,
		},
		{
			name:   "success",
			userId: 1,
			flag:   "success get",
			status: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/user/secret/", nil)
			response := httptest.NewRecorder()
			ctx := request.Context()
			if tt.userId != -1 {
				ctx = context.WithValue(ctx, "ctx_user_id", tt.userId)
				request = request.WithContext(ctx)
			}
			service := &mockSecretServiceList{}
			service.setFlag(tt.flag)
			h := NewListHandler(service)
			h.secretGetAllHandler(response, request)
			assert.Equal(t, tt.status, response.Code)
		})
	}
}
