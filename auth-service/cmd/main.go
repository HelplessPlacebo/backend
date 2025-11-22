package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/HelplessPlacebo/backend/auth-service/config"
	"github.com/HelplessPlacebo/backend/auth-service/internal/api"
	"github.com/HelplessPlacebo/backend/auth-service/internal/service"
	"github.com/HelplessPlacebo/backend/auth-service/internal/storage"
	"github.com/HelplessPlacebo/backend/pkg/shared"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func main() {
	logger := shared.NewLogger()

	cfg := config.Load()

	db := storage.NewPostgres(cfg.DBURL)

	userRepo := storage.NewUserRepo(db)
	authSvc := service.NewAuthService(userRepo, logger)

	r := chi.NewRouter()
	v := validator.New()

	r.Route(cfg.APIBase, func(r chi.Router) {
		api.RegisterRegistration(r, authSvc, v, logger)
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		logger.Infof("shutting down auth-service")
		if err := srv.Shutdown(ctx); err != nil {
			logger.Errorf("server shutdown error: %v", err)
		}
		close(idleConnsClosed)
	}()

	logger.Infof("auth-service running on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Errorf("server error: %v", err)
	}

	<-idleConnsClosed
}
