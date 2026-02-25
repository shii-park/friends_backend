package domain

import (
	"math"
	"math/rand/v2"
)

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
		OppoHP:      oppoHP,
	}
}

func (b *Battle) PlayerAttack(hand AttackType) {
	multATK := rand.Float64()*(1.0-0.8) + 0.8
	multTECH := rand.Float64()*(1.0-0.5) + 0.5
	multDEF := rand.Float64()*(0.5-0.3) + 0.3
	if b.isPlayerSpecialType(hand) {
		// 攻撃が必殺技だった場合
		dmg := float64(b.PlayerCard.ATK+b.PlayerEquip.BonusATK)*multATK +
			float64(b.PlayerCard.TECH+b.PlayerEquip.BonusTECH)*multTECH
		def := float64(b.OppoCard.ATK+b.OppoEquip.BonusATK) * multDEF
		b.OppoHP = int(math.Max(0, float64(b.OppoHP)-(dmg-def)))
	} else {
		// 攻撃が通常技だった場合
		dmg := float64(b.PlayerCard.ATK+b.PlayerEquip.BonusATK) * multATK
		def := float64(b.OppoCard.ATK+b.OppoEquip.BonusATK) * multDEF
		b.OppoHP = int(math.Max(0, float64(b.OppoHP)-(dmg-def)))
	}
}

func (b *Battle) OppoAttack(hand AttackType) {
	multATK := rand.Float64()*(1.0-0.8) + 0.8
	multTECH := rand.Float64()*(1.0-0.5) + 0.5
	multDEF := rand.Float64()*(0.5-0.3) + 0.3
	if b.isOppoSpecialType(hand) {
		// 攻撃が必殺技だった場合
		dmg := float64(b.OppoCard.ATK+b.OppoEquip.BonusATK)*multATK +
			float64(b.OppoCard.TECH+b.OppoEquip.BonusTECH)*multTECH
		def := float64(b.PlayerCard.ATK+b.PlayerEquip.BonusATK) * multDEF
		b.PlayerHP = int(math.Max(0, float64(b.PlayerHP)-(dmg-def)))
	} else {
		// 攻撃が通常技だった場合
		dmg := float64(b.OppoCard.ATK+b.OppoEquip.BonusATK) * multATK
		def := float64(b.PlayerCard.ATK+b.PlayerEquip.BonusATK) * multDEF
		b.PlayerHP = int(math.Max(0, float64(b.PlayerHP)-(dmg-def)))
	}
}

// isPlayerSpecialTypeは自分の出した手が必殺技の手であるかを判定します
func (b *Battle) isPlayerSpecialType(hand AttackType) bool {
	return hand == AttackType(b.PlayerCard.SpecialType)
}

// isOppoSpecialTypeは相手の出した手が必殺技の手であるかを判定します
func (b *Battle) isOppoSpecialType(hand AttackType) bool {
	return hand == AttackType(b.OppoCard.SpecialType)
}

// IsGameOverはキャラクラーが死んでいるかを判定します
// 死んでいればtrue, 生きていればfalseを返します
func (b *Battle) IsGameOver(chara *Character) bool {
	return chara.HP <= 0
}
