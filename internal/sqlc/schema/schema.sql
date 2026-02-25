CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    user_name VARCHAR(255) NOT NULL,
    icon_url TEXT,
    profile_message TEXT,
    birth_month SMALLINT CHECK (birth_month >= 1 AND birth_month <= 12),
    birth_day SMALLINT CHECK (birth_day >= 1 AND birth_day <= 31),
    latest_login_date TIMESTAMPTZ,
    streak_login_days INTEGER DEFAULT 0 NOT NULL,
    registered_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    coin INTEGER DEFAULT 0 NOT NULL,
    gacha_stone INTEGER DEFAULT 0 NOT NULL,
    rank_point INTEGER DEFAULT 0 NOT NULL
);
--
CREATE TABLE cards (
    card_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    card_name VARCHAR(255) NOT NULL,
    card_kind SMALLINT NOT NULL,
    card_icon_url TEXT,
    rarity VARCHAR(3) CHECK (rarity IN ('C', 'UC', 'R', 'SR', 'SSR'))
);

CREATE TABLE characters (
    character_id UUID PRIMARY KEY,
    card_id INTEGER UNIQUE NOT NULL REFERENCES cards(card_id) ON DELETE CASCADE,
    hp INTEGER NOT NULL,
    atk INTEGER NOT NULL,
    tech INTEGER NOT NULL,
    init_hp INTEGER NOT NULL,
    init_atk INTEGER NOT NULL,
    init_tech INTEGER NOT NULL,
    max_hp INTEGER NOT NULL,
    max_atk INTEGER NOT NULL,
    max_tech INTEGER NOT NULL,
    special_type VARCHAR(10) CHECK (special_type IN ('rock', 'paper', 'scissors'))
);

CREATE TABLE equipments (
    equipment_id UUID PRIMARY KEY,
    card_id INTEGER UNIQUE NOT NULL REFERENCES cards(card_id) ON DELETE CASCADE,
    bonus_hp INTEGER NOT NULL,
    bonus_atk INTEGER NOT NULL,
    bonus_tech INTEGER NOT NULL,
    init_bonus_hp INTEGER NOT NULL,
    init_bonus_atk INTEGER NOT NULL,
    init_bonus_tech INTEGER NOT NULL,
    max_bonus_hp INTEGER NOT NULL,
    max_bonus_atk INTEGER NOT NULL,
    max_bonus_tech INTEGER NOT NULL,
    buff_effect TEXT
);

CREATE TABLE user_cards (
    instance_id UUID PRIMARY KEY, -- 1枚ごとの固有ID（ダブりごとに別々のUUIDが発行される）
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    card_id INTEGER NOT NULL REFERENCES cards(card_id),
    level SMALLINT DEFAULT 1 CHECK (level >= 1 AND level <= 10),
    acquired_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
)