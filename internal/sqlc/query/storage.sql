-- name: ListUserCardDetails :many
SELECT
  uc.instance_id,

  c.card_id         AS card_id,
  c.card_name       AS card_name,
  c.card_kind       AS card_kind,
  c.rarity          AS rarity,
  c.card_icon_url   AS card_icon_url,

  ch.character_id   AS character_id,
  ch.hp             AS ch_hp,
  ch.atk            AS ch_atk,
  ch.tech           AS ch_tech,
  ch.init_hp        AS ch_init_hp,
  ch.init_atk       AS ch_init_atk,
  ch.init_tech      AS ch_init_tech,
  ch.max_hp         AS ch_max_hp,
  ch.max_atk        AS ch_max_atk,
  ch.max_tech       AS ch_max_tech,
  ch.special_type   AS ch_special_type,

  e.equipment_id    AS equipment_id,
  e.bonus_hp        AS eq_bonus_hp,
  e.bonus_atk       AS eq_bonus_atk,
  e.bonus_tech      AS eq_bonus_tech,
  e.init_bonus_hp   AS eq_init_bonus_hp,
  e.init_bonus_atk  AS eq_init_bonus_atk,
  e.init_bonus_tech AS eq_init_bonus_tech,
  e.max_bonus_hp    AS eq_max_bonus_hp,
  e.max_bonus_atk   AS eq_max_bonus_atk,
  e.max_bonus_tech  AS eq_max_bonus_tech,
  e.buff_effect     AS eq_buff_effect
FROM user_cards uc
JOIN cards c ON c.card_id = uc.card_id
LEFT JOIN characters ch ON ch.card_id = c.card_id
LEFT JOIN equipments e ON e.card_id = c.card_id
WHERE uc.user_id = sqlc.arg(user_id)
ORDER BY uc.instance_id DESC;