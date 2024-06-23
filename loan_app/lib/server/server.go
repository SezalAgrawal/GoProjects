package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goProjects/loan_app/lib/logger"
	"go.uber.org/zap"
)

type GracefulShutdownServer struct {
	*http.Server
}

func (s *GracefulShutdownServer) Run() error {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	ctx := context.Background()

	go func() {
		// start server in go routine and then block on interrupt signal
		ctx := context.Background()
		logger.I(ctx, "server staring up", zap.String("addr", s.Addr))
		if err := s.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logger.E(ctx, "unable to start server", zap.Error(err))
				done <- os.Kill
				return
			}
			logger.I(ctx, "server shutdown")
		}
	}()

	sig := <-done
	logger.I(ctx, "server shutting down", zap.String("sig", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		logger.E(ctx, "server failed to shutdown gracefully", zap.Error(err))
		return err
	}

	return nil
}
