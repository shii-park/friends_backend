package errs

import "errors"

var (
	// ===== User =====
	ErrUserNameRequired       = errors.New("ユーザー名は必須です")
	ErrUserNotFound           = errors.New("ユーザーが存在しません")
	ErrInvalidUserID          = errors.New("無効なユーザーIDです")
	ErrInvalidBirthday        = errors.New("無効な日付です")
	ErrInvalidRankPointDelta  = errors.New("無効なランクポイントの増加量です")
	ErrInvalidRankPoint       = errors.New("無効なランクポイントです")
	ErrInvalidCoinDelta       = errors.New("無効な通貨の増加量です")
	ErrInsufficientCoin       = errors.New("通貨の量が不十分です")
	ErrInvalidGachaStoneDelta = errors.New("無効なガチャ石の増加量です")
	ErrInsufficientGachaStone = errors.New("ガチャ石の量が不十分です")

	// ===== Card =====
	ErrCardIDRequired      = errors.New("カードIDは必須です")
	ErrCardNameRequired    = errors.New("カード名は必須です")
	ErrInvalidCardID       = errors.New("無効なカードIDです")
	ErrInvalidCardRarity   = errors.New("無効なカードのレア度です")
	ErrInvalidCardKind     = errors.New("無効なカードの種類です")
	ErrInvalidEnhanceTimes = errors.New("無効な強化レベルです")
	ErrInvalidEnhanceCost  = errors.New("無効な強化コストの値です")

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
	ErrInvalidCharaHPDelta     = errors.New("無効なHPの増加量です")
	ErrInvalidCharaATKDelta    = errors.New("無効なATKの増加量です")
	ErrInvalidCharaTECHDelta   = errors.New("無効なTECHの増加量です")
	ErrCharaAlreadyMaxLevel    = errors.New("既に最大レベルです")
	ErrInvalidCharaLevel       = errors.New("不正なレベルです")

	// ===== Equip =====
	ErrEquipIDRequired       = errors.New("装備IDは必須です")
	ErrInvalidEquipID        = errors.New("無効な装備IDです")
	ErrInvalidEquipHP        = errors.New("無効な付加HPです")
	ErrInvalidEquipATK       = errors.New("無効な付加ATKです")
	ErrInvalidEquipTECH      = errors.New("無効な付加TECKです")
	ErrInvalidEquipMaxHP     = errors.New("無効な最大付加HPです")
	ErrInvalidEquipMaxATK    = errors.New("無効な最大付加ATKです")
	ErrInvalidEquipMaxTECH   = errors.New("無効な最大付加TECKです")
	ErrInvalidEquipHPDelta   = errors.New("無効なHPの増加量です")
	ErrInvalidEquipATKDelta  = errors.New("無効なATKの増加量です")
	ErrInvalidEquipTECHDelta = errors.New("無効なTECHの増加量です")
	ErrEquipAlreadyMaxLevel  = errors.New("既に最大レベルです")
	ErrInvalidEquipLevel     = errors.New("不正なレベルです")
	ErrInvalidInstanceID     = errors.New("無効なInstanceIDです")
	ErrStorageNil            = errors.New("ストレージには何もありません")
	ErrCardAlreadyHas        = errors.New("既に所持しているカードです")
	ErrCardNotFound          = errors.New("カードが見つかりません")

	// ===== queue ======
	ErrNoUsersInQueue = errors.New("キューに人がいませんでした")
)
