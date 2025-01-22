package secret

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/korol8484/gophkeeper/internal/server/api/util"
	"github.com/korol8484/gophkeeper/internal/server/domain"
	"github.com/korol8484/gophkeeper/pkg"
	"github.com/korol8484/gophkeeper/pkg/model"
	"io"
	"net/http"
)

type AddHandler struct {
	service secretServiceI
}

type secretServiceI interface {
	Add(ctx context.Context, metaData map[string]interface{}, context []byte, userId domain.UserID) (*uuid.UUID, error)
}

func NewSecretAddHandler(service secretServiceI) *AddHandler {
	return &AddHandler{service: service}
}

func (h *AddHandler) secretAddHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := util.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "внутренняя ошибка сервера. Не опознан пользователь", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	dto := &model.SecretCreateRequest{}

	if err = json.Unmarshal(body, dto); err != nil {
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	_, err = h.service.Add(ctx, dto.MetaData, dto.Content, userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// RegisterRoutes - Register user api routes
func (h *AddHandler) RegisterRoutes(loader util.AuthSession) func(mux *chi.Mux) {
	return func(mux *chi.Mux) {
		mux.With(util.CheckAuth(loader)).Post(pkg.SecretAPIAdd, h.secretAddHandler)
	}
}
