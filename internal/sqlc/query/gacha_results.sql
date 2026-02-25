-- name: InsertGachaResult :exec
INSERT INTO gacha_results (
  result_id, user_id, instance_id, card_id, kind, rarity, is_pickup, created_at
) VALUES (
  sqlc.arg(result_id), sqlc.arg(user_id), sqlc.arg(instance_id), sqlc.arg(card_id),
  sqlc.arg(kind), sqlc.arg(rarity), sqlc.arg(is_pickup), sqlc.arg(created_at)
);