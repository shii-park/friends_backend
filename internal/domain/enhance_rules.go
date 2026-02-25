package domain

// NextLevelCoinCost 次のレベルへ上げるためのコインコストを返す
func NextLevelCoinCost(currentLevel int) int {
	if currentLevel < InitialLevel {
		currentLevel = InitialLevel
	}
	if currentLevel >= MaxLevel {
		return 0
	}

	next := currentLevel + 1

	switch next {
	case 2:
		return 100
	case 3:
		return 200
	case 4:
		return 300
	case 5:
		return 400
	case 6:
		return 500
	case 7:
		return 600
	case 8:
		return 700
	case 9:
		return 800
	case 10:
		return 10000
	default:
		return 0
	}
}

// CanEnhance 強化できるか
func CanEnhance(card *BaseCard) bool {
	if card == nil {
		return false
	}
	lv := card.Level
	if lv < InitialLevel {
		lv = InitialLevel
	}
	return lv < MaxLevel
}

func TotalCoinCostForLevelUps(currentLevel int, times int) int {
	if times <= 0 {
		return 0
	}
	if currentLevel < InitialLevel {
		currentLevel = InitialLevel
	}

	total := 0
	lv := currentLevel

	for i := 0; i < times; i++ {
		cost := NextLevelCoinCost(lv)
		if cost == 0 {
			break // Max か不正
		}
		total += cost
		lv++
	}

	return total
}
