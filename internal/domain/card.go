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
	InitialLevel = 1
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
	Kind         CardKind  `json:"cardKind"`
}

func NewBaseCard(id string, name string, icon string, rarity Rarity, kind CardKind) (*BaseCard, error) {
	formattedName := strings.TrimSpace(name)
	if formattedName == "" {
		return nil, errs.ErrCardNameRequired
	}

	if !rarity.IsValid() {
		return nil, errs.ErrInvalidCardRarity
	}

	if strings.TrimSpace(id) == "" {
		return nil, errs.ErrCardIDRequired
	}

	if !kind.IsValid() {
		return nil, errs.ErrInvalidCardKind
	}

	return &BaseCard{
		ID:           id,
		Name:         formattedName,
		Icon:         icon,
		Rarity:       rarity,
		AcquiredDate: time.Now(),
		Level:        InitialLevel,
		Kind:         kind,
	}, nil
}

func (c *BaseCard) GetID() string              { return c.ID }
func (c *BaseCard) GetName() string            { return c.Name }
func (c *BaseCard) GetIcon() string            { return c.Icon }
func (c *BaseCard) GetRarity() Rarity          { return c.Rarity }
func (c *BaseCard) GetAcquiredDate() time.Time { return c.AcquiredDate }
func (c *BaseCard) GetLevel() int              { return c.Level }
func (c *BaseCard) GetKind() CardKind          { return c.Kind }
