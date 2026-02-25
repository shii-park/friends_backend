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

	// ===== Chara =====
	ErrCharaIDRequired         = errors.New("キャラクターIDは必須です")
	ErrInvalidCharaID          = errors.New("無効なキャラクターIDです")
	ErrInvalidCharacterHP      = errors.New("無効なHPです")
	ErrInvalidCharacterATK     = errors.New("無効なATKです")
	ErrInvalidCharacterTECH    = errors.New("無効なTECKです")
	ErrInvalidCharacterMaxHP   = errors.New("無効な最大HPです")
	ErrInvalidCharacterMaxATK  = errors.New("無効な最大ATKです")
	ErrInvalidCharacterMaxTECH = errors.New("無効な最大TECKです")
	ErrInvalidSpecialType      = errors.New("無効な必殺技の種類です")

	// ===== Equip =====
	ErrEquipIDRequired     = errors.New("キャラクターIDは必須です")
	ErrInvalidEquipID      = errors.New("無効なキャラクターIDです")
	ErrInvalidEquipHP      = errors.New("無効な付加HPです")
	ErrInvalidEquipATK     = errors.New("無効な付加ATKです")
	ErrInvalidEquipTECH    = errors.New("無効な付加TECKです")
	ErrInvalidEquipMaxHP   = errors.New("無効な最大付加HPです")
	ErrInvalidEquipMaxATK  = errors.New("無効な最大付加ATKです")
	ErrInvalidEquipMaxTECH = errors.New("無効な最大付加TECKです")
)
