package secret

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/korol8484/gophkeeper/internal/server/domain"
	"github.com/korol8484/gophkeeper/pkg/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockSecretServiceGet struct {
	flag string
}

func (s *mockSecretServiceGet) setFlag(flag string) {
	s.flag = flag
}
func (s *mockSecretServiceGet) Get(ctx context.Context, userID domain.UserID, ID uuid.UUID) (*model.Secret, error) {
	if s.flag == "success get" {
		list := model.Secret{}
		return &list, nil
	}
	return nil, fmt.Errorf("some error")
}

func TestGetHandler_secretGetHandler(t *testing.T) {
	tests := []struct {
		name   string
		userId domain.UserID
		id     string
		status int
		flag   string
	}{
		{
			name:   "user not exists",
			id:     "",
			userId: -1,
			flag:   "",
			status: http.StatusUnauthorized,
		},
		{
			name:   "not uuid",
			id:     "12",
			userId: 1,
			flag:   "",
			status: http.StatusInternalServerError,
		},
		{
			name:   "sql error",
			id:     uuid.New().String(),
			userId: 1,
			flag:   "",
			status: http.StatusInternalServerError,
		},
		{
			name:   "success",
			id:     uuid.New().String(),
			userId: 1,
			flag:   "success get",
			status: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/user/secret/"+tt.id, nil)
			response := httptest.NewRecorder()
			chiCtx := chi.NewRouteContext()

			ctx := request.Context()
			if tt.userId != -1 {
				ctx = context.WithValue(ctx, "ctx_user_id", tt.userId)
				request = request.WithContext(ctx)
			}
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx))
			chiCtx.URLParams.Add("id", tt.id)

			service := &mockSecretServiceGet{}
			service.setFlag(tt.flag)
			h := NewGetHandler(service)
			h.secretGetHandler(response, request)
			assert.Equal(t, tt.status, response.Code)
		})
	}
}
