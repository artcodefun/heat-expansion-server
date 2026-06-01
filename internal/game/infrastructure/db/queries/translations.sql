-- name: GetAllTranslations :many
SELECT * FROM game.translations ORDER BY locale, key;
