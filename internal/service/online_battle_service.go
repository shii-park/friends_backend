package service

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/shii-park/friends/internal/domain"
	"github.com/shii-park/friends/internal/errs"
	"github.com/shii-park/friends/internal/sqlc"
)

// WaitingPlayer はマッチング待ちのプレイヤー情報
type WaitingPlayer struct {
	Conn   *websocket.Conn
	wmu    sync.Mutex
	UserID uuid.UUID
	Chara  *domain.Character
	Equip  *domain.Equip
	RoomCh chan *OnlineBattleRoom
}

func NewWaitingPlayer(conn *websocket.Conn, userID uuid.UUID, chara *domain.Character, equip *domain.Equip) *WaitingPlayer {
	return &WaitingPlayer{
		Conn:   conn,
		UserID: userID,
		Chara:  chara,
		Equip:  equip,
		RoomCh: make(chan *OnlineBattleRoom, 1),
	}
}

func (p *WaitingPlayer) WriteJSON(v interface{}) error {
	p.wmu.Lock()
	defer p.wmu.Unlock()
	return p.Conn.WriteJSON(v)
}

// Matchmaker はマッチングを管理する
type Matchmaker struct {
	mu      sync.Mutex
	waiting *WaitingPlayer
	queries *sqlc.Queries
}

func NewMatchmaker(queries *sqlc.Queries) *Matchmaker {
	return &Matchmaker{queries: queries}
}

// Enqueue はプレイヤーをマッチングキューに追加する。
// nil=待機中, non-nil=即マッチ（返り値は対戦相手）
func (m *Matchmaker) Enqueue(p *WaitingPlayer) *WaitingPlayer {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.waiting == nil {
		m.waiting = p
		return nil
	}

	opponent := m.waiting
	m.waiting = nil
	return opponent
}

// Dequeue はタイムアウト時などにキューからプレイヤーを除去する
func (m *Matchmaker) Dequeue(p *WaitingPlayer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.waiting == p {
		m.waiting = nil
	}
}

// CreateRoom は2プレイヤーのバトルルームを作成する
func (m *Matchmaker) CreateRoom(playerA, playerB *WaitingPlayer) *OnlineBattleRoom {
	return NewOnlineBattleRoom(playerA, playerB, m.queries)
}

// OnlinePlayerMsg は各プレイヤーへ送るラウンド結果のデータ
type OnlinePlayerMsg struct {
	PlayerHP       int
	NpcHP          int
	NpcHand        domain.AttackType
	Outcome        string
	RankPointDelta int
	CoinReward     int
	StoneReward    int
}

// OnlineRoundResult はSubmitHandの結果
type OnlineRoundResult struct {
	IsOver bool
	ForA   OnlinePlayerMsg
	ForB   OnlinePlayerMsg
}

// OnlineBattleRoom は2人のプレイヤーのバトルルーム
type OnlineBattleRoom struct {
	WriterA *WaitingPlayer
	WriterB *WaitingPlayer
	battle  *domain.Battle
	queries *sqlc.Queries
	mu      sync.Mutex
	handA   *domain.AttackType
	handB   *domain.AttackType
	ended   bool
}

func NewOnlineBattleRoom(playerA, playerB *WaitingPlayer, queries *sqlc.Queries) *OnlineBattleRoom {
	battle := domain.NewBattle(
		playerA.Chara, playerB.Chara,
		playerA.Equip, playerB.Equip,
		playerA.Chara.HP, playerB.Chara.HP,
	)
	return &OnlineBattleRoom{
		WriterA: playerA,
		WriterB: playerB,
		battle:  battle,
		queries: queries,
	}
}

// Battle はバトル状態を返す（ready メッセージ送信用）
func (r *OnlineBattleRoom) Battle() *domain.Battle {
	return r.battle
}

// IsEnded はゲームが正常終了したかを返す（切断通知の抑制に使用）
func (r *OnlineBattleRoom) IsEnded() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.ended
}

