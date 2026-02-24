package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/errs"
)

type CardID string

func NewCardID() string {
	return uuid.NewString()
}

const (
	InitialLevel = 0
)

type Rarity string

const (
	RarityUnknown        Rarity = ""
	RarityCommon         Rarity = "C"
	RarityUnCommon       Rarity = "UC"
	RarityRare           Rarity = "R"
	RaritySuperRare      Rarity = "SR"
	RaritySuperSuperRare Rarity = "SSR"
)

func (r Rarity) IsValid() bool {
	switch r {
	case RarityCommon, RarityUnCommon, RarityRare, RaritySuperRare, RaritySuperSuperRare:
		return true
	default:
		return false
	}
}

type Card struct {
	ID string `json:"cardID"`

	Name string `json:"cardName"`
	Icon string `json:"cardIcon,omitempty"` // アイコン、stringにしてるが実際どうなるかはわからない

	Rarity Rarity `json:"rarity"`

	AcquiredDate time.Time `json:"acquiredDate"`
	Level        int       `json:"level"`
}

func NewCard(id string, name string, icon string, rarity Rarity) (*Card, error) {
	formattedName := strings.TrimSpace(name)
	if formattedName == "" {
		return nil, errs.ErrCardNameRequired
	}

	if !rarity.IsValid() {
		return nil, errs.ErrInvalidCardRarity
	}

	if strings.TrimSpace(id) == "" {
		id = string(NewCardID()) // CardID型ならstring変換
	}

	return &Card{
		ID:           id,
		Name:         formattedName,
		Icon:         icon,
		Rarity:       rarity,
		AcquiredDate: time.Now(),
		Level:        InitialLevel,
	}, nil
}
