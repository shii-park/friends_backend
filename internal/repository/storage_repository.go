package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/domain"
	"github.com/shii-park/friends/internal/errs"
	"github.com/shii-park/friends/internal/sqlc"
)

type StorageRepository struct {
	q *sqlc.Queries
}

func NewStorageRepository(q *sqlc.Queries) domain.StorageRepository {
	return &StorageRepository{q: q}
}

func (r *StorageRepository) AddCard(
	ctx context.Context,
	userID uuid.UUID,
	instanceID uuid.UUID,
	cardID int,
) error {
	return r.q.AddUserCard(ctx, sqlc.AddUserCardParams{
		InstanceID: instanceID,
		UserID:     userID,
		CardID:     int32(cardID),
	})
}

func (r *StorageRepository) RemoveCard(
	ctx context.Context,
	userID uuid.UUID,
	instanceID uuid.UUID,
) error {
	rows, err := r.q.RemoveUserCard(ctx, sqlc.RemoveUserCardParams{
		UserID:     userID,
		InstanceID: instanceID,
	})
	if err != nil {
		return err
	}
	if rows == 0 {
		return errs.ErrCardNotFound
	}
	return nil
}

func (r *StorageRepository) ListCards(
	ctx context.Context,
	userID uuid.UUID,
) ([]domain.CardInstance, error) {
	rows, err := r.q.ListUserCards(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := make([]domain.CardInstance, 0, len(rows))
	for _, row := range rows {
		res = append(res, domain.CardInstance{
			InstanceID: row.InstanceID,
			CardID:     int(row.CardID),
		})
	}

	return res, nil
}

func calcLinearStat(init, max, level int) int {
	// domain側に既にあるならそれを使ってOK（重複させない）
	if level <= 1 {
		return init
	}
	if level >= domain.MaxLevel {
		return max
	}
	denom := domain.MaxLevel - 1
	return init + (max-init)*(level-1)/denom
}

func (r *StorageRepository) ListCardDetails(
	ctx context.Context,
	userID uuid.UUID,
) ([]domain.CardInstanceDetail, error) {
	rows, err := r.q.ListUserCardDetails(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := make([]domain.CardInstanceDetail, 0, len(rows))
	for _, row := range rows {
		level := nullInt16ToInt(row.Level, 1)

		item := domain.CardInstanceDetail{
			InstanceID: row.InstanceID,
			Level:      level,
			Card: domain.CardMaster{
				CardID:      int(row.CardID),
				CardName:    row.CardName,
				CardKind:    row.CardKind,
				Rarity:      row.Rarity.String,
				CardIconURL: row.CardIconUrl.String,
			},
		}

		// character（現在値を計算）
		if row.CharacterID.Valid {
			initHP := int(row.ChInitHp.Int32)
			initATK := int(row.ChInitAtk.Int32)
			initTECH := int(row.ChInitTech.Int32)

			maxHP := int(row.ChMaxHp.Int32)
			maxATK := int(row.ChMaxAtk.Int32)
			maxTECH := int(row.ChMaxTech.Int32)

			hp := calcLinearStat(initHP, maxHP, level)
			atk := calcLinearStat(initATK, maxATK, level)
			tech := calcLinearStat(initTECH, maxTECH, level)

			item.Character = &domain.CharacterDetail{
				CharacterID: row.CharacterID.UUID,

				HP:   hp,
				ATK:  atk,
				TECH: tech,

				InitHP:   initHP,
				InitATK:  initATK,
				InitTECH: initTECH,
				MaxHP:    maxHP,
				MaxATK:   maxATK,
				MaxTECH:  maxTECH,

				SpecialType: row.ChSpecialType.String,
			}
		}

		// equipment（現在値を計算）
		if row.EquipmentID.Valid {
			bhp := calcLinearStat(int(row.EqInitBonusHp.Int32), int(row.EqMaxBonusHp.Int32), level)
			batk := calcLinearStat(int(row.EqInitBonusAtk.Int32), int(row.EqMaxBonusAtk.Int32), level)
			btech := calcLinearStat(int(row.EqInitBonusTech.Int32), int(row.EqMaxBonusTech.Int32), level)

			var buff *string
			if row.EqBuffEffect.Valid {
				v := row.EqBuffEffect.String
				buff = &v
			}

			item.Equipment = &domain.EquipmentDetail{
				EquipmentID: row.EquipmentID.UUID,

				BonusHP:   bhp,
				BonusATK:  batk,
				BonusTECH: btech,

				// ★ これを追加（init/max を返す）
				InitBonusHP:   int(row.EqInitBonusHp.Int32),
				InitBonusATK:  int(row.EqInitBonusAtk.Int32),
				InitBonusTECH: int(row.EqInitBonusTech.Int32),
				MaxBonusHP:    int(row.EqMaxBonusHp.Int32),
				MaxBonusATK:   int(row.EqMaxBonusAtk.Int32),
				MaxBonusTECH:  int(row.EqMaxBonusTech.Int32),

				BuffEffect: buff,
			}
		}

		res = append(res, item)
	}

	return res, nil
}

func nullInt16ToInt(v sql.NullInt16, def int) int {
	if v.Valid {
		return int(v.Int16)
	}
	return def
}
