package bootstrap

import (
	"database/sql"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
	dbgen "github.com/artcodefun/heat-expansion-server/internal/admin/infrastructure/db/gen"
	dbrepo "github.com/artcodefun/heat-expansion-server/internal/admin/infrastructure/db/repo"
	"github.com/artcodefun/heat-expansion-server/internal/admin/infrastructure/i18n"
	readgen "github.com/artcodefun/heat-expansion-server/internal/admin/infrastructure/readstore/gen"
	readrepo "github.com/artcodefun/heat-expansion-server/internal/admin/infrastructure/readstore/repo"
	"github.com/artcodefun/heat-expansion-server/internal/admin/infrastructure/security"
	platformsecurity "github.com/artcodefun/heat-expansion-server/internal/platform/security"
)

// Adapters wires all secondary adapters for the admin module.
type Adapters struct {
	Admins           ports.AdminRepository
	Sessions         ports.SessionRepository
	AdminRead        ports.AdminReadRepository
	SessionValidator ports.SessionValidator
	TxMgr            ports.TransactionManager
	Hasher           ports.PasswordHasher
	TokenGen         ports.SessionTokenGenerator
	Translator       ports.Translator
}

func NewAdapters(db *sql.DB) (*Adapters, error) {
	q := dbgen.New(db)
	rq := readgen.New(db)

	translator, err := i18n.NewSimpleTranslator()
	if err != nil {
		return nil, err
	}

	sessions := dbrepo.NewSessionRepository(q)
	admins := dbrepo.NewAdminRepository(q)

	return &Adapters{
		Admins:           admins,
		Sessions:         sessions,
		AdminRead:        readrepo.NewAdminReadRepo(rq),
		SessionValidator: security.NewAdminSessionValidator(sessions, admins),
		TxMgr:            dbrepo.NewDBTxManager(db),
		Hasher:           platformsecurity.NewBcryptHasher(),
		TokenGen:         security.NewRandomSessionTokenGenerator(),
		Translator:       translator,
	}, nil
}
