package secret

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/korol8484/gophkeeper/internal/server/api/util"
	"github.com/korol8484/gophkeeper/internal/server/domain"
	"github.com/korol8484/gophkeeper/pkg/model"
	"net/http"
)

const SecretAPIGet = "/user/secret/{id}"

type GetHandler struct {
	service secretServiceGetI
}

type secretServiceGetI interface {
	Get(ctx context.Context, userID domain.UserID, ID uuid.UUID) (*model.Secret, error)
}

func NewGetHandler(service secretServiceGetI) *GetHandler {
	return &GetHandler{service: service}
}

func (h *GetHandler) secretGetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := util.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "внутренняя ошибка сервера. Не опознан пользователь", http.StatusUnauthorized)
		return
	}

	secretIDStr := chi.URLParam(r, "id")
	secretID, err := uuid.Parse(secretIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	list, err := h.service.Get(ctx, userId, secretID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, _ := json.Marshal(list)

	w.Write(bytes)
}

// RegisterRoutes - Register user api routes
func (h *GetHandler) RegisterRoutes(loader util.AuthSession) func(mux *chi.Mux) {
	return func(mux *chi.Mux) {
		mux.With(util.CheckAuth(loader)).Get(SecretAPIGet, h.secretGetHandler)
	}
}
