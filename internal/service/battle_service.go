package service

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/domain"
	"github.com/shii-park/friends/internal/errs"
	"github.com/shii-park/friends/internal/sqlc"
)

type BattleService struct {
	queries *sqlc.Queries
}

func NewBattleService(queries *sqlc.Queries) *BattleService {
	return &BattleService{queries: queries}
}

// BattleSession はWebSocketコネクション中に保持するバトル状態
type BattleSession struct {
	Battle *domain.Battle
	UserID uuid.UUID
}

// RoundResult は1ラウンドの結果
type RoundResult struct {
	PlayerHP int
	NpcHP    int
	NpcHand  domain.AttackType
	IsOver   bool
	Outcome  string
	RankPointDelta int
}

// BattleResult はゲーム終了時の結果
type BattleResult struct {
	Outcome        string `json:"outcome"`
	RankPointDelta int    `json:"rankPointDelta"`
}

var attackTypes = []domain.AttackType{domain.Rock, domain.Scissors, domain.Paper}

func randomAttackType() domain.AttackType {
	return attackTypes[rand.IntN(len(attackTypes))]
}

// StartBattle はバトルを初期化してセッションを返す
func (s *BattleService) StartBattle(ctx context.Context, userID uuid.UUID, charaID, equipID uuid.UUID) (*BattleSession, error) {
	// 自分のキャラクター・装備を取得する
	dbPlayerChara, err := s.queries.GetCharacter(ctx, charaID)
	if err != nil {
		return nil, errs.ErrInvalidCharaID
	}
	dbPlayerEquip, err := s.queries.GetEquipment(ctx, equipID)
	if err != nil {
		return nil, errs.ErrInvalidEquipID
	}

	// 相手NPCのキャラクター・装備をランダムで選ぶ
	allCharas, err := s.queries.GetAllCharacters(ctx)
	if err != nil || len(allCharas) == 0 {
		return nil, errs.ErrInvalidCharaID
	}
	allEquips, err := s.queries.GetAllEquipments(ctx)
	if err != nil || len(allEquips) == 0 {
		return nil, errs.ErrInvalidEquipID
	}
	dbNpcChara := allCharas[rand.IntN(len(allCharas))]
	dbNpcEquip := allEquips[rand.IntN(len(allEquips))]

	// sqlc型をdomain型に変換
	playerChara := todomainChara(dbPlayerChara)
	playerEquip := todomainEquip(dbPlayerEquip)
	npcChara := todomainChara(dbNpcChara)
	npcEquip := todomainEquip(dbNpcEquip)

	if playerChara == nil || npcChara == nil || playerEquip == nil || npcEquip == nil {
		return nil, errs.ErrInvalidCharaID
	}

	battle := domain.NewBattle(playerChara, npcChara, playerEquip, npcEquip, playerChara.HP, npcChara.HP)

	return &BattleSession{
		Battle: battle,
		UserID: userID,
	}, nil
}

// RoundBattle は1ラウンドの攻撃処理を行い結果を返す
func (s *BattleService) RoundBattle(ctx context.Context, session *BattleSession, playerHand domain.AttackType) (*RoundResult, error) {
	npcHand := randomAttackType()

	session.Battle.PlayerAttack(playerHand)
	session.Battle.OppoAttack(npcHand)

	playerHP := session.Battle.PlayerHP
	npcHP := session.Battle.OppoHP

	// どちらかのHPが0以下になったらゲーム終了
	if playerHP <= 0 || npcHP <= 0 {
		outcome, delta, err := s.finishBattle(ctx, session, playerHP, npcHP)
		if err != nil {
			return nil, err
		}
		return &RoundResult{
			PlayerHP:       playerHP,
			NpcHP:          npcHP,
			NpcHand:        npcHand,
			IsOver:         true,
			Outcome:        outcome,
			RankPointDelta: delta,
		}, nil
	}

	return &RoundResult{
		PlayerHP: playerHP,
		NpcHP:    npcHP,
		NpcHand:  npcHand,
		IsOver:   false,
	}, nil
}

// finishBattle は勝敗を判定しランクポイントをDBに保存する
func (s *BattleService) finishBattle(ctx context.Context, session *BattleSession, playerHP, npcHP int) (string, int, error) {
	const winDelta = 10
	const loseDelta = 5

	var outcome string
	var delta int

	switch {
	case playerHP > 0 && npcHP <= 0:
		outcome = string(domain.Win)
		delta = winDelta
	case playerHP <= 0 && npcHP > 0:
		outcome = string(domain.Lose)
		delta = -loseDelta
	default:
		outcome = string(domain.Draw)
		delta = 0
	}

	// 現在のランクポイントを取得してから更新する
	dbUser, err := s.queries.GetUser(ctx, session.UserID)
	if err != nil {
		return "", 0, errs.ErrUserNotFound
	}

	newRankPoint := int(dbUser.RankPoint) + delta
	if newRankPoint < 0 {
		newRankPoint = 0
	}

	err = s.queries.UpdateUserRankPoint(ctx, sqlc.UpdateUserRankPointParams{
		UserID:    session.UserID,
		RankPoint: int32(newRankPoint),
	})
	if err != nil {
		return "", 0, err
	}

	return outcome, delta, nil
}

func todomainChara(c sqlc.Character) *domain.Character {
	chara, _ := domain.NewCharacter(
		fmt.Sprintf("%d", c.CardID),
		c.CharacterID.String(),
		"",
		"",
		domain.Rarity(""),
		int(c.InitHp), int(c.InitAtk), int(c.InitTech),
		int(c.MaxHp), int(c.MaxAtk), int(c.MaxTech),
		domain.SpecialType(c.SpecialType.String),
	)
	return chara
}

func todomainEquip(e sqlc.Equipment) *domain.Equip {
	equip, _ := domain.NewEquip(
		fmt.Sprintf("%d", e.CardID),
		e.EquipmentID.String(),
		"", "",
		domain.Rarity(""),
		int(e.InitBonusHp), int(e.InitBonusAtk), int(e.InitBonusTech),
		int(e.MaxBonusHp), int(e.MaxBonusAtk), int(e.MaxBonusTech),
		nil,
	)
	return equip
}
