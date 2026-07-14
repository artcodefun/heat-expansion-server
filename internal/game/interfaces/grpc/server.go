package grpc

import (
	"context"
	"errors"
	"net"

	"google.golang.org/grpc"
)

// Server wraps a gRPC server with listen/serve lifecycle.
type Server struct {
	srv  *grpc.Server
	addr string
}

// NewServer wraps a configured Router with the listen/serve lifecycle.
func NewServer(router Router, addr string) *Server {
	return &Server{srv: router.srv, addr: addr}
}

// Start binds the listener and serves until Shutdown is called.
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	if err := s.srv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return err
	}
	return nil
}

// Shutdown drains in-flight RPCs gracefully, falling back to a hard Stop if ctx is cancelled.
func (s *Server) Shutdown(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		s.srv.GracefulStop()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		s.srv.Stop()
		return ctx.Err()
	}
}
