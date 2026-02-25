package domain

type BuffEffect interface {
	Apply(chara *Character)

	Remove(chara *Character)

	Name() string
}
