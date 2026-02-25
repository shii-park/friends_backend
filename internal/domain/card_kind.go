package domain

type CardKind int

const (
	CardKindEquip     CardKind = 0
	CardKindCharacter CardKind = 1
)

func (k CardKind) IsValid() bool {
	switch k {
	case CardKindEquip, CardKindCharacter:
		return true
	default:
		return false
	}
}
