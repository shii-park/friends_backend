-- name: GetCharacter :one
SELECT * FROM characters WHERE character_id=sqlc.arg(character_id) LiMIT 1;
-- name: GetAllCharacters :many
SELECT * FROM characters;
-- name: GetEquipment :one
SELECT * FROM equipments WHERE equipment_id=sqlc.arg(equipment_id) LIMIT 1;
-- name: GetAllEquipments :many
SELECT * FROM equipments;
