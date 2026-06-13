package dtos

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
)

const defaultLocale = "en"

// LocaleFromContext extracts the caller's locale from accept-language metadata. Defaults to English.
func LocaleFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return defaultLocale
	}
	vals := md.Get("accept-language")
	if len(vals) == 0 || vals[0] == "" {
		return defaultLocale
	}
	return strings.Split(vals[0], ",")[0]
}

// StatusFromError maps a billing domain/CQRS error to a gRPC status with a translated message.
func StatusFromError(ctx context.Context, tr ports.Translator, err error) error {
	if err == nil {
		return nil
	}
	locale := LocaleFromContext(ctx)

	var appErr cqrs.AppError
	if errors.As(err, &appErr) {
		return status.Error(codeForKind(appErr.Kind), tr.T(locale, appErr.Code, appErr.Params))
	}

	slog.ErrorContext(ctx, "internal error occurred", "error", err.Error())
	return status.Error(codes.Internal, tr.T(locale, "error.application.internal_server_error", nil))
}

func codeForKind(kind cqrs.ErrorKind) codes.Code {
	switch kind {
	case cqrs.KindNotFound:
		return codes.NotFound
	case cqrs.KindForbidden:
		return codes.PermissionDenied
	case cqrs.KindConflict:
		return codes.AlreadyExists
	case cqrs.KindInvalidInput:
		return codes.InvalidArgument
	case cqrs.KindUnavailable:
		return codes.Unavailable
	default:
		return codes.Internal
	}
}