// SubmitHand は手を提出し、両者揃ったらA・B向けの結果を返す。
// nil=まだ片方のみ提出
func (r *OnlineBattleRoom) SubmitHand(ctx context.Context, isA bool, hand domain.AttackType) (*OnlineRoundResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if isA {
		r.handA = &hand
	} else {
		r.handB = &hand
	}

	if r.handA == nil || r.handB == nil {
		return nil, nil
	}

	handA := *r.handA
	handB := *r.handB
	r.handA = nil
	r.handB = nil

	// じゃんけん判定とダメージ計算
	jankenResult := domain.JudgeJanken(handA, handB)
	if jankenResult == domain.Win {
		r.battle.PlayerAttack(handA)
	} else if jankenResult == domain.Lose {
		r.battle.OppoAttack(handB)
	}

	playerHP := r.battle.PlayerHP
	oppoHP := r.battle.OppoHP

	if playerHP <= 0 || oppoHP <= 0 {
		r.ended = true
		outcomeA, outcomeB, deltaA, deltaB, coinA, stoneA, coinB, stoneB, err := r.finishOnlineBattle(ctx, playerHP, oppoHP)
		if err != nil {
			return nil, err
		}
		return &OnlineRoundResult{
			IsOver: true,
			ForA: OnlinePlayerMsg{
				PlayerHP:       playerHP,
				NpcHP:          oppoHP,
				NpcHand:        handB,
				Outcome:        outcomeA,
				RankPointDelta: deltaA,
				CoinReward:     coinA,
				StoneReward:    stoneA,
			},
			ForB: OnlinePlayerMsg{
				PlayerHP:       oppoHP,
				NpcHP:          playerHP,
				NpcHand:        handA,
				Outcome:        outcomeB,
				RankPointDelta: deltaB,
				CoinReward:     coinB,
				StoneReward:    stoneB,
			},
		}, nil
	}

	return &OnlineRoundResult{
		IsOver: false,
		ForA: OnlinePlayerMsg{
			PlayerHP: playerHP,
			NpcHP:    oppoHP,
			NpcHand:  handB,
		},
		ForB: OnlinePlayerMsg{
			PlayerHP: oppoHP,
			NpcHP:    playerHP,
			NpcHand:  handA,
		},
	}, nil
}

func (r *OnlineBattleRoom) finishOnlineBattle(ctx context.Context, playerHP, oppoHP int) (
	outcomeA, outcomeB string, deltaA, deltaB, coinA, stoneA, coinB, stoneB int, err error,
) {
	const winDelta = 10
	const loseDelta = 5
	const winCoin = 100
	const winStone = 1

	switch {
	case playerHP > 0 && oppoHP <= 0:
		outcomeA = string(domain.Win)
		outcomeB = string(domain.Lose)
		deltaA = winDelta
		deltaB = -loseDelta
		coinA = winCoin
		stoneA = winStone
	case playerHP <= 0 && oppoHP > 0:
		outcomeA = string(domain.Lose)
		outcomeB = string(domain.Win)
		deltaA = -loseDelta
		deltaB = winDelta
		coinB = winCoin
		stoneB = winStone
	default:
		outcomeA = string(domain.Draw)
		outcomeB = string(domain.Draw)
	}

	if updateErr := r.updateUserRankPoint(ctx, r.WriterA.UserID, deltaA); updateErr != nil {
		err = updateErr
		return
	}
	if updateErr := r.updateUserRankPoint(ctx, r.WriterB.UserID, deltaB); updateErr != nil {
		err = updateErr
		return
	}

	if outcomeA == string(domain.Win) {
		if rewardErr := r.giveRewards(ctx, r.WriterA.UserID, coinA, stoneA); rewardErr != nil {
			err = rewardErr
			return
		}
	}
	if outcomeB == string(domain.Win) {
		if rewardErr := r.giveRewards(ctx, r.WriterB.UserID, coinB, stoneB); rewardErr != nil {
			err = rewardErr
			return
		}
	}

	return
}

func (r *OnlineBattleRoom) updateUserRankPoint(ctx context.Context, userID uuid.UUID, delta int) error {
	dbUser, err := r.queries.GetUser(ctx, userID)
	if err != nil {
		return errs.ErrUserNotFound
	}

	newRankPoint := int(dbUser.RankPoint) + delta
	if newRankPoint < 0 {
		newRankPoint = 0
	}

	return r.queries.UpdateUserRankPoint(ctx, sqlc.UpdateUserRankPointParams{
		UserID:    userID,
		RankPoint: int32(newRankPoint),
	})
}

func (r *OnlineBattleRoom) giveRewards(ctx context.Context, userID uuid.UUID, coin, stone int) error {
	if coin > 0 {
		if _, err := r.queries.AddUserCoin(ctx, sqlc.AddUserCoinParams{
			Amount: int32(coin),
			UserID: userID,
		}); err != nil {
			return err
		}
	}
	if stone > 0 {
		if _, err := r.queries.AddUserGachaStone(ctx, sqlc.AddUserGachaStoneParams{
			Amount: int32(stone),
			UserID: userID,
		}); err != nil {
			return err
		}
	}
	return nil
}
