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

func (e *Equip) AddBonusHP(delta int) error {
	if delta <= 0 {
		return errs.ErrInvalidEquipHPDelta
	}
	e.BonusHP += delta
	if e.BonusHP > e.MaxBonusHP {
		e.BonusHP = e.MaxBonusHP
	}
	return nil
}

func (e *Equip) AddBonusATK(delta int) error {
	if delta <= 0 {
		return errs.ErrInvalidEquipATKDelta
	}
	e.BonusATK += delta
	if e.BonusATK > e.MaxBonusATK {
		e.BonusATK = e.MaxBonusATK
	}
	return nil
}

func (e *Equip) AddBonusTECH(delta int) error {
	if delta <= 0 {
		return errs.ErrInvalidEquipTECHDelta
	}
	e.BonusTECH += delta
	if e.BonusTECH > e.MaxBonusTECH {
		e.BonusTECH = e.MaxBonusTECH
	}
	return nil
}

func (e *Equip) LevelUp() error {
	if e.Level >= MaxLevel {
		return errs.ErrEquipAlreadyMaxLevel
	}
	return e.SetLevel(e.Level + 1)
}

func (e *Equip) SetLevel(level int) error {
	if level < 1 || level > MaxLevel {
		return errs.ErrInvalidEquipLevel
	}

	e.Level = level

	e.BonusHP = calcLinearBonus(e.InitBonusHP, e.MaxBonusHP, level)
	e.BonusATK = calcLinearBonus(e.InitBonusATK, e.MaxBonusATK, level)
	e.BonusTECH = calcLinearBonus(e.InitBonusTECH, e.MaxBonusTECH, level)

	return nil
}

func calcLinearBonus(init int, max int, level int) int {
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
