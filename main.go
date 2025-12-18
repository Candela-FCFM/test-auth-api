package main

import (
	"net/http"

	"github.com/alejandrogzzcandela/auth-api/internal/handeler"
	"github.com/alejandrogzzcandela/auth-api/internal/middleware"
	"github.com/alejandrogzzcandela/auth-api/internal/repository"
	"github.com/alejandrogzzcandela/auth-api/internal/service"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	repo := repository.NewAuthRepository()
	serv := service.NewAuthService(repo)
	hand := handeler.NewAuthHandeler(serv)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", hand.HealthCheck)

	wrappedMux := middleware.ObservabilityMiddleware(logger)(mux)

	http.ListenAndServe("internal.learning-lenguages:8080", wrappedMux)
}
