package domain

import "github.com/shii-park/friends/internal/errs"

// LevelUpgradable: 強化（レベル変更）できる対象
// Character / Equip がこれを満たせばOK
type LevelUpgradable interface {
	GetLevel() int
	SetLevel(level int) error
}

type EnhancementService struct{}

// EnhanceByCoin: コインを消費して強化する（最大 times レベル上げる）
func (s EnhancementService) EnhanceByCoin(
	user *User,
	target LevelUpgradable,
	times int,
	costPerLevel int,
) (leveledUp int, err error) {
	if times <= 0 {
		return 0, errs.ErrInvalidEnhanceTimes
	}
	if costPerLevel <= 0 {
		return 0, errs.ErrInvalidEnhanceCost
	}

	current := target.GetLevel()
	if current < 1 {
		current = 1
	}

	if current >= MaxLevel {
		return 0, nil
	}

	remain := MaxLevel - current
	if times > remain {
		times = remain
	}

	totalCost := costPerLevel * times
	if err := user.ConsumeCoin(totalCost); err != nil {
		return 0, err
	}

	if err := target.SetLevel(current + times); err != nil {
		return 0, err
	}

	return times, nil
}
