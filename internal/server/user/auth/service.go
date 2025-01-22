package auth

import (
	"context"
	"github.com/korol8484/gophkeeper/internal/server/domain"
	"github.com/zitadel/passwap"
	"github.com/zitadel/passwap/argon2"
)

type userStore interface {
	AddUser(ctx context.Context, user *domain.User) (*domain.User, error)
	FindByLogin(ctx context.Context, login string) (*domain.User, error)
}

// Service - user auth service
type Service struct {
	store   userStore
	hashSvc *passwap.Swapper
}

// NewService - factory
func NewService(uStore userStore) *Service {
	return &Service{
		store: uStore,
		// hashSvc - можно добавить интерфейс, принимать как зависимость, но в данном случае лишнее
		hashSvc: passwap.NewSwapper(
			argon2.NewArgon2id(argon2.RecommendedIDParams),
			argon2.Verifier,
		),
	}
}

// CreateUser - add new user
func (s *Service) CreateUser(ctx context.Context, user *domain.User, password string) (*domain.User, error) {
	encPas, err := s.hashSvc.Hash(password)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = encPas
	if user, err = s.store.AddUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Auth - authenticate user by login\password
func (s *Service) Auth(ctx context.Context, login, password string) (*domain.User, error) {
	user, err := s.store.FindByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	_, err = s.hashSvc.Verify(user.PasswordHash, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
