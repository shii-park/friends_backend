package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/domain"
	"github.com/shii-park/friends/internal/sqlc"
)

var ErrNoRows = sql.ErrNoRows

type GachaRepository interface {
	Draw(ctx context.Context, userID uuid.UUID, count int) (newStone int, results []DrawResult, err error)
	ListCharacters(ctx context.Context) ([]GachaCharacter, error)
	ListEquipments(ctx context.Context) ([]GachaEquipment, error)
}

type DrawResult struct {
	InstanceID uuid.UUID
	CardID     int
	Kind       string // "character" or "equip"
	Rarity     string // "C" "UC" "R" "SR" "SSR" など
	IsPickup   bool
}

type gachaRepository struct {
	db          *sql.DB
	q           *sqlc.Queries
	storageRepo domain.StorageRepository
}

func NewGachaRepository(db *sql.DB, q *sqlc.Queries, storageRepo domain.StorageRepository) GachaRepository {
	return &gachaRepository{db: db, q: q, storageRepo: storageRepo}
}

func (r *gachaRepository) Draw(ctx context.Context, userID uuid.UUID, count int) (int, []DrawResult, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, nil, err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := r.q.WithTx(tx)

	// 1回=石1
	newStone, err := qtx.ConsumeUserGachaStone(ctx, sqlc.ConsumeUserGachaStoneParams{
		UserID: userID,
		Amount: int32(count),
	})
	if err != nil {
		return 0, nil, err
	}

	results := make([]DrawResult, 0, count)

	for i := 0; i < count; i++ {
		_ = i
	}

	if err := tx.Commit(); err != nil {
		return 0, nil, err
	}
	return int(newStone), results, nil
}

