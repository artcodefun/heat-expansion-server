-- Down migration: drops base item tables in reverse order

DROP TABLE game.base_storage_items;

DROP TABLE game.base_tech_items;

DROP TABLE game.base_build_items;

DROP TABLE game.base_army_items;
