package http

import "github.com/gin-gonic/gin"

// Server owns the HTTP server lifecycle.
type Server struct {
	Engine *gin.Engine
}

func NewServer(engine *gin.Engine) *Server {
	return &Server{Engine: engine}
}

func (s *Server) Start(addr string) error {
	return s.Engine.Run(addr)
}

// HealthHandler reports service liveness.
func HealthHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
