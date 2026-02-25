package domain

import (
	"github.com/shii-park/friends/internal/errs"
)

// SpecialType 一旦ここで実装
// そのうちAtackTypeの定義をしてくれるのでそっちを用いる
type SpecialType string

const (
	SpecialTypeUnknown  SpecialType = ""
	SpecialTypeRock     SpecialType = "rock"
	SpecialTypeScissors SpecialType = "scissors"
	SpecialTypePaper    SpecialType = "paper"
)

func (s SpecialType) IsValid() bool {
	switch s {
	case SpecialTypeRock, SpecialTypeScissors, SpecialTypePaper:
		return true
	default:
		return false
	}
}

type Character struct {
	*BaseCard `json:"card"`

	CharaID string `json:"charaID"`

	HP   int `json:"hp"`
	ATK  int `json:"atk"`
	TECH int `json:"tech"`

	InitHP   int `json:"initHP"`
	InitATK  int `json:"initATK"`
	InitTECH int `json:"initTECH"`

	MaxHP   int `json:"maxHP"`
	MaxATK  int `json:"maxATK"`
	MaxTECH int `json:"maxTECH"`

	SpecialType SpecialType `json:"specialType"`
}

func NewCharacter(
	cardID string,
	charaID string,
	name string,
	icon string,
	rarity Rarity,

	initHP, initATK, initTECH int,
	maxHP, maxATK, maxTECH int,

	specialType SpecialType,
) (*Character, error) {
	base, err := NewBaseCard(
		cardID,
		name,
		icon,
		rarity,
		CardKindCharacter,
	)
	if err != nil {
		return nil, err
	}

	if charaID == "" {
		return nil, errs.ErrCharaIDRequired
	}

	if initHP <= 0 {
		return nil, errs.ErrInvalidCharacterHP
	}

	if initATK < 0 {
		return nil, errs.ErrInvalidCharacterATK
	}

	if initTECH < 0 {
		return nil, errs.ErrInvalidCharacterTECH
	}

	if maxHP < initHP {
		return nil, errs.ErrInvalidCharacterMaxHP
	}

	if maxATK < initATK {
		return nil, errs.ErrInvalidCharacterMaxATK
	}

	if maxTECH < initTECH {
		return nil, errs.ErrInvalidCharacterMaxTECH
	}

	if !specialType.IsValid() {
		return nil, errs.ErrInvalidSpecialType
	}

	return &Character{
		BaseCard: base,

		CharaID: charaID,

		HP:   initHP,
		ATK:  initATK,
		TECH: initTECH,

		InitHP:   initHP,
		InitATK:  initATK,
		InitTECH: initTECH,

		MaxHP:   maxHP,
		MaxATK:  maxATK,
		MaxTECH: maxTECH,

		SpecialType: specialType,
	}, nil
}

func (c *Character) AddHP(delta int) error {
	if delta <= 0 {
		return errs.ErrInvalidCharaHPDelta
	}
	c.HP += delta
	if c.HP > c.MaxHP {
		c.HP = c.MaxHP
	}
	return nil
}

func (c *Character) AddATK(delta int) error {
	if delta <= 0 {
		return errs.ErrInvalidCharaATKDelta
	}
	c.ATK += delta
	if c.ATK > c.MaxATK {
		c.ATK = c.MaxATK
	}
	return nil
}

func (c *Character) AddTECH(delta int) error {
	if delta <= 0 {
		return errs.ErrInvalidCharaTECHDelta
	}
	c.TECH += delta
	if c.TECH > c.MaxTECH {
		c.TECH = c.MaxTECH
	}
	return nil
}

func (c *Character) LevelUp() error {
	if c.Level >= MaxLevel {
		return errs.ErrCharaAlreadyMaxLevel
	}
	return c.SetLevel(c.Level + 1)
}

func (c *Character) SetLevel(level int) error {
	if level < 1 || level > MaxLevel {
		return errs.ErrInvalidCharaLevel
	}

	c.Level = level

	c.HP = calcLinearStat(c.InitHP, c.MaxHP, level)
	c.ATK = calcLinearStat(c.InitATK, c.MaxATK, level)
	c.TECH = calcLinearStat(c.InitTECH, c.MaxTECH, level)

	return nil
}

func calcLinearStat(init int, max int, level int) int {
	// Lv1はinit、LvMaxは必ずmaxにする
	if level <= 1 {
		return init
	}
	if level >= MaxLevel {
		return max
	}

	denom := MaxLevel - 1
	step := (max - init) / denom
	return init + step*(level-1)
}
