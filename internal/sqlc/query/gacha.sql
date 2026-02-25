-- name: GetRandomCharacterByRarity :one
SELECT card_id
FROM cards
WHERE card_kind = 1
  AND rarity = sqlc.arg(rarity)
ORDER BY random()
LIMIT 1;

-- name: GetRandomEquipByRarity :one
SELECT card_id
FROM cards
WHERE card_kind = 0
  AND rarity = sqlc.arg(rarity)
ORDER BY random()
LIMIT 1;