// 保存専用（serviceから呼ぶ用）
func (r *gachaRepository) SaveDraw(
	ctx context.Context,
	tx *sql.Tx,
	userID uuid.UUID,
	dr DrawResult,
	saveHistory bool,
) error {
	qtx := r.q.WithTx(tx)

	// storageに追加（既存）
	if err := r.storageRepo.AddCard(ctx, userID, dr.InstanceID, dr.CardID); err != nil {
		return err
	}

	if saveHistory {
		err := qtx.InsertGachaResult(ctx, sqlc.InsertGachaResultParams{
			ResultID:   uuid.New(),
			UserID:     userID,
			InstanceID: dr.InstanceID,
			CardID:     int32(dr.CardID),
			Kind:       dr.Kind,
			Rarity:     dr.Rarity,
			IsPickup:   dr.IsPickup,
			CreatedAt:  time.Now(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Begin/Commit を service に移すためのヘルパ
func (r *gachaRepository) BeginTx(ctx context.Context) (*sql.Tx, *sqlc.Queries, func(error) error, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, nil, err
	}
	qtx := r.q.WithTx(tx)

	finish := func(commitErr error) error {
		if commitErr != nil {
			_ = tx.Rollback()
			return commitErr
		}
		if err := tx.Commit(); err != nil {
			_ = tx.Rollback()
			return err
		}
		return nil
	}

	return tx, qtx, finish, nil
}

func isNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

type GachaCharacter struct {
	Card      GachaCard            `json:"card"`
	Character GachaCharacterDetail `json:"character"`
}

type GachaCard struct {
	CardID      int    `json:"cardID"`
	CardName    string `json:"cardName"`
	CardKind    int    `json:"cardKind"`
	Rarity      string `json:"rarity"`
	CardIconURL string `json:"cardIconURL"`
}

type GachaCharacterDetail struct {
	CharacterID string `json:"characterID"`
	HP          int    `json:"hp"`
	ATK         int    `json:"atk"`
	TECH        int    `json:"tech"`
	InitHP      int    `json:"initHP"`
	InitATK     int    `json:"initATK"`
	InitTECH    int    `json:"initTECH"`
	MaxHP       int    `json:"maxHP"`
	MaxATK      int    `json:"maxATK"`
	MaxTECH     int    `json:"maxTECH"`
	SpecialType string `json:"specialType"`
}

func (r *gachaRepository) ListCharacters(ctx context.Context) ([]GachaCharacter, error) {
	rows, err := r.q.ListGachaCharacters(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]GachaCharacter, 0, len(rows))
	for _, row := range rows {
		rarity := ""
		if row.Rarity.Valid {
			rarity = row.Rarity.String
		}
		iconURL := ""
		if row.CardIconUrl.Valid {
			iconURL = row.CardIconUrl.String
		}
		specialType := ""
		if row.SpecialType.Valid {
			specialType = row.SpecialType.String
		}

		out = append(out, GachaCharacter{
			Card: GachaCard{
				CardID:      int(row.CardID),
				CardName:    row.CardName,
				CardKind:    int(row.CardKind),
				Rarity:      rarity,
				CardIconURL: iconURL,
			},
			Character: GachaCharacterDetail{
				CharacterID: row.CharacterID.String(),
				HP:          int(row.Hp),
				ATK:         int(row.Atk),
				TECH:        int(row.Tech),
				InitHP:      int(row.InitHp),
				InitATK:     int(row.InitAtk),
				InitTECH:    int(row.InitTech),
				MaxHP:       int(row.MaxHp),
				MaxATK:      int(row.MaxAtk),
				MaxTECH:     int(row.MaxTech),
				SpecialType: specialType,
			},
		})
	}
	return out, nil
}

type GachaEquipment struct {
	Card      GachaCard            `json:"card"`
	Equipment GachaEquipmentDetail `json:"equipment"`
}

type GachaEquipmentDetail struct {
	EquipmentID string `json:"equipmentID"`

	BonusHP   int `json:"bonusHP"`
	BonusATK  int `json:"bonusATK"`
	BonusTECH int `json:"bonusTECH"`

	InitBonusHP   int `json:"initBonusHP"`
	InitBonusATK  int `json:"initBonusATK"`
	InitBonusTECH int `json:"initBonusTECH"`

	MaxBonusHP   int `json:"maxBonusHP"`
	MaxBonusATK  int `json:"maxBonusATK"`
	MaxBonusTECH int `json:"maxBonusTECH"`

	BuffEffect string `json:"buffEffect,omitempty"`
}

func (r *gachaRepository) ListEquipments(ctx context.Context) ([]GachaEquipment, error) {
	rows, err := r.q.ListGachaEquipments(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]GachaEquipment, 0, len(rows))
	for _, row := range rows {
		// cards 側の nullable 対応（あなたの rarity が NullString だったので同様に）
		rarity := ""
		if row.Rarity.Valid {
			rarity = row.Rarity.String
		}
		iconURL := ""
		if row.CardIconUrl.Valid {
			iconURL = row.CardIconUrl.String
		}

		buff := ""
		if row.BuffEffect.Valid {
			buff = row.BuffEffect.String
		}

		out = append(out, GachaEquipment{
			Card: GachaCard{
				CardID:      int(row.CardID),
				CardName:    row.CardName,
				CardKind:    int(row.CardKind),
				Rarity:      rarity,
				CardIconURL: iconURL,
			},
			Equipment: GachaEquipmentDetail{
				EquipmentID: row.EquipmentID.String(),

				BonusHP:   int(row.BonusHp),
				BonusATK:  int(row.BonusAtk),
				BonusTECH: int(row.BonusTech),

				InitBonusHP:   int(row.InitBonusHp),
				InitBonusATK:  int(row.InitBonusAtk),
				InitBonusTECH: int(row.InitBonusTech),

				MaxBonusHP:   int(row.MaxBonusHp),
				MaxBonusATK:  int(row.MaxBonusAtk),
				MaxBonusTECH: int(row.MaxBonusTech),

				BuffEffect: buff,
			},
		})
	}
	return out, nil
}
