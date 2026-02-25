package service

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/domain"
	"github.com/shii-park/friends/internal/errs"
	"github.com/shii-park/friends/internal/repository"
	"github.com/shii-park/friends/internal/sqlc"
)

type GachaService struct {
	db        *sql.DB
	q         *sqlc.Queries
	gachaRepo *repository.GachaRepository
	storage   domain.StorageRepository

	rng         *rand.Rand
	saveHistory bool
}

func NewGachaService(db *sql.DB, q *sqlc.Queries, storage domain.StorageRepository) *GachaService {
	return &GachaService{
		db:          db,
		q:           q,
		storage:     storage,
		rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
		saveHistory: false,
	}
}

type DrawResult struct {
	InstanceID uuid.UUID `json:"instanceID"`
	CardID     int       `json:"cardID"`
	Kind       string    `json:"kind"`   // character / equip
	Rarity     string    `json:"rarity"` // C UC R SR SSR
	IsPickup   bool      `json:"isPickup"`
}

func (s *GachaService) Draw(ctx context.Context, userID uuid.UUID, count int) (newStone int, results []DrawResult, err error) {
	if count != 1 && count != 10 {
		return 0, nil, errs.ErrInvalidGachaCount
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, nil, err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := s.q.WithTx(tx)

	// 石消費（1回=1）
	ns, err := qtx.ConsumeUserGachaStone(ctx, sqlc.ConsumeUserGachaStoneParams{
		UserID: userID,
		Amount: int32(count),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil, errs.ErrInsufficientGachaStone
		}
		return 0, nil, err
	}

	out := make([]DrawResult, 0, count)

	for i := 0; i < count; i++ {
		kind, rarity, isPickup := s.roll()

		rarityNS := sql.NullString{String: rarity, Valid: true}

		var cardID int32

		switch kind {
		case "character":
			cardID, err = qtx.GetRandomCharacterByRarity(ctx, rarityNS)
		case "equip":
			cardID, err = qtx.GetRandomEquipByRarity(ctx, rarityNS)
		default:
			return 0, nil, errs.ErrInvalidGachaKind
		}
		if err != nil {
			return 0, nil, err
		}

		instanceID := uuid.New()

		// 既存 storage に付与
		if err := s.storage.AddCard(ctx, userID, instanceID, int(cardID)); err != nil {
			return 0, nil, err
		}

		out = append(out, DrawResult{
			InstanceID: instanceID,
			CardID:     int(cardID),
			Kind:       kind,
			Rarity:     rarity,
			IsPickup:   isPickup,
		})
	}

	if err := tx.Commit(); err != nil {
		return 0, nil, err
	}

	return int(ns), out, nil
}

// roll は仕様どおりの確率で kind/rarity/pickup を決める
func (s *GachaService) roll() (kind string, rarity string, isPickup bool) {
	// kind
	if s.randPct() < 40.0 {
		kind = "character"
		// character rarity（合計40%）
		r := s.randPct()
		switch {
		case r < 21.25:
			rarity = "C"
		case r < 21.25+10.0:
			rarity = "UC"
		case r < 21.25+10.0+5.0:
			rarity = "R"
		case r < 21.25+10.0+5.0+2.5:
			rarity = "SR"
		default:
			rarity = "SSR" // 1.25
		}

		// pickup（SSRキャラのみ：恒常0.5 / PU0.75 → SSRキャラ内で PU=0.75/1.25=60%）
		if rarity == "SSR" {
			isPickup = s.randPct() < 60.0
		}
		return
	}

	kind = "equip"
	// equip rarity（合計60%）
	r := s.randPct()
	switch {
	case r < 31.875:
		rarity = "C"
	case r < 31.875+15.0:
		rarity = "UC"
	case r < 31.875+15.0+7.5:
		rarity = "R"
	case r < 31.875+15.0+7.5+3.75:
		rarity = "SR"
	default:
		rarity = "SSR" // 1.875
	}
	return
}

func (s *GachaService) randPct() float64 {
	return s.rng.Float64() * 100.0
}
