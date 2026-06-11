-- name: GetAllTranslations :many
SELECT * FROM game.translations ORDER BY locale, key;

-- name: UpsertTranslation :exec
INSERT INTO game.translations (key, locale, value)
VALUES (@key, @locale, @value)
ON CONFLICT (key, locale) DO UPDATE SET value = EXCLUDED.value;

-- name: NotifyTranslationsChanged :exec
NOTIFY game_translations_changed;
