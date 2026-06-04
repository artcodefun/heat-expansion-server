package bootstrap

import (
	"database/sql"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/auth/domain"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/db/repo"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/email"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/i18n"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/security"
	platformevents "github.com/artcodefun/heat-expansion-server/internal/platform/events"
	platformsecurity "github.com/artcodefun/heat-expansion-server/internal/platform/security"
)

type Adapters struct {
	Repo              ports.AccountRepository
	Hasher            ports.PasswordHasher
	TokenProvider     ports.TokenProvider
	Outbox            ports.OutboxEventRepository
	TxMgr             ports.TransactionManager
	Events            ports.EventPublisher
	IntegrationOutbox ports.IntegrationOutboxRepository
	IntegrationEvents ports.IntegrationEventPublisher
	Translator        ports.Translator
	ResetRepo         ports.PasswordResetRepository
	EmailSender       ports.EmailSender
}

type SMTPConfig struct {
	Host     string
	User     string
	Password string
	From     string
}

func NewAdapters(db *sql.DB, jwtPrivateKeyPEM string, intPublisher ports.IntegrationEventPublisher, smtpCfg SMTPConfig) (*Adapters, error) {
	translator, err := i18n.NewSimpleTranslator()
	if err != nil {
		return nil, err
	}

	tokenProvider, err := security.NewSimpleTokenProvider(jwtPrivateKeyPEM)
	if err != nil {
		return nil, err
	}

	q := gen.New(db)

	return &Adapters{
		Repo:              repo.NewAccountRepository(q),
		Hasher:            platformsecurity.NewBcryptHasher(),
		TokenProvider:     tokenProvider,
		Outbox:            repo.NewOutboxEventRepo(q),
		TxMgr:             repo.NewDBTxManager(db),
		Events:            platformevents.NewSimplePublisher[domain.DomainEvent](),
		IntegrationOutbox: repo.NewIntegrationOutboxRepo(q),
		IntegrationEvents: intPublisher,
		Translator:        translator,
		ResetRepo:         repo.NewPasswordResetRepository(q),
		EmailSender:       email.NewSMTPSender(smtpCfg.Host, smtpCfg.User, smtpCfg.Password, smtpCfg.From),
	}, nil
}
