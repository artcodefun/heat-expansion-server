package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Server owns the HTTP server lifecycle.
type Server struct {
	srv *http.Server
}

func NewServer(engine *gin.Engine, addr string) *Server {
	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: engine,
		},
	}
}

func (s *Server) Start() error {
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

// HealthHandler reports service liveness.
func HealthHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
