-- User系

-- name: GetUser :one
SELECT * FROM users
WHERE user_id = sqlc.arg(user_id) LIMIT 1;


-- name: CreateUser :exec
INSERT INTO users (
    user_id, user_name, icon_url, profile_message, birth_month, birth_day, registered_date, latest_login_date, streak_login_days, rank_point, coin, gacha_stone
) VALUES (
    sqlc.arg(user_id), sqlc.arg(user_name), sqlc.arg(icon_url), sqlc.arg(profile_message), sqlc.arg(birth_month), sqlc.arg(birth_day), sqlc.arg(registered_date), sqlc.arg(latest_login_date), sqlc.arg(streak_login_days),  sqlc.arg(rank_point), sqlc.arg(coin), sqlc.arg(gacha_stone)
);

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = sqlc.arg(user_id);

-- name: UpdateUserCoin :exec
UPDATE users
SET coin = sqlc.arg(coin)
WHERE user_id = sqlc.arg(user_id);

-- name: UpdateUserGachaStone :exec
UPDATE users
SET gacha_stone = sqlc.arg(gacha_stone)
WHERE user_id = sqlc.arg(user_id);

-- name: UpdateUserRankPoint :exec
UPDATE users
SET rank_point = sqlc.arg(rank_point)
WHERE user_id = sqlc.arg(user_id);

