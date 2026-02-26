package service

import (
	"context"
	"errors"
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
	Battle   *domain.Battle
	UserID   uuid.UUID
	NpcChara *domain.Character
	NpcEquip *domain.Equip

	// 追加：次ラウンドでNPCが出す手を「先に確定」して保持する
	PendingNpcHand *domain.AttackType
}

// RoundResult は1ラウンドの結果
type RoundResult struct {
	PlayerHP       int
	NpcHP          int
	NpcHand        domain.AttackType
	IsOver         bool
	Outcome        string
	RankPointDelta int
	CoinReward     int
	StoneReward    int
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

// 追加：NPCの手に応じた「予告テキスト」生成
func npcHandText(hand domain.AttackType) string {
	switch hand {
	case domain.Rock:
		return "次はグーを出すぞ！"
	case domain.Paper:
		return "次はパーでいく！"
	case domain.Scissors:
		return "次はチョキだ！"
	default:
		return ""
	}
}

// 追加：prepare_round の戻り値（WSでテキストを返す用）
type PreparedRound struct {
	Text string
}

// 追加：NPCの手を先に決めてセッションに保持する
func (s *BattleService) PrepareRound(session *BattleSession) (*PreparedRound, error) {
	if session == nil || session.Battle == nil {
		return nil, errors.New("invalid battle session")
	}

	// すでに確定済みならそのまま（連打対策）
	if session.PendingNpcHand != nil {
		return &PreparedRound{Text: npcHandText(*session.PendingNpcHand)}, nil
	}

	hand := randomAttackType()
	session.PendingNpcHand = &hand

	return &PreparedRound{Text: npcHandText(hand)}, nil
}

// StartBattle はバトルを初期化してセッションを返す
func (s *BattleService) StartBattle(ctx context.Context, userID uuid.UUID, charaID, equipID uuid.UUID) (*BattleSession, error) {
	// 自分のキャラクター・装備を取得する（カード情報付き）
	dbPlayerChara, err := s.queries.GetCharacterWithCard(ctx, charaID)
	if err != nil {
		return nil, errs.ErrInvalidCharaID
	}
	dbPlayerEquip, err := s.queries.GetEquipmentWithCard(ctx, equipID)
	if err != nil {
		return nil, errs.ErrInvalidEquipID
	}

	// 相手NPCのキャラクター・装備をランダムで選ぶ（カード情報付き）
	allCharas, err := s.queries.GetAllCharactersWithCard(ctx)
	if err != nil || len(allCharas) == 0 {
		return nil, errs.ErrInvalidCharaID
	}
	allEquips, err := s.queries.GetAllEquipmentsWithCard(ctx)
	if err != nil || len(allEquips) == 0 {
		return nil, errs.ErrInvalidEquipID
	}
	dbNpcChara := allCharas[rand.IntN(len(allCharas))]
	dbNpcEquip := allEquips[rand.IntN(len(allEquips))]

	// sqlc型をdomain型に変換
	playerChara := todomainCharaWithCard(dbPlayerChara)
	playerEquip := todomainEquipWithCard(dbPlayerEquip)
	npcChara := todomainCharaWithCardFromMany(dbNpcChara)
	npcEquip := todomainEquipWithCardFromMany(dbNpcEquip)

	if playerChara == nil || npcChara == nil || playerEquip == nil || npcEquip == nil {
		return nil, errs.ErrInvalidCharaID
	}

	battle := domain.NewBattle(playerChara, npcChara, playerEquip, npcEquip, playerChara.HP, npcChara.HP)

	return &BattleSession{
		Battle:   battle,
		UserID:   userID,
		NpcChara: npcChara,
		NpcEquip: npcEquip,

		// 追加：開始時は未確定
		PendingNpcHand: nil,
	}, nil
}

// RoundBattle は1ラウンドの攻撃処理を行い結果を返す
func (s *BattleService) RoundBattle(ctx context.Context, session *BattleSession, playerHand domain.AttackType) (*RoundResult, error) {
	// NPCの手は「先に確定されたもの」を優先して使う
	var npcHand domain.AttackType
	if session.PendingNpcHand != nil {
		npcHand = *session.PendingNpcHand
		session.PendingNpcHand = nil // 消費
	} else {
		// 保険：prepare_round が呼ばれてない場合でも動く
		npcHand = randomAttackType()
	}

	jankenResult := domain.JudgeJanken(playerHand, npcHand)
	if jankenResult == domain.Win {
		session.Battle.PlayerAttack(playerHand)
	} else if jankenResult == domain.Lose {
		session.Battle.OppoAttack(npcHand)
	}

	playerHP := session.Battle.PlayerHP
	npcHP := session.Battle.OppoHP

	// どちらかのHPが0以下になったらゲーム終了
	if playerHP <= 0 || npcHP <= 0 {
		outcome, delta, coinReward, stoneReward, err := s.finishBattle(ctx, session, playerHP, npcHP)
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
			CoinReward:     coinReward,
			StoneReward:    stoneReward,
		}, nil
	}

	return &RoundResult{
		PlayerHP: playerHP,
		NpcHP:    npcHP,
		NpcHand:  npcHand,
		IsOver:   false,
	}, nil
}

