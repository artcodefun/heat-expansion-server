-- Building prototypes queries

-- name: GetBuildPrototypeByID :one
SELECT id, name, category, faction, unlock_technology_id, short_description, full_description,
       price,
       production_time, space, image_url,
       control_data, resources_data, defense_data, military_data, intelligence_data
FROM build_item_prototypes
WHERE id = @id;

-- name: ListBuildPrototypes :many
SELECT id, name, category, faction, unlock_technology_id, short_description, full_description,
       price,
       production_time, space, image_url,
       control_data, resources_data, defense_data, military_data, intelligence_data
FROM build_item_prototypes
ORDER BY id;
