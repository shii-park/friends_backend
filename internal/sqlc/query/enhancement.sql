-- name: GetUserCoinForUpdate :one
SELECT coin
FROM users
WHERE user_id = $1
FOR UPDATE;

-- name: GetUserCardForUpdate :one
SELECT instance_id, user_id, card_id, level
FROM user_cards
WHERE instance_id = $1 AND user_id = $2
FOR UPDATE;

-- name: UpdateUserCardLevel :exec
UPDATE user_cards
SET level = $3
WHERE instance_id = $1 AND user_id = $2;

-- name: GetCardKindByInstanceID :one
SELECT c.card_id, c.card_kind, c.rarity
FROM user_cards uc
JOIN cards c ON c.card_id = uc.card_id
WHERE uc.instance_id = $1 AND uc.user_id = $2;

-- name: GetCharacterByCardID :one
SELECT *
FROM characters
WHERE card_id = $1;

-- name: GetEquipmentByCardID :one
SELECT *
FROM equipments
WHERE card_id = $1;