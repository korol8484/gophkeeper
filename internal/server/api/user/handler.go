package user

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/korol8484/gophkeeper/internal/server/domain"
	"net/http"
)

// AuthUser - user api interface
type AuthUser interface {
	CreateUser(ctx context.Context, user *domain.User, password string) (*domain.User, error)
	Auth(ctx context.Context, login, password string) (*domain.User, error)
}

// AuthSession - session interface
type AuthSession interface {
	CreateSession(w http.ResponseWriter, r *http.Request, id domain.UserID) error
}

// Handler -
// В задаче не учитывается, что делать, если сессия уже существует (аутентифицированый пользователь)
// поэтому такие нюансы тут не учитываем и считаем, что обращения на данные методы всегда выполняют поставленные задачи
type Handler struct {
	auth    AuthUser
	session AuthSession
}

// NewAuthHandler - factory
func NewAuthHandler(auth AuthUser, session AuthSession) *Handler {
	return &Handler{
		auth:    auth,
		session: session,
	}
}

type request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// RegisterRoutes - Register user api routes
func (h *Handler) RegisterRoutes() func(mux *chi.Mux) {
	return func(mux *chi.Mux) {
		mux.Post("/api/user/register", h.registerHandler)
		mux.Post("/api/user/login", h.authHandler)
	}
}
