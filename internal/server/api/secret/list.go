package secret

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/korol8484/gophkeeper/internal/server/api/util"
	"github.com/korol8484/gophkeeper/internal/server/domain"
	"github.com/korol8484/gophkeeper/pkg/model"
	"net/http"
)

const SecretAPIList = "/user/secret"

type ListHandler struct {
	service secretServiceListI
}

type secretServiceListI interface {
	GetAllByUserID(ctx context.Context, userID domain.UserID) ([]*model.Secret, error)
}

func NewListHandler(service secretServiceListI) *ListHandler {
	return &ListHandler{service: service}
}

func (h *ListHandler) secretGetAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := util.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "внутренняя ошибка сервера. Не опознан пользователь", http.StatusUnauthorized)
		return
	}

	list, err := h.service.GetAllByUserID(ctx, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, _ := json.Marshal(list)

	w.Write(bytes)
}

// RegisterRoutes - Register user api routes
func (h *ListHandler) RegisterRoutes() func(mux *chi.Mux) {
	return func(mux *chi.Mux) {
		mux.Get(SecretAPIList, h.secretGetAllHandler)
	}
}
