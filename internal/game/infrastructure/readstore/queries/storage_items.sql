-- Storage items queries

-- name: ListPresentStorageItems :many
SELECT bsi.id, bsi.base_id, bsi.prototype_id, bsi.status, bsi.present_data, bsi.state, p.id AS proto_id, p.name, p.category, p.estimated_worth, p.short_description, p.full_description, p.image_url, p.buff_data, p.intel_data, p.damaged_data, p.artifact_data, p.consumable_data
FROM game.base_storage_items bsi
JOIN game.storage_item_prototypes p ON p.id = bsi.prototype_id
WHERE bsi.base_id = $1 AND p.category = $2 AND bsi.status = 'PRESENT'
ORDER BY p.id;
