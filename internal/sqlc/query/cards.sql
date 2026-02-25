-- name: GetCharacter :one
SELECT * FROM characters WHERE character_id=sqlc.arg(character_id) LIMIT 1;

-- name: GetAllCharacters :many
SELECT * FROM characters;

-- name: GetEquipment :one
SELECT * FROM equipments WHERE equipment_id=sqlc.arg(equipment_id) LIMIT 1;

-- name: GetAllEquipments :many
SELECT * FROM equipments;

-- name: CreateCard :one
INSERT INTO cards (card_name, card_kind, card_icon_url, rarity, card_detail)
VALUES ($1, $2, $3, $4, $5)
RETURNING card_id;

-- name: CreateCharacter :exec
INSERT INTO characters (
    character_id, card_id, hp, atk, tech,
    init_hp, init_atk, init_tech,
    max_hp, max_atk, max_tech, special_type
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

-- name: CreateEquipment :exec
INSERT INTO equipments (
    equipment_id, card_id, bonus_hp, bonus_atk, bonus_tech,
    init_bonus_hp, init_bonus_atk, init_bonus_tech,
    max_bonus_hp, max_bonus_atk, max_bonus_tech, buff_effect
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

-- name: GetCharacterWithCard :one
SELECT ch.character_id, c.card_id, c.card_name, c.card_icon_url, c.rarity,
  ch.hp, ch.atk, ch.tech, ch.init_hp, ch.init_atk, ch.init_tech,
  ch.max_hp, ch.max_atk, ch.max_tech, ch.special_type
FROM characters ch JOIN cards c ON c.card_id = ch.card_id
WHERE ch.character_id = sqlc.arg(character_id) LIMIT 1;

-- name: GetAllCharactersWithCard :many
SELECT ch.character_id, c.card_id, c.card_name, c.card_icon_url, c.rarity,
  ch.hp, ch.atk, ch.tech, ch.init_hp, ch.init_atk, ch.init_tech,
  ch.max_hp, ch.max_atk, ch.max_tech, ch.special_type
FROM characters ch JOIN cards c ON c.card_id = ch.card_id;

-- name: GetEquipmentWithCard :one
SELECT e.equipment_id, c.card_id, c.card_name, c.card_icon_url, c.rarity,
  e.bonus_hp, e.bonus_atk, e.bonus_tech, e.init_bonus_hp, e.init_bonus_atk, e.init_bonus_tech,
  e.max_bonus_hp, e.max_bonus_atk, e.max_bonus_tech, e.buff_effect
FROM equipments e JOIN cards c ON c.card_id = e.card_id
WHERE e.equipment_id = sqlc.arg(equipment_id) LIMIT 1;

-- name: GetAllEquipmentsWithCard :many
SELECT e.equipment_id, c.card_id, c.card_name, c.card_icon_url, c.rarity,
  e.bonus_hp, e.bonus_atk, e.bonus_tech, e.init_bonus_hp, e.init_bonus_atk, e.init_bonus_tech,
  e.max_bonus_hp, e.max_bonus_atk, e.max_bonus_tech, e.buff_effect
FROM equipments e JOIN cards c ON c.card_id = e.card_id;
