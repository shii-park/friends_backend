-- name: ListCollection :many
SELECT
  c.card_id,
  c.card_name,
  c.card_kind,
  c.rarity,
  c.card_icon_url,
  c.card_detail,

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
  eq.buff_effect     AS eq_buff_effect,

  COALESCE(uc.cnt, 0) AS owned_count,

  COALESCE(
    uc.latest_acquired_date,
    '0001-01-01 00:00:00+00'::timestamptz
  )::timestamptz AS latest_acquired_date

FROM cards c

LEFT JOIN (
  SELECT
    card_id,
    COUNT(*) AS cnt,
    MAX(acquired_date)::timestamptz AS latest_acquired_date
  FROM user_cards
  WHERE user_id = $1
  GROUP BY card_id
) uc ON uc.card_id = c.card_id

LEFT JOIN characters ch ON ch.card_id = c.card_id
LEFT JOIN equipments eq ON eq.card_id = c.card_id

ORDER BY c.card_kind, c.rarity DESC, c.card_name;