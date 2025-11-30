-- Technology prototypes queries

-- name: GetTechPrototypeByID :one
SELECT id, name, category, unlock_technology_id, short_description, full_description,
       price,
       research_time, image_url, effects
FROM tech_item_prototypes
WHERE id = @id;

-- name: ListTechPrototypes :many
SELECT id, name, category, unlock_technology_id, short_description, full_description,
       price,
       research_time, image_url, effects
FROM tech_item_prototypes
ORDER BY id;
