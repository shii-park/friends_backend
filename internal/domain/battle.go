package domain

type Battle struct {
	PlayerCard  *Character
	OppoCard    *Character
	PlayerEquip *Equip
	OppoEquip   *Equip
	PlayerHP    int
	OppoHP      int
}

func NewBattle(playerCard *Character, oppoCard *Character, playerEquip *Equip, oppoEquip *Equip, playerHP int, oppoHP int) *Battle {
	return &Battle{PlayerCard: playerCard,
		OppoCard:    oppoCard,
		PlayerEquip: playerEquip,
		OppoEquip:   oppoEquip,
		PlayerHP:    playerHP,
		OppoHP:      oppoHP}
}
