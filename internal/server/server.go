package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend-store/config"
	"backend-store/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg    *config.Config
	log    logger.Log
	router *gin.Engine
	server *http.Server
}

func New(cfg *config.Config, log logger.Log, router *gin.Engine) *Server {
	return &Server{
		cfg:    cfg,
		log:    log,
		router: router,
	}
}

// Start запускает HTTP сервер
func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:         s.cfg.ServerHost + ":" + s.cfg.ServerPort,
		Handler:      s.router,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
	}

	// Запускаем сервер в горутине
	go func() {
		s.log.Info("Server starting on", "address", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatal("Failed to start server:", err)
		}
	}()

	return nil
}

// WaitForShutdown ожидает сигналов завершения работы
func (s *Server) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Fatal("Server forced to shutdown:", err)
	}

	s.log.Info("Server exited")
}

// Stop принудительно останавливает сервер
func (s *Server) Stop() {
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.server.Shutdown(ctx); err != nil {
			s.log.Error("Error stopping server:", err)
		}
	}
}
