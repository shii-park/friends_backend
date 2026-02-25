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

-- name: ListGachaEquipments :many
SELECT
  c.card_id,
  c.card_name,
  c.card_kind,
  c.rarity,
  c.card_icon_url,

  e.equipment_id,
  e.bonus_hp,
  e.bonus_atk,
  e.bonus_tech,
  e.init_bonus_hp,
  e.init_bonus_atk,
  e.init_bonus_tech,
  e.max_bonus_hp,
  e.max_bonus_atk,
  e.max_bonus_tech,
  e.buff_effect
FROM cards c
JOIN equipments e ON e.card_id = c.card_id
WHERE c.card_kind = 0
ORDER BY c.card_id ASC;