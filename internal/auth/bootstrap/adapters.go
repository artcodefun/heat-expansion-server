package bootstrap

import (
	"database/sql"
	"fmt"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/db/repo"
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
}

func NewAdapters(db *sql.DB, jwtSecret string, rabbitURL string, integrationExchange string) (*Adapters, error) {
	intPublisher, err := events.NewRabbitMQPublisher(rabbitURL, integrationExchange)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize RabbitMQ publisher: %w", err)
	}

	translator := i18n.NewJSONTranslator()
	if err := translator.LoadFromFS(locales.Files, "."); err != nil {
		return nil, fmt.Errorf("failed to load translations: %w", err)
	}

	return &Adapters{
		Repo:              repo.NewAccountRepository(db),
		Hasher:            security.NewBcryptHasher(),
		TokenProvider:     security.NewSimpleTokenProvider(jwtSecret),
		Outbox:            repo.NewOutboxEventRepo(db),
		TxMgr:             repo.NewDBTxManager(db),
		Events:            events.NewSimplePublisher(),
		IntegrationOutbox: repo.NewIntegrationOutboxRepo(db),
		IntegrationEvents: intPublisher,
		Translator:        translator,
	}, nil
}
