package domain

import "github.com/shii-park/friends/internal/errs"

// LevelUpgradable 強化（レベル変更）できる対象
type LevelUpgradable interface {
	GetLevel() int
	SetLevel(level int) error
}

type EnhancementService struct{}

// EnhanceByCoin コインを消費して強化する（最大 times レベル上げる）
func (s EnhancementService) EnhanceByCoinTable(
	user *User,
	target LevelUpgradable,
	times int,
) (leveledUp int, err error) {
	if times <= 0 {
		return 0, errs.ErrInvalidEnhanceTimes
	}

	current := target.GetLevel()
	if current < InitialLevel {
		current = InitialLevel
	}

	// もう最大
	if current >= MaxLevel {
		return 0, nil
	}

	// 上限までに収める
	remain := MaxLevel - current
	if times > remain {
		times = remain
	}

	// コスト計算（enhance_rules.go の関数）
	totalCost := TotalCoinCostForLevelUps(current, times)
	if totalCost <= 0 {
		return 0, errs.ErrInvalidEnhanceCost
	}

	// 支払い
	if err := user.ConsumeCoin(totalCost); err != nil {
		return 0, err
	}

	// レベル更新（SetLevel 内でステが再計算される想定）
	if err := target.SetLevel(current + times); err != nil {
		return 0, err
	}

	return times, nil
}
