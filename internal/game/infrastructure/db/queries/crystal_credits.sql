-- Crystal credits queries

-- name: InsertCrystalCredit :exec
INSERT INTO game.crystal_credits (order_id, user_id, crystals, credited_at)
VALUES (@order_id, @user_id, @crystals, @credited_at);

-- name: CrystalCreditExists :one
SELECT EXISTS (
    SELECT 1 FROM game.crystal_credits WHERE order_id = @order_id
) AS exists;
