package matching

import (
	"github.com/shii-park/friends/internal/errs"
)

type Queue struct {
	users []string
}

// 2. Queue構造体のポインタを返す
func NewQueue() *Queue {
	return &Queue{
		users: []string{}, // 初期化
	}
}

func (q *Queue) Enqueue(userID string) {
	q.users = append(q.users, userID)
}

func (q *Queue) Dequeue() (string, error) {
	if len(q.users) == 0 {
		return "", errs.ErrNoUsersInQueue
	}

	userID := q.users[0]

	q.users = q.users[1:]

	return userID, nil
}
