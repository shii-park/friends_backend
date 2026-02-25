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
		// 石不足だと ErrNoRows になる想定
		return 0, nil, err
	}

	results := make([]DrawResult, 0, count)

	// 抽選→ストレージ付与
	for i := 0; i < count; i++ {
		// ここは service 側で決めて repository に渡してもOKだが、
		// 今回は「DBにカードを取りに行く」だけをrepoが担当する設計にする
		// → 実際の抽選は service がやる（下でやります）
		_ = i
	}

	// このrepoは service から「引いたカード情報（cardID/instanceID）」を受けて保存する形の方が綺麗なので、
	// 最終形は service で抽選して repo に保存を依頼する、にします（下で実装）。

	// ここでは何もしないので戻す
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
		// gacha_results がある場合だけ
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
