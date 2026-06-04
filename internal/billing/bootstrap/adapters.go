package bootstrap

import (
	"database/sql"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
	dbgen "github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/db/repo"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/i18n"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/payment"
	readstoregen "github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/readstore/gen"
	readrepo "github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/readstore/repo"
	platformevents "github.com/artcodefun/heat-expansion-server/internal/platform/events"
	"github.com/artcodefun/heat-expansion-server/internal/platform/security"
)

type Adapters struct {
	Orders            ports.PurchaseOrderRepository
	Packages          ports.CrystalPackageRepository
	Users             ports.UserRepository
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

func NewAdapters(db *sql.DB, jwtPublicKeyPEM string, intPublisher ports.IntegrationEventPublisher, yookassaShopID, yookassaSecretKey string) (*Adapters, error) {
	q := dbgen.New(db)
	rq := readstoregen.New(db)

	translator, err := i18n.NewSimpleTranslator()
	if err != nil {
		return nil, err
	}

	tokens, err := security.NewSimpleTokenValidator(jwtPublicKeyPEM)
	if err != nil {
		return nil, err
	}

	return &Adapters{
		Orders:            repo.NewOrderRepo(q),
		Packages:          repo.NewPackageRepo(q),
		Users:             repo.NewUserRepo(q),
		PackageRead:       readrepo.NewPackageReadRepo(rq),
		OrderRead:         readrepo.NewOrderReadRepo(rq),
		Gateway:           payment.NewYooKassaGateway(yookassaShopID, yookassaSecretKey),
		Outbox:            repo.NewOutboxEventRepo(q),
		IntegrationOutbox: repo.NewIntegrationOutboxRepo(q),
		TxMgr:             repo.NewDBTxManager(db),
		Events:            platformevents.NewSimplePublisher[domain.DomainEvent](),
		IntegrationEvents: intPublisher,
		Tokens:            tokens,
		Translator:        translator,
	}, nil
}
