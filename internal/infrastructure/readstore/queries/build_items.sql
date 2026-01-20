-- Building items lifecycle queries

-- name: ListBuildPrototypesByIDs :many
SELECT p.id, p.name, p.category, p.faction, p.unlock_technology_id, p.short_description, p.full_description, p.price, p.production_time, p.space, p.image_url, p.control_data, p.resources_data, p.defense_data, p.military_data, p.intelligence_data
FROM build_item_prototypes p
WHERE p.id = ANY($1::bigint[])
ORDER BY p.id;

-- name: ListPendingBuildItems :many
SELECT bbi.id, bbi.base_id, bbi.prototype_id, bbi.status, bbi.pending_data, p.id AS proto_id, p.name, p.category, p.faction, p.unlock_technology_id, p.short_description, p.full_description, p.price, p.production_time, p.space, p.image_url, p.control_data, p.resources_data, p.defense_data, p.military_data, p.intelligence_data
FROM base_build_items bbi
JOIN build_item_prototypes p ON p.id = bbi.prototype_id
WHERE bbi.base_id = $1 AND p.category = $2 AND bbi.status = 'PENDING'
ORDER BY p.id;

-- name: ListInProductionBuildItems :many
SELECT bbi.id, bbi.base_id, bbi.prototype_id, bbi.status, bbi.in_prod_data, p.id AS proto_id, p.name, p.category, p.faction, p.unlock_technology_id, p.short_description, p.full_description, p.price, p.production_time, p.space, p.image_url, p.control_data, p.resources_data, p.defense_data, p.military_data, p.intelligence_data
FROM base_build_items bbi
JOIN build_item_prototypes p ON p.id = bbi.prototype_id
WHERE bbi.base_id = $1 AND p.category = $2 AND bbi.status = 'IN_PRODUCTION'
ORDER BY (bbi.in_prod_data->>'completion_date')::bigint ASC NULLS LAST;

-- name: ListPresentBuildItems :many
SELECT bbi.id, bbi.base_id, bbi.prototype_id, bbi.status, bbi.present_data, p.id AS proto_id, p.name, p.category, p.faction, p.unlock_technology_id, p.short_description, p.full_description, p.price, p.production_time, p.space, p.image_url, p.control_data, p.resources_data, p.defense_data, p.military_data, p.intelligence_data
FROM base_build_items bbi
JOIN build_item_prototypes p ON p.id = bbi.prototype_id
WHERE bbi.base_id = $1 AND p.category = $2 AND bbi.status = 'PRESENT'
ORDER BY p.id;

-- name: ListPendingBuildItemsAll :many
SELECT bbi.id, bbi.base_id, bbi.prototype_id, bbi.status, bbi.pending_data, p.id AS proto_id, p.name, p.category, p.faction, p.unlock_technology_id, p.short_description, p.full_description, p.price, p.production_time, p.space, p.image_url, p.control_data, p.resources_data, p.defense_data, p.military_data, p.intelligence_data
FROM base_build_items bbi
JOIN build_item_prototypes p ON p.id = bbi.prototype_id
WHERE bbi.base_id = $1 AND bbi.status = 'PENDING'
ORDER BY p.id;

-- name: ListInProductionBuildItemsAll :many
SELECT bbi.id, bbi.base_id, bbi.prototype_id, bbi.status, bbi.in_prod_data, p.id AS proto_id, p.name, p.category, p.faction, p.unlock_technology_id, p.short_description, p.full_description, p.price, p.production_time, p.space, p.image_url, p.control_data, p.resources_data, p.defense_data, p.military_data, p.intelligence_data
FROM base_build_items bbi
JOIN build_item_prototypes p ON p.id = bbi.prototype_id
WHERE bbi.base_id = $1 AND bbi.status = 'IN_PRODUCTION'
ORDER BY (bbi.in_prod_data->>'completion_date')::bigint ASC NULLS LAST;

-- name: ListPresentBuildItemsAll :many
SELECT bbi.id, bbi.base_id, bbi.prototype_id, bbi.status, bbi.present_data, p.id AS proto_id, p.name, p.category, p.faction, p.unlock_technology_id, p.short_description, p.full_description, p.price, p.production_time, p.space, p.image_url, p.control_data, p.resources_data, p.defense_data, p.military_data, p.intelligence_data
FROM base_build_items bbi
JOIN build_item_prototypes p ON p.id = bbi.prototype_id
WHERE bbi.base_id = $1 AND bbi.status = 'PRESENT'
ORDER BY p.id;
