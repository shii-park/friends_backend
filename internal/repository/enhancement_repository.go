package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/domain"
	"github.com/shii-park/friends/internal/errs"
	"github.com/shii-park/friends/internal/sqlc"
)

type EnhancementRepository interface {
	EnhanceCardByCoin(ctx context.Context, userID uuid.UUID, instanceID uuid.UUID, times int) (newLevel int, cost int, err error)
}

type enhancementRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

func NewEnhancementRepository(db *sql.DB, q *sqlc.Queries) EnhancementRepository {
	return &enhancementRepository{db: db, q: q}
}

func (r *enhancementRepository) EnhanceCardByCoin(
	ctx context.Context,
	userID uuid.UUID,
	instanceID uuid.UUID,
	times int,
) (newLevel int, cost int, err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	qtx := r.q.WithTx(tx)

	// 所持カード
	uc, err := qtx.GetUserCardForUpdate(ctx, sqlc.GetUserCardForUpdateParams{
		InstanceID: instanceID,
		UserID:     userID,
	})
	if err != nil {
		return 0, 0, err
	}

	currentLevel := int(uc.Level.Int16)
	if currentLevel < domain.InitialLevel {
		currentLevel = domain.InitialLevel
	}

	// 上限で止める
	if currentLevel >= domain.MaxLevel {
		if err := tx.Commit(); err != nil {
			return 0, 0, err
		}
		return currentLevel, 0, nil
	}

	remain := domain.MaxLevel - currentLevel
	if times > remain {
		times = remain
	}
	if times <= 0 {
		if err := tx.Commit(); err != nil {
			return 0, 0, err
		}
		return currentLevel, 0, nil
	}

	// 合計コスト（あなたのテーブル）
	cost = domain.TotalCoinCostForLevelUps(currentLevel, times)
	if cost <= 0 {
		return 0, 0, errs.ErrInvalidEnhanceCost
	}

	// ユーザーコイン（ロック）
	coin, err := qtx.GetUserCoinForUpdate(ctx, userID)
	if err != nil {
		return 0, 0, err
	}
	if int(coin) < cost {
		return 0, 0, errs.ErrInsufficientCoin
	}

	newLevel = currentLevel + times
	if newLevel > domain.MaxLevel {
		newLevel = domain.MaxLevel
	}

	// コイン更新
	if err := qtx.UpdateUserCoin(ctx, sqlc.UpdateUserCoinParams{
		UserID: userID,
		Coin:   int32(int(coin) - cost),
	}); err != nil {
		return 0, 0, err
	}

	// レベル更新
	if err := qtx.UpdateUserCardLevel(ctx, sqlc.UpdateUserCardLevelParams{
		InstanceID: instanceID,
		UserID:     userID,
		Level: sql.NullInt16{
			Int16: int16(newLevel),
			Valid: true,
		},
	}); err != nil {
		return 0, 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, err
	}

	return newLevel, cost, nil
}
