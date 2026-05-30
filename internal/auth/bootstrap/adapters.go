package bootstrap

import (
	"database/sql"
	"fmt"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/db/repo"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/email"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/events"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/i18n"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/i18n/locales"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/security"
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

func NewAdapters(db *sql.DB, jwtSecret string, intPublisher ports.IntegrationEventPublisher, smtpCfg SMTPConfig) (*Adapters, error) {
	translator := i18n.NewJSONTranslator()
	if err := translator.LoadFromFS(locales.Files, "."); err != nil {
		return nil, fmt.Errorf("failed to load translations: %w", err)
	}

	q := gen.New(db)

	return &Adapters{
		Repo:              repo.NewAccountRepository(q),
		Hasher:            security.NewBcryptHasher(),
		TokenProvider:     security.NewSimpleTokenProvider(jwtSecret),
		Outbox:            repo.NewOutboxEventRepo(q),
		TxMgr:             repo.NewDBTxManager(db),
		Events:            events.NewSimplePublisher(),
		IntegrationOutbox: repo.NewIntegrationOutboxRepo(q),
		IntegrationEvents: intPublisher,
		Translator:        translator,
		ResetRepo:         repo.NewPasswordResetRepository(q),
		EmailSender:       email.NewSMTPSender(smtpCfg.Host, smtpCfg.User, smtpCfg.Password, smtpCfg.From),
	}, nil
}
