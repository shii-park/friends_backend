package domain

import "time"

type Card interface {
	GetID() string
	GetName() string
	GetIcon() string
	GetRarity() Rarity
	GetAcquiredDate() time.Time
	GetLevel() int
	Kind() CardKind
}
