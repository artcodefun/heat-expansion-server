-- Local projection of auth accounts, populated by consuming
-- auth.account.registered.v1 integration events. Billing keeps it so it can
-- attach a customer email to YooKassa payment receipts (54-FZ).
CREATE TABLE billing.users (
    id    UUID PRIMARY KEY,
    email TEXT NOT NULL
);
