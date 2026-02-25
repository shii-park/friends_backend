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

-- name: ListGachaCharacters :many
SELECT
  c.card_id,
  c.card_name,
  c.card_kind,
  c.rarity,
  c.card_icon_url,
  ch.character_id,
  ch.hp, ch.atk, ch.tech,
  ch.init_hp, ch.init_atk, ch.init_tech,
  ch.max_hp, ch.max_atk, ch.max_tech,
  ch.special_type
FROM cards c
JOIN characters ch ON ch.card_id = c.card_id
WHERE c.card_kind = 1
ORDER BY c.card_id ASC;