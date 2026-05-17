ALTER TABLE game.build_item_prototypes
    ADD COLUMN creation_sources JSONB NOT NULL DEFAULT '["PLAYER_BASE"]'::jsonb;

UPDATE game.build_item_prototypes
SET creation_sources = CASE
    WHEN faction = 'EXO_COALITION' THEN '["PLAYER_BASE"]'::jsonb
    ELSE '["NPC_LOCATION"]'::jsonb
END;

ALTER TABLE game.army_item_prototypes
    ADD COLUMN creation_sources JSONB NOT NULL DEFAULT '["PLAYER_BASE"]'::jsonb;

UPDATE game.army_item_prototypes
SET creation_sources = CASE
    WHEN faction = 'EXO_COALITION' THEN '["PLAYER_BASE"]'::jsonb
    ELSE '["NPC_LOCATION"]'::jsonb
END;

ALTER TABLE game.storage_item_prototypes
    ADD COLUMN creation_sources JSONB NOT NULL DEFAULT '["NPC_LOCATION", "CONSUMABLE_BOX"]'::jsonb;
