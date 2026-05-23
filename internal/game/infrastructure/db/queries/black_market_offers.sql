-- Black market offers queries

-- name: GetBlackMarketOfferByID :one
SELECT *
FROM game.black_market_offers
WHERE id = @id;

-- name: GetBlackMarketOfferByIDForUpdate :one
SELECT *
FROM game.black_market_offers
WHERE id = @id
FOR UPDATE;

-- name: ListActiveLimitedBlackMarketOffers :many
SELECT *
FROM game.black_market_offers
WHERE is_limited = TRUE
    AND ends_at > @now
ORDER BY priority DESC, id ASC;

-- name: ListExpiredLimitedBlackMarketOffers :many
SELECT *
FROM game.black_market_offers
WHERE is_limited = TRUE
    AND ends_at <= @now
ORDER BY ends_at ASC, priority DESC, id ASC;

-- name: InsertBlackMarketOffer :one
INSERT INTO game.black_market_offers (
    kind,
    prototype_id,
    price_in_crystals,
    ends_at,
    is_limited,
    priority
) VALUES (
    @kind,
    @prototype_id,
    @price_in_crystals,
    @ends_at,
    @is_limited,
    @priority
)
RETURNING *;

-- name: UpdateBlackMarketOffer :one
UPDATE game.black_market_offers
SET kind = @kind,
    prototype_id = @prototype_id,
    price_in_crystals = @price_in_crystals,
    ends_at = @ends_at,
    is_limited = @is_limited,
    priority = @priority
WHERE id = @id
RETURNING *;