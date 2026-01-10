-- Technology items lifecycle queries

-- name: ListTechPrototypesByIDs :many
SELECT p.id, p.name, p.category, p.unlock_technology_id, p.short_description, p.full_description, p.price, p.research_time, p.image_url, p.improvement
FROM tech_item_prototypes p
WHERE p.id = ANY($1::bigint[])
ORDER BY p.id;

-- name: ListInResearchTechItems :many
SELECT bti.id, bti.base_id, bti.prototype_id, bti.status, bti.in_progress_data, p.id AS proto_id, p.name, p.category, p.unlock_technology_id, p.short_description, p.full_description, p.price, p.research_time, p.image_url, p.improvement
FROM base_tech_items bti
JOIN tech_item_prototypes p ON p.id = bti.prototype_id
WHERE bti.base_id = $1 AND bti.status = 'IN_PROGRESS'
ORDER BY (bti.in_progress_data->>'completion_date')::bigint ASC NULLS LAST;

-- name: ListDoneTechItems :many
SELECT bti.id, bti.base_id, bti.prototype_id, bti.status, bti.done_data, p.id AS proto_id, p.name, p.category, p.unlock_technology_id, p.short_description, p.full_description, p.price, p.research_time, p.image_url, p.improvement
FROM base_tech_items bti
JOIN tech_item_prototypes p ON p.id = bti.prototype_id
WHERE bti.base_id = $1 AND bti.status = 'DONE'
ORDER BY p.id;
