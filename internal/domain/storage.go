package domain

import (
	"github.com/shii-park/friends/internal/errs"
)

// 型定義

// Storage はユーザーの所持カード一覧
// instanceの概念は無く、card_id の集合として管理する
type Storage struct {
	UserID string

	// 所持カード
	cards map[string]struct{}
}

// NewEmptyStorage 新規作成（空）
func NewEmptyStorage(userID string) *Storage {
	return &Storage{
		UserID: userID,
		cards:  make(map[string]struct{}),
	}
}

// ReconstructStorage DBから復元
func ReconstructStorage(
	userID string,
	cardIDs []string,
) *Storage {
	m := make(map[string]struct{}, len(cardIDs))

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
func (s *Storage) Add(cardID string) error {
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
func (s *Storage) Remove(cardID string) error {
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
func (s *Storage) Has(cardID string) bool {
	if s == nil {
		return false
	}

	_, ok := s.cards[cardID]

	return ok
}

// AllList 一覧取得
func (s *Storage) AllList() []string {
	if s == nil {
		return nil
	}

	result := make([]string, 0, len(s.cards))

	for id := range s.cards {
		result = append(result, id)
	}

	return result
}

// AllCount 枚数
func (s *Storage) AllCount() int {
	if s == nil {
		return 0
	}

	return len(s.cards)
}
