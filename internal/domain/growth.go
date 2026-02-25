package domain

import "math"

func CalcLinear(init, max, level, maxLevel int) int {
	if maxLevel <= 1 {
		return int(init)
	}
	if level < 1 {
		level = 1
	}
	if level > maxLevel {
		level = maxLevel
	}

	r := float64(level-1) / float64(maxLevel-1)
	v := float64(init) + (float64(max)-float64(init))*r
	return int(math.Round(v))
}
