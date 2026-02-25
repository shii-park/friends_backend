package domain

import (
	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/errs"
)

// CardInstance 所持カード1枚を表す
type CardInstance struct {
	InstanceID uuid.UUID
	CardID     int
}

// Storage ユーザーの所持カード一覧
type Storage struct {
	UserID string

	// instanceID → CardInstance
	cards map[uuid.UUID]CardInstance
}

// NewEmptyStorage 新規作成（空）
func NewEmptyStorage(userID string) *Storage {
	return &Storage{
		UserID: userID,
		cards:  make(map[uuid.UUID]CardInstance),
	}
}

// ReconstructStorage DBから復元
func ReconstructStorage(userID string, cardInstances []CardInstance) *Storage {
	m := make(map[uuid.UUID]CardInstance, len(cardInstances))
	for _, ci := range cardInstances {
		m[ci.InstanceID] = ci
	}
	return &Storage{
		UserID: userID,
		cards:  m,
	}
}

// Add Card追加（instanceID を新規生成して追加）
func (s *Storage) Add(cardID int) (CardInstance, error) {
	if s == nil {
		return CardInstance{}, errs.ErrStorageNil
	}
	if cardID <= 0 {
		return CardInstance{}, errs.ErrInvalidCardID
	}

	instanceID := uuid.New()

	ci := CardInstance{
		InstanceID: instanceID,
		CardID:     cardID,
	}

	s.cards[instanceID] = ci
	return ci, nil
}

// Remove InstanceIDで削除
func (s *Storage) Remove(instanceID uuid.UUID) error {
	if s == nil {
		return errs.ErrStorageNil
	}
	if instanceID == uuid.Nil {
		return errs.ErrInvalidInstanceID
	}
	if _, ok := s.cards[instanceID]; !ok {
		return errs.ErrCardNotFound
	}
	delete(s.cards, instanceID)
	return nil
}

// Has InstanceID存在確認
func (s *Storage) Has(instanceID uuid.UUID) bool {
	if s == nil {
		return false
	}
	_, ok := s.cards[instanceID]
	return ok
}

// FindByCardID cardIDで検索
func (s *Storage) FindByCardID(cardID int) []CardInstance {
	if s == nil {
		return nil
	}
	result := []CardInstance{}
	for _, v := range s.cards {
		if v.CardID == cardID {
			result = append(result, v)
		}
	}
	return result
}

// AllList 一覧取得
func (s *Storage) AllList() []CardInstance {
	if s == nil {
		return nil
	}
	result := make([]CardInstance, 0, len(s.cards))
	for _, v := range s.cards {
		result = append(result, v)
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
