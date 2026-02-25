package domain

type CardState string

const (
	NotFound CardState = "notFound"
	Find     CardState = "find"
	Get      CardState = "get"
)
