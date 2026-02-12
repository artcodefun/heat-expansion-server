-- Down migration: drops prototypes tables in reverse order

DROP TABLE game.storage_item_prototypes;

DROP TABLE game.build_item_prototypes;

DROP TABLE game.army_item_prototypes;

DROP TABLE game.tech_item_prototypes;
