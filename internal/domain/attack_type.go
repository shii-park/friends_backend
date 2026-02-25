package domain

type AttackType string

const (
	Rock     AttackType = "rock"
	Scissors AttackType = "scissors"
	Paper    AttackType = "paper"
)

func JudgeJanken(player, oppo AttackType) JankenResult {
	if player == oppo {
		return Draw
	}
	if (player == Rock && oppo == Scissors) ||
		(player == Paper && oppo == Rock) ||
		(player == Scissors && oppo == Paper) {
		return Win
	}
	return Lose
}
