package repository

import (
	"context"

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
		item := domain.CardInstanceDetail{
			InstanceID: row.InstanceID,
			Card: domain.CardMaster{
				CardID:      int(row.CardID),
				CardName:    row.CardName,
				CardKind:    row.CardKind,
				Rarity:      row.Rarity.String, // rarity が NULL許容なら NullString になる
				CardIconURL: row.CardIconUrl.String,
			},
		}

		// character（LEFT JOINなのでNULLあり）
		if row.CharacterID.Valid {
			item.Character = &domain.CharacterDetail{
				CharacterID: row.CharacterID.UUID,
				HP:          int(row.ChHp.Int32),
				ATK:         int(row.ChAtk.Int32),
				TECH:        int(row.ChTech.Int32),
				InitHP:      int(row.ChInitHp.Int32),
				InitATK:     int(row.ChInitAtk.Int32),
				InitTECH:    int(row.ChInitTech.Int32),
				MaxHP:       int(row.ChMaxHp.Int32),
				MaxATK:      int(row.ChMaxAtk.Int32),
				MaxTECH:     int(row.ChMaxTech.Int32),
				SpecialType: row.ChSpecialType.String,
			}
		}

		// equipment
		if row.EquipmentID.Valid {
			var buff *string
			if row.EqBuffEffect.Valid {
				v := row.EqBuffEffect.String
				buff = &v
			}
			item.Equipment = &domain.EquipmentDetail{
				EquipmentID:   row.EquipmentID.UUID,
				BonusHP:       int(row.EqBonusHp.Int32),
				BonusATK:      int(row.EqBonusAtk.Int32),
				BonusTECH:     int(row.EqBonusTech.Int32),
				InitBonusHP:   int(row.EqInitBonusHp.Int32),
				InitBonusATK:  int(row.EqInitBonusAtk.Int32),
				InitBonusTECH: int(row.EqInitBonusTech.Int32),
				MaxBonusHP:    int(row.EqMaxBonusHp.Int32),
				MaxBonusATK:   int(row.EqMaxBonusAtk.Int32),
				MaxBonusTECH:  int(row.EqMaxBonusTech.Int32),
				BuffEffect:    buff,
			}
		}

		res = append(res, item)
	}

	return res, nil
}
