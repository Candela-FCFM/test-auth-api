package service

import (
	"errors"

	"github.com/alejandrogzzcandela/auth-api/internal/repository"
)

type IAuthService interface {
	HealthCheck() (bool, error)
}

type AuthService struct {
	r repository.IAuthRepository
}

func NewAuthService(r repository.IAuthRepository) IAuthService {
	return &AuthService{r: r}
}

func (s *AuthService) HealthCheck() (bool, error) {
	return true, errors.New("funcion no esta implementada, no tenemos nada que validar")
}
