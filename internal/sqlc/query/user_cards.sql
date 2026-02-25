-- name: AddUserCard :exec
INSERT INTO user_cards (instance_id, user_id, card_id)
VALUES ($1, $2, $3);

-- name: RemoveUserCard :execrows
DELETE FROM user_cards
WHERE user_id = $1 AND instance_id = $2;

-- name: ListUserCards :many
SELECT instance_id, card_id
FROM user_cards
WHERE user_id = $1
ORDER BY acquired_date DESC;

ALTER TABLE user_cards
ADD COLUMN level INT NOT NULL DEFAULT 1;

-- name: IncreaseUserCardLevel :execrows
UPDATE user_cards
SET level = level + 1
WHERE user_id = $1 AND instance_id = $2 AND level < $3;