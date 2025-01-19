package secret

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/korol8484/gophkeeper/internal/server/domain"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockSecretService struct {
	flag string
}

func (s *mockSecretService) setFlag(flag string) {
	s.flag = flag
}

func (s *mockSecretService) Add(ctx context.Context, metaData map[string]interface{}, context []byte, userId domain.UserID) (*uuid.UUID, error) {
	if s.flag == "success save" {
		myuuid, _ := uuid.NewUUID()
		return &myuuid, nil
	}
	return nil, fmt.Errorf("some error")
}

func TestAddHandler_secretAddHandler(t *testing.T) {
	tests := []struct {
		name   string
		userId domain.UserID
		flag   string
		status int
		body   string
	}{
		{
			name:   "user not exists",
			userId: -1,
			flag:   "",
			status: http.StatusUnauthorized,
			body:   "",
		},
		{
			name:   "bad body",
			userId: 1,
			flag:   "",
			status: http.StatusBadRequest,
			body:   "",
		},
		{
			name:   "fail save",
			userId: 1,
			flag:   "fail save",
			status: http.StatusInternalServerError,
			body:   `{"meta_data":{"foo":"bar"},"content":[49,48,58,50,52,58,50,54]}`,
		},
		{
			name:   "success save",
			userId: 1,
			flag:   "success save",
			status: http.StatusCreated,
			body:   `{"meta_data":{"foo":"bar"},"content":[49,48,58,50,52,58,50,54]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/articles/1", strings.NewReader(tt.body))
			response := httptest.NewRecorder()
			ctx := request.Context()
			if tt.userId != -1 {
				ctx = context.WithValue(ctx, "ctx_user_id", tt.userId)
				request = request.WithContext(ctx)
			}
			service := &mockSecretService{}
			service.setFlag(tt.flag)
			h := NewSecretAddHandler(service)
			h.secretAddHandler(response, request)
			assert.Equal(t, tt.status, response.Code)
		})
	}
}
