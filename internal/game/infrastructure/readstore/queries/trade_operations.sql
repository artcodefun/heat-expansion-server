-- Trade operation read queries

-- name: GetTradeOperation :one
SELECT *
FROM game.trade_operations
WHERE id = $1;

-- name: ListActiveTradeOperations :many
SELECT *
FROM game.trade_operations
WHERE (sender_base_id = $1 OR receiver_base_id = $1)
  AND phase != 'COMPLETED'
ORDER BY id DESC;