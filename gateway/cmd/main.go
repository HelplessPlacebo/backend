package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/HelplessPlacebo/backend/gateway/config"
	"github.com/HelplessPlacebo/backend/gateway/internal/proxy"
	"github.com/HelplessPlacebo/backend/gateway/internal/router"
	"github.com/HelplessPlacebo/backend/pkg/shared"
)

func main() {
	logger := shared.NewLogger()
	cfg := config.Load()

	// build auth base url using config.AuthBaseURL()
	p := proxy.NewClient(cfg.AuthBaseURL())

	r := router.NewRouter(p, logger)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// graceful shutdown
	done := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		logger.Infof("shutting down gateway")
		_ = srv.Shutdown(ctx)
		close(done)
	}()

	logger.Infof("gateway running on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Errorf("server error: %v", err)
	}

	<-done
}
