package grpcclient

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
)

// grpcErrToSentinel maps gRPC status codes to port-level sentinels so the
// application layer can map them to CQRS errors without an infra→cqrs import.
func grpcErrToSentinel(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return err
	}
	switch st.Code() {
	case codes.NotFound:
		return ports.ErrClientNotFound
	case codes.InvalidArgument:
		return ports.ErrClientInvalidInput
	case codes.PermissionDenied, codes.Unauthenticated:
		return ports.ErrClientForbidden
	default:
		return err
	}
}