// finishBattle は勝敗を判定しランクポイントをDBに保存し、勝利時は報酬を付与する
func (s *BattleService) finishBattle(ctx context.Context, session *BattleSession, playerHP, npcHP int) (string, int, int, int, error) {
	const winDelta = 10
	const loseDelta = 5
	const winCoinReward = 100
	const winStoneReward = 1

	var outcome string
	var delta int
	var coinReward int
	var stoneReward int

	switch {
	case playerHP > 0 && npcHP <= 0:
		outcome = string(domain.Win)
		delta = winDelta
		coinReward = winCoinReward
		stoneReward = winStoneReward
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
		return "", 0, 0, 0, errs.ErrUserNotFound
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
		return "", 0, 0, 0, err
	}

	// 勝利時はコインとガチャ石を付与
	if outcome == string(domain.Win) {
		_, err = s.queries.AddUserCoin(ctx, sqlc.AddUserCoinParams{
			Amount: int32(coinReward),
			UserID: session.UserID,
		})
		if err != nil {
			return "", 0, 0, 0, err
		}

		_, err = s.queries.AddUserGachaStone(ctx, sqlc.AddUserGachaStoneParams{
			Amount: int32(stoneReward),
			UserID: session.UserID,
		})
		if err != nil {
			return "", 0, 0, 0, err
		}
	}

	return outcome, delta, coinReward, stoneReward, nil
}

// LoadPlayerData はプレイヤーのキャラクターと装備をロードする（オンラインバトル用）
func (s *BattleService) LoadPlayerData(ctx context.Context, charaID, equipID uuid.UUID) (*domain.Character, *domain.Equip, error) {
	dbPlayerChara, err := s.queries.GetCharacterWithCard(ctx, charaID)
	if err != nil {
		return nil, nil, errs.ErrInvalidCharaID
	}
	dbPlayerEquip, err := s.queries.GetEquipmentWithCard(ctx, equipID)
	if err != nil {
		return nil, nil, errs.ErrInvalidEquipID
	}

	playerChara := todomainCharaWithCard(dbPlayerChara)
	playerEquip := todomainEquipWithCard(dbPlayerEquip)

	if playerChara == nil || playerEquip == nil {
		return nil, nil, errs.ErrInvalidCharaID
	}

	return playerChara, playerEquip, nil
}

func todomainCharaWithCard(c sqlc.GetCharacterWithCardRow) *domain.Character {
	rarity := ""
	if c.Rarity.Valid {
		rarity = c.Rarity.String
	}
	iconURL := ""
	if c.CardIconUrl.Valid {
		iconURL = c.CardIconUrl.String
	}
	chara, _ := domain.NewCharacter(
		fmt.Sprintf("%d", c.CardID),
		c.CharacterID.String(),
		c.CardName,
		iconURL,
		domain.Rarity(rarity),
		int(c.InitHp), int(c.InitAtk), int(c.InitTech),
		int(c.MaxHp), int(c.MaxAtk), int(c.MaxTech),
		domain.SpecialType(c.SpecialType.String),
	)
	if chara != nil {
		chara.HP = int(c.Hp)
		chara.ATK = int(c.Atk)
		chara.TECH = int(c.Tech)
	}
	return chara
}

func todomainEquipWithCard(e sqlc.GetEquipmentWithCardRow) *domain.Equip {
	rarity := ""
	if e.Rarity.Valid {
		rarity = e.Rarity.String
	}
	iconURL := ""
	if e.CardIconUrl.Valid {
		iconURL = e.CardIconUrl.String
	}
	equip, _ := domain.NewEquip(
		fmt.Sprintf("%d", e.CardID),
		e.EquipmentID.String(),
		e.CardName,
		iconURL,
		domain.Rarity(rarity),
		int(e.InitBonusHp), int(e.InitBonusAtk), int(e.InitBonusTech),
		int(e.MaxBonusHp), int(e.MaxBonusAtk), int(e.MaxBonusTech),
		nil,
	)
	if equip != nil {
		equip.BonusHP = int(e.BonusHp)
		equip.BonusATK = int(e.BonusAtk)
		equip.BonusTECH = int(e.BonusTech)
	}
	return equip
}

// GetAllCharactersWithCard用のオーバーロード
func todomainCharaWithCardFromMany(c sqlc.GetAllCharactersWithCardRow) *domain.Character {
	rarity := ""
	if c.Rarity.Valid {
		rarity = c.Rarity.String
	}
	iconURL := ""
	if c.CardIconUrl.Valid {
		iconURL = c.CardIconUrl.String
	}
	chara, _ := domain.NewCharacter(
		fmt.Sprintf("%d", c.CardID),
		c.CharacterID.String(),
		c.CardName,
		iconURL,
		domain.Rarity(rarity),
		int(c.InitHp), int(c.InitAtk), int(c.InitTech),
		int(c.MaxHp), int(c.MaxAtk), int(c.MaxTech),
		domain.SpecialType(c.SpecialType.String),
	)
	if chara != nil {
		chara.HP = int(c.Hp)
		chara.ATK = int(c.Atk)
		chara.TECH = int(c.Tech)
	}
	return chara
}

func todomainEquipWithCardFromMany(e sqlc.GetAllEquipmentsWithCardRow) *domain.Equip {
	rarity := ""
	if e.Rarity.Valid {
		rarity = e.Rarity.String
	}
	iconURL := ""
	if e.CardIconUrl.Valid {
		iconURL = e.CardIconUrl.String
	}
	equip, _ := domain.NewEquip(
		fmt.Sprintf("%d", e.CardID),
		e.EquipmentID.String(),
		e.CardName,
		iconURL,
		domain.Rarity(rarity),
		int(e.InitBonusHp), int(e.InitBonusAtk), int(e.InitBonusTech),
		int(e.MaxBonusHp), int(e.MaxBonusAtk), int(e.MaxBonusTech),
		nil,
	)
	if equip != nil {
		equip.BonusHP = int(e.BonusHp)
		equip.BonusATK = int(e.BonusAtk)
		equip.BonusTECH = int(e.BonusTech)
	}
	return equip
}
