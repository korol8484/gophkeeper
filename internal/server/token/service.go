package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/korol8484/gophkeeper/internal/server/domain"
	"net/http"
	"time"
)

// Service - jwt service
type Service struct {
	secret     string
	expire     time.Duration
	signMethod jwt.SigningMethod
	tokenName  string
}

type claims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"user_id,omitempty"`
}

// NewJwtService - factory
func NewJwtService(cfg *Config) *Service {
	return &Service{
		secret:     cfg.Secret,
		expire:     cfg.Expire,
		signMethod: jwt.SigningMethodHS256,
		tokenName:  cfg.Name,
	}
}

// LoadUserID - load token from header, parse token, return user
func (s *Service) LoadUserID(r *http.Request) (domain.UserID, error) {
	cToken, err := r.Cookie(s.tokenName)
	if err != nil {
		return 0, errors.New("user session not' start")
	}

	claim, err := s.loadClaims(cToken.Value)
	if err != nil {
		return 0, errors.New("token not valid")
	}

	return domain.UserID(claim.UserID), nil
}

// CreateSession - По хорошему надо создать структуру сессии и возвращать ее,
// для структуры методы save, read через интефрейс репозитория и там уже Cookie итд
// тут упрощенно
func (s *Service) CreateSession(w http.ResponseWriter, r *http.Request, id domain.UserID) error {
	token, err := s.buildJWTString(id)
	if err != nil {
		return err
	}

	w.Header().Set(s.tokenName, token)

	return nil
}

func (s *Service) buildJWTString(id domain.UserID) (string, error) {
	token := jwt.NewWithClaims(s.signMethod, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expire)),
		},
		UserID: int64(id),
	})

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) loadClaims(tokenStr string) (*claims, error) {
	cl := &claims{}
	token, err := jwt.ParseWithClaims(tokenStr, cl,
		func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != s.signMethod.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
			}

			return []byte(s.secret), nil
		})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token not valid")
	}

	return cl, nil
}
