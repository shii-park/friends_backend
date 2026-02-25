package domain

import (
	"github.com/shii-park/friends/internal/errs"
)

// SpecialType 一旦ここで実装
// そのうちAtackTypeの定義をしてくれるのでそっちを用いる
type SpecialType string

const (
	SpecialTypeUnknown SpecialType = ""
	SpecialTypeGu      SpecialType = "gu"
	SpecialTypeChoki   SpecialType = "choki"
	SpecialTypePa      SpecialType = "pa"
)

func (s SpecialType) IsValid() bool {
	switch s {
	case SpecialTypeGu, SpecialTypeChoki, SpecialTypePa:
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
