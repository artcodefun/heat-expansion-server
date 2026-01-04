-- Prototypes queries for read-store

-- name: ListArmyPrototypes :many
SELECT id, name, category, unlock_technology_id, short_description, full_description,
       price,
       production_time, space, image_url,
       attack, defence, capacity, stealth, speed
FROM army_item_prototypes
ORDER BY id;

-- name: ListBuildPrototypes :many
SELECT id, name, category, unlock_technology_id, short_description, full_description,
       price,
       production_time, space, image_url,
       control_data, resources_data, defense_data, military_data, intelligence_data
FROM build_item_prototypes
ORDER BY id;

-- name: ListStoragePrototypes :many
SELECT id, name, category, short_description, full_description, image_url,
       buff_data, map_data, damaged_data, artifact_data, consumable_data
FROM storage_item_prototypes
ORDER BY id;

-- name: ListTechPrototypes :many
SELECT id, name, category, unlock_technology_id, short_description, full_description,
       price,
       research_time, image_url, effects
FROM tech_item_prototypes
ORDER BY id;
