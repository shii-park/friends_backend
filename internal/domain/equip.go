package domain

import (
	"github.com/shii-park/friends/internal/errs"
)

type Equip struct {
	*BaseCard `json:"card"`

	EquipID string `json:"equipID"`

	BonusHP   int `json:"bonusHP"`
	BonusATK  int `json:"bonusATK"`
	BonusTECH int `json:"bonusTECH"`

	InitBonusHP   int `json:"initBonusHP"`
	InitBonusATK  int `json:"initBonusATK"`
	InitBonusTECH int `json:"initBonusTECH"`

	MaxBonusHP   int `json:"maxBonusHP"`
	MaxBonusATK  int `json:"maxBonusATK"`
	MaxBonusTECH int `json:"maxBonusTECH"`

	BuffEffect BuffEffect `json:"-"`
}

func NewEquip(
	cardID string,
	equipID string,
	name string,
	icon string,
	rarity Rarity,

	initHP int,
	initATK int,
	initTECH int,

	maxHP int,
	maxATK int,
	maxTECH int,

	buff BuffEffect,
) (*Equip, error) {
	base, err := NewBaseCard(
		cardID,
		name,
		icon,
		rarity,
		CardKindEquip,
	)
	if err != nil {
		return nil, err
	}

	if equipID == "" {
		return nil, errs.ErrEquipIDRequired
	}

	if initHP < 0 {
		return nil, errs.ErrInvalidEquipHP
	}

	if initATK < 0 {
		return nil, errs.ErrInvalidEquipATK
	}

	if initTECH < 0 {
		return nil, errs.ErrInvalidEquipTECH
	}

	if maxHP < initHP {
		return nil, errs.ErrInvalidEquipMaxHP
	}

	if maxATK < initATK {
		return nil, errs.ErrInvalidEquipMaxATK
	}

	if maxTECH < initTECH {
		return nil, errs.ErrInvalidEquipMaxTECH
	}

	return &Equip{
		BaseCard: base,

		EquipID: equipID,

		BonusHP:   initHP,
		BonusATK:  initATK,
		BonusTECH: initTECH,

		InitBonusHP:   initHP,
		InitBonusATK:  initATK,
		InitBonusTECH: initTECH,

		MaxBonusHP:   maxHP,
		MaxBonusATK:  maxATK,
		MaxBonusTECH: maxTECH,

		BuffEffect: buff,
	}, nil
}
