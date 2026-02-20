package errs

import "errors"

var (
	// ===== User =====
	ErrUserNameRequired      = errors.New("ユーザーネームは必須です")
	ErrUserNotFound          = errors.New("ユーザーが存在しません")
	ErrInvalidBirthday       = errors.New("無効な日付です")
	ErrInvalidRankPointDelta = errors.New("無効なランクポイントの増加量です")
	ErrInsufficientRankPoint = errors.New("")
)
