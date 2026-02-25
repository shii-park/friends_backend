package domain

import (
	"github.com/shii-park/friends/internal/errs"
)

// 型定義
type (
	UserID string
	CardID string
)

// Storage はユーザーの所持カード一覧
// instanceの概念は無く、card_id の集合として管理する
type Storage struct {
	UserID UserID

	// 所持カード
	cards map[CardID]struct{}
}

// NewEmptyStorage 新規作成（空）
func NewEmptyStorage(userID UserID) *Storage {
	return &Storage{
		UserID: userID,
		cards:  make(map[CardID]struct{}),
	}
}

// ReconstructStorage DBから復元
func ReconstructStorage(
	userID UserID,
	cardIDs []CardID,
) *Storage {
	m := make(map[CardID]struct{}, len(cardIDs))

	for _, id := range cardIDs {
		m[id] = struct{}{}
	}

	return &Storage{
		UserID: userID,
		cards:  m,
	}
}

// Add Card追加
// 既に持っていたらエラー
func (s *Storage) Add(cardID CardID) error {
	if s == nil {
		return errs.ErrStorageNil
	}

	if cardID == "" {
		return errs.ErrInvalidCardID
	}

	if s.Has(cardID) {
		return errs.ErrCardAlreadyHas
	}

	s.cards[cardID] = struct{}{}

	return nil
}

// Remove Card削除
func (s *Storage) Remove(cardID CardID) error {
	if s == nil {
		return errs.ErrStorageNil
	}

	if !s.Has(cardID) {
		return errs.ErrCardNotFound
	}

	delete(s.cards, cardID)

	return nil
}

// Has 所持確認
func (s *Storage) Has(cardID CardID) bool {
	if s == nil {
		return false
	}

	_, ok := s.cards[cardID]

	return ok
}

// List 一覧取得
func (s *Storage) List() []CardID {
	if s == nil {
		return nil
	}

	result := make([]CardID, 0, len(s.cards))

	for id := range s.cards {
		result = append(result, id)
	}

	return result
}

// Count 枚数
func (s *Storage) Count() int {
	if s == nil {
		return 0
	}

	return len(s.cards)
}
