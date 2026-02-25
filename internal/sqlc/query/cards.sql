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
