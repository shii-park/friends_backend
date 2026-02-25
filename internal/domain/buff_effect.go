package domain

type BuffEffect interface {
	Apply(chara *Character)

	Name() string
}
