-- name: ListUserCardDetails :many
SELECT
  uc.instance_id,
  uc.card_id,
  uc.level,

  c.card_name,
  c.card_kind,
  c.rarity,
  c.card_icon_url,

  ch.character_id,
  ch.init_hp        AS ch_init_hp,
  ch.init_atk       AS ch_init_atk,
  ch.init_tech      AS ch_init_tech,
  ch.max_hp         AS ch_max_hp,
  ch.max_atk        AS ch_max_atk,
  ch.max_tech       AS ch_max_tech,
  ch.special_type   AS ch_special_type,

  eq.equipment_id,
  eq.init_bonus_hp   AS eq_init_bonus_hp,
  eq.init_bonus_atk  AS eq_init_bonus_atk,
  eq.init_bonus_tech AS eq_init_bonus_tech,
  eq.max_bonus_hp    AS eq_max_bonus_hp,
  eq.max_bonus_atk   AS eq_max_bonus_atk,
  eq.max_bonus_tech  AS eq_max_bonus_tech,
  eq.buff_effect     AS eq_buff_effect

FROM user_cards uc
JOIN cards c ON c.card_id = uc.card_id
LEFT JOIN characters ch ON ch.card_id = c.card_id
LEFT JOIN equipments eq ON eq.card_id = c.card_id
WHERE uc.user_id = $1
ORDER BY uc.acquired_date DESC;