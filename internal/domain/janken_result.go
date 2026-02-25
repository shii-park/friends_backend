package domain

type JankenResult string

const (
	Win  JankenResult = "win"
	Lose JankenResult = "lose"
	Draw JankenResult = "draw"
)
