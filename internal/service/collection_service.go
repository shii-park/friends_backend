package service

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/domain"
	"github.com/shii-park/friends/internal/errs"
	"github.com/shii-park/friends/internal/repository"
	"github.com/shii-park/friends/internal/sqlc"
)

type CollectionService interface {
	List(ctx context.Context, userID uuid.UUID) ([]CollectionEntryDTO, error)
}

type CollectionEntryDTO struct {
	Card  CardDTO          `json:"card"`
	State domain.CardState `json:"state"`
	Count int              `json:"count"`
}

type CardDTO struct {
	Base BaseCardDTO `json:"base"`

	Character *CharacterDTO `json:"character,omitempty"`
	Equip     *EquipDTO     `json:"equip,omitempty"`
}

type BaseCardDTO struct {
	ID     string `json:"cardID"` // domainに合わせて文字列で返す（DBはintでもOK）
	Name   string `json:"cardName"`
	Icon   string `json:"cardIcon,omitempty"`
	Detail string `json:"detail,omitempty"`

	Rarity domain.Rarity   `json:"rarity"`
	Kind   domain.CardKind `json:"cardKind"`

	LatestAcquiredDate time.Time `json:"latestAcquiredDate"`
}

type CharacterDTO struct {
	CharaID     string             `json:"charaID"`
	InitHP      int                `json:"initHP"`
	InitATK     int                `json:"initATK"`
	InitTECH    int                `json:"initTECH"`
	MaxHP       int                `json:"maxHP"`
	MaxATK      int                `json:"maxATK"`
	MaxTECH     int                `json:"maxTECH"`
	SpecialType domain.SpecialType `json:"specialType"`
}

type EquipDTO struct {
	EquipID       string `json:"equipID"`
	InitBonusHP   int    `json:"initBonusHP"`
	InitBonusATK  int    `json:"initBonusATK"`
	InitBonusTECH int    `json:"initBonusTECH"`
	MaxBonusHP    int    `json:"maxBonusHP"`
	MaxBonusATK   int    `json:"maxBonusATK"`
	MaxBonusTECH  int    `json:"maxBonusTECH"`
	BuffEffect    string `json:"buffEffect,omitempty"`
}

type collectionService struct {
	repo repository.CollectionRepository
}

func NewCollectionService(repo repository.CollectionRepository) CollectionService {
	return &collectionService{repo: repo}
}

func (s *collectionService) List(ctx context.Context, userID uuid.UUID) ([]CollectionEntryDTO, error) {
	rows, err := s.repo.ListCollection(ctx, userID)
	if err != nil {
		return nil, err
	}

	out := make([]CollectionEntryDTO, 0, len(rows))
	for _, r := range rows {
		entry, err := toCollectionEntryDTO(r)
		if err != nil {
			return nil, err
		}
		out = append(out, entry)
	}
	return out, nil
}

func toCollectionEntryDTO(r sqlc.ListCollectionRow) (CollectionEntryDTO, error) {
	kind := domain.CardKind(int(r.CardKind))
	if !kind.IsValid() {
		return CollectionEntryDTO{}, errs.ErrInvalidCardKind
	}

	base := BaseCardDTO{
		ID:   strconv.Itoa(int(r.CardID)),
		Name: r.CardName,
		Kind: kind,
	}

	// icon
	base.Icon = nullString(r.CardIconUrl)
	base.Detail = nullString(r.CardDetail)

	// rarity: sql.NullString -> domain.Rarity
	if r.Rarity.Valid {
		base.Rarity = domain.Rarity(r.Rarity.String)
	} else {
		base.Rarity = domain.RarityUnknown
	}

	base.LatestAcquiredDate = r.LatestAcquiredDate
	hasHistory := !r.LatestAcquiredDate.IsZero()

	count := int(r.OwnedCount)

	// state
	state := domain.NotFound
	if count > 0 {
		state = domain.Get
	} else if hasHistory {
		state = domain.Find
	}

	card := CardDTO{Base: base}

	switch kind {
	case domain.CardKindCharacter:
		if !r.CharacterID.Valid {
			return CollectionEntryDTO{}, errs.ErrInvalidCardKind
		}
		card.Character = &CharacterDTO{
			CharaID:     r.CharacterID.UUID.String(),
			InitHP:      nullInt32(r.ChInitHp),
			InitATK:     nullInt32(r.ChInitAtk),
			InitTECH:    nullInt32(r.ChInitTech),
			MaxHP:       nullInt32(r.ChMaxHp),
			MaxATK:      nullInt32(r.ChMaxAtk),
			MaxTECH:     nullInt32(r.ChMaxTech),
			SpecialType: nullSpecialType(r.ChSpecialType),
		}

	case domain.CardKindEquip:
		if !r.EquipmentID.Valid {
			return CollectionEntryDTO{}, errs.ErrInvalidCardKind
		}
		eq := &EquipDTO{
			EquipID:       r.EquipmentID.UUID.String(),
			InitBonusHP:   nullInt32(r.EqInitBonusHp),
			InitBonusATK:  nullInt32(r.EqInitBonusAtk),
			InitBonusTECH: nullInt32(r.EqInitBonusTech),
			MaxBonusHP:    nullInt32(r.EqMaxBonusHp),
			MaxBonusATK:   nullInt32(r.EqMaxBonusAtk),
			MaxBonusTECH:  nullInt32(r.EqMaxBonusTech),
		}
		eq.BuffEffect = nullString(r.EqBuffEffect)
		card.Equip = eq

	default:
		return CollectionEntryDTO{}, errs.ErrInvalidCardKind
	}

	return CollectionEntryDTO{
		Card:  card,
		State: state,
		Count: count,
	}, nil
}

func nullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func nullInt32(n sql.NullInt32) int {
	if n.Valid {
		return int(n.Int32)
	}
	return 0
}

func nullSpecialType(ns sql.NullString) domain.SpecialType {
	if ns.Valid {
		return domain.SpecialType(ns.String)
	}
	return domain.SpecialTypeUnknown
}
