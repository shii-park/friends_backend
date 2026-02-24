package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/errs"
)

func NewCardID() string {
	return uuid.NewString()
}

const (
	InitialLevel = 0
	MaxLevel     = 10
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

type BaseCard struct {
	ID string `json:"cardID"`

	Name string `json:"cardName"`
	Icon string `json:"cardIcon,omitempty"` // アイコン、stringにしてるが実際どうなるかはわからない

	Rarity Rarity `json:"rarity"`

	AcquiredDate time.Time `json:"acquiredDate"`
	Level        int       `json:"level"`
}

func NewCard(id string, name string, icon string, rarity Rarity) (*BaseCard, error) {
	formattedName := strings.TrimSpace(name)
	if formattedName == "" {
		return nil, errs.ErrCardNameRequired
	}

	if !rarity.IsValid() {
		return nil, errs.ErrInvalidCardRarity
	}

	if strings.TrimSpace(id) == "" {
		id = NewCardID()
	}

	return &BaseCard{
		ID:           id,
		Name:         formattedName,
		Icon:         icon,
		Rarity:       rarity,
		AcquiredDate: time.Now(),
		Level:        InitialLevel,
	}, nil
}
