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
	db      *sql.DB
	q       *sqlc.Queries
	repo    repository.GachaRepository
	storage domain.StorageRepository

	rng         *rand.Rand
	saveHistory bool
}

func NewGachaService(
	db *sql.DB,
	q *sqlc.Queries,
	repo repository.GachaRepository,
	storage domain.StorageRepository,
) *GachaService {
	return &GachaService{
		db:          db,
		q:           q,
		repo:        repo,
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

func (s *GachaService) roll() (kind string, rarity string, isPickup bool) {
	if s.randPct() < 40.0 {
		kind = "character"
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
			rarity = "SSR"
		}
		if rarity == "SSR" {
			isPickup = s.randPct() < 60.0
		}
		return
	}

	kind = "equip"
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
		rarity = "SSR"
	}
	return
}

func (s *GachaService) randPct() float64 {
	return s.rng.Float64() * 100.0
}

func (s *GachaService) ListCharacters(ctx context.Context) ([]repository.GachaCharacter, error) {
	return s.repo.ListCharacters(ctx)
}

func (s *GachaService) ListEquipments(ctx context.Context) ([]repository.GachaEquipment, error) {
	return s.repo.ListEquipments(ctx)
}

type LineupResponse struct {
	Characters []repository.GachaCharacter `json:"characters"`
	Equipments []repository.GachaEquipment `json:"equipments"`
}

func (s *GachaService) Lineup(ctx context.Context) (LineupResponse, error) {
	chars, err := s.repo.ListCharacters(ctx)
	if err != nil {
		return LineupResponse{}, err
	}
	eqs, err := s.repo.ListEquipments(ctx)
	if err != nil {
		return LineupResponse{}, err
	}
	return LineupResponse{
		Characters: chars,
		Equipments: eqs,
	}, nil
}
