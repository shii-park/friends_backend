package errs

import "errors"

var (
	// ===== User =====
	ErrUserNameRequired       = errors.New("ユーザー名は必須です")
	ErrUserNotFound           = errors.New("ユーザーが存在しません")
	ErrInvalidBirthday        = errors.New("無効な日付です")
	ErrInvalidRankPointDelta  = errors.New("無効なランクポイントの増加量です")
	ErrInvalidCoinDelta       = errors.New("無効な通貨の増加量です")
	ErrInsufficientCoin       = errors.New("通貨の量が不十分です")
	ErrInvalidGachaStoneDelta = errors.New("無効なガチャ石の増加量です")
	ErrInsufficientGachaStone = errors.New("ガチャ石の量が不十分です")

	// ===== Card =====
	ErrCardIDRequired    = errors.New("カードIDは必須です")
	ErrCardNameRequired  = errors.New("カード名は必須です")
	ErrInvalidCardRarity = errors.New("無効なカードのレア度です")
	ErrInvalidCardKind   = errors.New("無効なカードの種類です")
)
