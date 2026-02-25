package domain

type CardKind string

const (
	CardKindUnknown   CardKind = ""
	CardKindCharacter CardKind = "chara"
	CardKindEquip     CardKind = "equip"
)

func (k CardKind) IsValid() bool {
	switch k {
	case CardKindCharacter, CardKindEquip:
		return true
	default:
		return false
	}
}
