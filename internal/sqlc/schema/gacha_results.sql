CREATE TABLE gacha_results (
    result_id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,

    instance_id UUID NOT NULL,
    card_id INT NOT NULL,

    kind TEXT NOT NULL,
    rarity TEXT NOT NULL,
    is_pickup BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);