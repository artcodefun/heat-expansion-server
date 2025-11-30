-- Army items lifecycle queries

-- name: ListNewArmyItems :many
SELECT p.id, p.name, p.category, p.unlock_technology_id, p.short_description, p.full_description, p.price, p.production_time, p.space, p.image_url, p.attack, p.defence, p.capacity, p.stealth, p.speed
FROM army_item_prototypes p
WHERE p.category = $2 AND NOT EXISTS (
    SELECT 1 FROM base_army_items bai WHERE bai.base_id = $1 AND bai.prototype_id = p.id
);

-- name: ListPendingArmyItems :many
SELECT bai.id, bai.base_id, bai.prototype_id, bai.status, bai.pending_data, p.id AS proto_id, p.name, p.category, p.unlock_technology_id, p.short_description, p.full_description, p.price, p.production_time, p.space, p.image_url, p.attack, p.defence, p.capacity, p.stealth, p.speed
FROM base_army_items bai
JOIN army_item_prototypes p ON p.id = bai.prototype_id
WHERE bai.base_id = $1 AND p.category = $2 AND bai.status = 'PENDING'
ORDER BY p.id;

-- name: ListInProductionArmyItems :many
SELECT bai.id, bai.base_id, bai.prototype_id, bai.status, bai.in_prod_data, p.id AS proto_id, p.name, p.category, p.unlock_technology_id, p.short_description, p.full_description, p.price, p.production_time, p.space, p.image_url, p.attack, p.defence, p.capacity, p.stealth, p.speed
FROM base_army_items bai
JOIN army_item_prototypes p ON p.id = bai.prototype_id
WHERE bai.base_id = $1 AND p.category = $2 AND bai.status = 'IN_PRODUCTION'
ORDER BY (bai.in_prod_data->>'completion_date')::bigint ASC NULLS LAST;

-- name: ListPresentArmyItems :many
SELECT bai.id, bai.base_id, bai.prototype_id, bai.status, bai.present_data, p.id AS proto_id, p.name, p.category, p.unlock_technology_id, p.short_description, p.full_description, p.price, p.production_time, p.space, p.image_url, p.attack, p.defence, p.capacity, p.stealth, p.speed
FROM base_army_items bai
JOIN army_item_prototypes p ON p.id = bai.prototype_id
WHERE bai.base_id = $1 AND p.category = $2 AND bai.status = 'PRESENT'
ORDER BY p.id;

-- name: ListPendingArmyItemsAll :many
SELECT bai.id, bai.base_id, bai.prototype_id, bai.status, bai.pending_data, p.id AS proto_id, p.name, p.category, p.unlock_technology_id, p.short_description, p.full_description, p.price, p.production_time, p.space, p.image_url, p.attack, p.defence, p.capacity, p.stealth, p.speed
FROM base_army_items bai
JOIN army_item_prototypes p ON p.id = bai.prototype_id
WHERE bai.base_id = $1 AND bai.status = 'PENDING'
ORDER BY p.id;

-- name: ListInProductionArmyItemsAll :many
SELECT bai.id, bai.base_id, bai.prototype_id, bai.status, bai.in_prod_data, p.id AS proto_id, p.name, p.category, p.unlock_technology_id, p.short_description, p.full_description, p.price, p.production_time, p.space, p.image_url, p.attack, p.defence, p.capacity, p.stealth, p.speed
FROM base_army_items bai
JOIN army_item_prototypes p ON p.id = bai.prototype_id
WHERE bai.base_id = $1 AND bai.status = 'IN_PRODUCTION'
ORDER BY (bai.in_prod_data->>'completion_date')::bigint ASC NULLS LAST;

-- name: ListPresentArmyItemsAll :many
SELECT bai.id, bai.base_id, bai.prototype_id, bai.status, bai.present_data, p.id AS proto_id, p.name, p.category, p.unlock_technology_id, p.short_description, p.full_description, p.price, p.production_time, p.space, p.image_url, p.attack, p.defence, p.capacity, p.stealth, p.speed
FROM base_army_items bai
JOIN army_item_prototypes p ON p.id = bai.prototype_id
WHERE bai.base_id = $1 AND bai.status = 'PRESENT'
ORDER BY p.id;
