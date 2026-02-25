package domain

import "github.com/google/uuid"

type CardMaster struct {
	CardID      int    `json:"cardID"`
	CardName    string `json:"cardName"`
	CardKind    int16  `json:"cardKind"` // 0 equip / 1 character
	Rarity      string `json:"rarity"`
	CardIconURL string `json:"cardIconURL"`
}

type CharacterDetail struct {
	CharacterID uuid.UUID `json:"characterID"`
	HP          int       `json:"hp"`
	ATK         int       `json:"atk"`
	TECH        int       `json:"tech"`
	InitHP      int       `json:"initHP"`
	InitATK     int       `json:"initATK"`
	InitTECH    int       `json:"initTECH"`
	MaxHP       int       `json:"maxHP"`
	MaxATK      int       `json:"maxATK"`
	MaxTECH     int       `json:"maxTECH"`
	SpecialType string    `json:"specialType"`
}

type EquipmentDetail struct {
	EquipmentID   uuid.UUID `json:"equipmentID"`
	BonusHP       int       `json:"bonusHP"`
	BonusATK      int       `json:"bonusATK"`
	BonusTECH     int       `json:"bonusTECH"`
	InitBonusHP   int       `json:"initBonusHP"`
	InitBonusATK  int       `json:"initBonusATK"`
	InitBonusTECH int       `json:"initBonusTECH"`
	MaxBonusHP    int       `json:"maxBonusHP"`
	MaxBonusATK   int       `json:"maxBonusATK"`
	MaxBonusTECH  int       `json:"maxBonusTECH"`
	BuffEffect    *string   `json:"buffEffect,omitempty"`
}

type CardInstanceDetail struct {
	InstanceID uuid.UUID        `json:"instanceID"`
	Level      int              `json:"level"`
	Card       CardMaster       `json:"card"`
	Character  *CharacterDetail `json:"character,omitempty"`
	Equipment  *EquipmentDetail `json:"equipment,omitempty"`
}
