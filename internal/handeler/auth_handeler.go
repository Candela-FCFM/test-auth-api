package handeler

import (
	"net/http"

	"github.com/alejandrogzzcandela/auth-api/internal/service"
)

type IAuthHandeler interface {
	HealthCheck(http.ResponseWriter, *http.Request)
}

type AuthHandeler struct {
	s service.IAuthService
}

func NewAuthHandeler(s service.IAuthService) IAuthHandeler {
	return &AuthHandeler{s: s}
}

func (ah AuthHandeler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if result, err := ah.s.HealthCheck(); result && err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
