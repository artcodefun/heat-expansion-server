package bootstrap

import (
	"database/sql"
	"fmt"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	dbgen "github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/db/repo"
	readstoregen "github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/readstore/gen"
	readrepo "github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/readstore/repo"
	infraevents "github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/events"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/i18n"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/i18n/locales"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/payment"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/security"
)

type Adapters struct {
	Orders            ports.PurchaseOrderRepository
	Packages          ports.CrystalPackageRepository
	PackageRead       ports.PackageReadRepository
	OrderRead         ports.OrderReadRepository
	Gateway           ports.PaymentGateway
	Outbox            ports.OutboxEventRepository
	IntegrationOutbox ports.IntegrationOutboxRepository
	TxMgr             ports.TransactionManager
	Events            ports.EventPublisher
	IntegrationEvents ports.IntegrationEventPublisher
	Tokens            ports.TokenValidator
	Translator        ports.Translator
}

func NewAdapters(db *sql.DB, jwtSecret, rabbitURL, intExchange, yookassaShopID, yookassaSecretKey string) (*Adapters, error) {
	q := dbgen.New(db)
	rq := readstoregen.New(db)

	intPublisher, err := infraevents.NewRabbitMQPublisher(rabbitURL, intExchange)
	if err != nil {
		return nil, fmt.Errorf("billing: failed to initialize RabbitMQ publisher: %w", err)
	}

	translator := i18n.NewJSONTranslator()
	if err := translator.LoadFromFS(locales.Files, "."); err != nil {
		return nil, fmt.Errorf("billing: failed to load translations: %w", err)
	}

	return &Adapters{
		Orders:            repo.NewOrderRepo(q),
		Packages:          repo.NewPackageRepo(q),
		PackageRead:       readrepo.NewPackageReadRepo(rq),
		OrderRead:         readrepo.NewOrderReadRepo(rq),
		Gateway:           payment.NewYooKassaGateway(yookassaShopID, yookassaSecretKey),
		Outbox:            repo.NewOutboxEventRepo(q),
		IntegrationOutbox: repo.NewIntegrationOutboxRepo(q),
		TxMgr:             repo.NewDBTxManager(db),
		Events:            infraevents.NewSimplePublisher(),
		IntegrationEvents: intPublisher,
		Tokens:            security.NewSimpleTokenValidator(jwtSecret),
		Translator:        translator,
	}, nil
}
