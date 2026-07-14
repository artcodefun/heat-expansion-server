// Package interceptor holds gRPC server interceptors for the billing private API.
package interceptor

import (
	"context"
	"crypto/subtle"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// InternalKeyHeader is the metadata key carrying the shared static module key.
// gRPC lowercases metadata keys, so it must be lowercase here.
const InternalKeyHeader = "x-internal-key"

// KeyAuth returns a unary interceptor that authenticates callers via constant-time
// comparison of the x-internal-key metadata value against the configured key.
func KeyAuth(expectedKey string) grpc.UnaryServerInterceptor {
	expected := []byte(expectedKey)
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing request metadata")
		}
		vals := md.Get(InternalKeyHeader)
		if len(vals) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing internal key")
		}
		if subtle.ConstantTimeCompare([]byte(vals[0]), expected) != 1 {
			return nil, status.Error(codes.Unauthenticated, "invalid internal key")
		}
		return handler(ctx, req)
	}
}
