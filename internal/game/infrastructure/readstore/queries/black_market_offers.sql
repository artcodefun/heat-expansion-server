-- Black Market offer queries

-- name: ListActiveBlackMarketBuildingOffers :many
SELECT o.id, o.kind, o.price_in_crystals, o.ends_at, o.is_limited, o.priority,
       p.id AS proto_id, p.name, p.category, p.faction, p.unlock_technology_id, p.short_description, p.full_description,
       p.price, p.production_time, p.space, p.image_url, p.control_data, p.resources_data, p.defense_data, p.military_data, p.intelligence_data
FROM game.black_market_offers o
JOIN game.build_item_prototypes p ON p.id = o.prototype_id
WHERE o.kind = 'BUILDING'
  AND (sqlc.narg(limited) IS NULL OR o.is_limited = sqlc.narg(limited))
  AND (NOT o.is_limited OR o.ends_at > @now)
  AND p.creation_sources @> '["BLACK_MARKET"]'::jsonb
ORDER BY o.priority DESC, o.id ASC;

-- name: ListActiveBlackMarketArmyOffers :many
SELECT o.id, o.kind, o.price_in_crystals, o.ends_at, o.is_limited, o.priority,
       p.id AS proto_id, p.name, p.category, p.faction, p.unlock_technology_id, p.short_description, p.full_description,
       p.price, p.production_time, p.space, p.image_url, p.attack, p.defence, p.capacity, p.stealth, p.speed
FROM game.black_market_offers o
JOIN game.army_item_prototypes p ON p.id = o.prototype_id
WHERE o.kind = 'ARMY'
  AND (sqlc.narg(limited) IS NULL OR o.is_limited = sqlc.narg(limited))
  AND (NOT o.is_limited OR o.ends_at > @now)
  AND p.creation_sources @> '["BLACK_MARKET"]'::jsonb
ORDER BY o.priority DESC, o.id ASC;

-- name: ListActiveBlackMarketStorageOffers :many
SELECT o.id, o.kind, o.price_in_crystals, o.ends_at, o.is_limited, o.priority,
       p.id AS proto_id, p.name, p.category, p.estimated_worth, p.short_description, p.full_description,
       p.image_url, p.buff_data, p.intel_data, p.damaged_data, p.artifact_data, p.consumable_data
FROM game.black_market_offers o
JOIN game.storage_item_prototypes p ON p.id = o.prototype_id
WHERE o.kind = 'STORAGE'
  AND (sqlc.narg(limited) IS NULL OR o.is_limited = sqlc.narg(limited))
  AND (NOT o.is_limited OR o.ends_at > @now)
  AND p.creation_sources @> '["BLACK_MARKET"]'::jsonb
ORDER BY o.priority DESC, o.id ASC;