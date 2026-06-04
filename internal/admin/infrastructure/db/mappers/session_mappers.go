package mappers

import (
	"github.com/artcodefun/heat-expansion-server/internal/admin/domain"
	"github.com/artcodefun/heat-expansion-server/internal/admin/infrastructure/db/gen"
)

// SessionFromDB maps a sqlc AdminSession row to a domain.Session.
func SessionFromDB(row gen.AdminSession) *domain.Session {
	return &domain.Session{
		Token:     row.Token,
		AdminID:   row.AdminID,
		ExpiresAt: row.ExpiresAt,
		CreatedAt: row.CreatedAt,
	}
}

// CreateSessionParamsFromDomain maps a domain.Session to sqlc insert params.
func CreateSessionParamsFromDomain(session *domain.Session) gen.CreateSessionParams {
	return gen.CreateSessionParams{
		Token:     session.Token,
		AdminID:   session.AdminID,
		ExpiresAt: session.ExpiresAt,
		CreatedAt: session.CreatedAt,
	}
}
