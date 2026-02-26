package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/domain"
	"github.com/shii-park/friends/internal/service"
)

type OnlineBattleGinHandler struct {
	svc       *service.BattleService
	matchmaker *service.Matchmaker
}

func NewOnlineBattleGinHandler(svc *service.BattleService, matchmaker *service.Matchmaker) *OnlineBattleGinHandler {
	return &OnlineBattleGinHandler{svc: svc, matchmaker: matchmaker}
}

// GET /battle/online/ws
func (h *OnlineBattleGinHandler) WS(c *gin.Context) {
	userID, err := mustUserIDFromSession(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// "start" メッセージを待つ（30秒タイムアウト）
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	var startMsg inMessage
	if err := conn.ReadJSON(&startMsg); err != nil || startMsg.Type != "start" {
		conn.WriteJSON(outMessage{Type: "error", Error: "startメッセージが必要です"})
		return
	}
	conn.SetReadDeadline(time.Time{}) // タイムアウトリセット

	charaID, err := uuid.Parse(startMsg.CharaID)
	if err != nil {
		conn.WriteJSON(outMessage{Type: "error", Error: "無効なキャラクターIDです"})
		return
	}
	equipID, err := uuid.Parse(startMsg.EquipID)
	if err != nil {
		conn.WriteJSON(outMessage{Type: "error", Error: "無効な装備IDです"})
		return
	}

	// プレイヤーデータをロード
	chara, equip, err := h.svc.LoadPlayerData(c.Request.Context(), userID, charaID, equipID)
	if err != nil {
		conn.WriteJSON(outMessage{Type: "error", Error: err.Error()})
		return
	}

	player := service.NewWaitingPlayer(conn, userID, chara, equip)

	// マッチングキューに追加
	opponent := h.matchmaker.Enqueue(player)

	var room *service.OnlineBattleRoom
	isA := false

	// 切断時に相手へ通知する defer（room が nil や正常終了の場合はスキップ）
	defer func() {
		if room == nil || room.IsEnded() {
			return
		}
		var opponentWriter *service.WaitingPlayer
		if isA {
			opponentWriter = room.WriterB
		} else {
			opponentWriter = room.WriterA
		}
		opponentWriter.WriteJSON(outMessage{Type: "opponent_disconnected"}) //nolint:errcheck
	}()

	if opponent == nil {
		// 待機中: マッチング待ちを通知
		conn.WriteJSON(outMessage{Type: "matching"}) //nolint:errcheck

		// ルームが届くのを 60 秒待つ
		select {
		case room = <-player.RoomCh:
			// マッチ成立
		case <-time.After(60 * time.Second):
			h.matchmaker.Dequeue(player)
			conn.WriteJSON(outMessage{Type: "error", Error: "マッチングタイムアウトです"}) //nolint:errcheck
			return
		}
		isA = room.WriterA.UserID == userID
	} else {
		// 即マッチ: ルームを作成して待機中プレイヤーへ通知
		room = h.matchmaker.CreateRoom(opponent, player)
		isA = room.WriterA.UserID == userID
		opponent.RoomCh <- room
	}

	// ready メッセージを送信（自分視点）
	battle := room.Battle()
	var readyMsg outMessage
	if isA {
		readyMsg = outMessage{
			Type:            "ready",
			PlayerHP:        battle.PlayerHP,
			NpcHP:           battle.OppoHP,
			NpcCharaName:    room.WriterB.Chara.Name,
			NpcCharaRarity:  string(room.WriterB.Chara.Rarity),
			NpcCharaIconURL: room.WriterB.Chara.Icon,
			NpcSpecialType:  string(room.WriterB.Chara.SpecialType),
			NpcEquipName:    room.WriterB.Equip.Name,
			NpcEquipRarity:  string(room.WriterB.Equip.Rarity),
			NpcEquipIconURL: room.WriterB.Equip.Icon,
		}
	} else {
		readyMsg = outMessage{
			Type:            "ready",
			PlayerHP:        battle.OppoHP,
			NpcHP:           battle.PlayerHP,
			NpcCharaName:    room.WriterA.Chara.Name,
			NpcCharaRarity:  string(room.WriterA.Chara.Rarity),
			NpcCharaIconURL: room.WriterA.Chara.Icon,
			NpcSpecialType:  string(room.WriterA.Chara.SpecialType),
			NpcEquipName:    room.WriterA.Equip.Name,
			NpcEquipRarity:  string(room.WriterA.Equip.Rarity),
			NpcEquipIconURL: room.WriterA.Equip.Icon,
		}
	}
	conn.WriteJSON(readyMsg) //nolint:errcheck

	// バトルループ
	for {
		var msg inMessage
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}

		if msg.Type != "round" {
			conn.WriteJSON(outMessage{Type: "error", Error: "roundメッセージが必要です"}) //nolint:errcheck
			continue
		}

		hand := domain.AttackType(msg.Hand)
		if hand != domain.Rock && hand != domain.Paper && hand != domain.Scissors {
			conn.WriteJSON(outMessage{Type: "error", Error: "無効な手です"}) //nolint:errcheck
			continue
		}

		result, err := room.SubmitHand(c.Request.Context(), isA, hand)
		if err != nil {
			conn.WriteJSON(outMessage{Type: "error", Error: err.Error()}) //nolint:errcheck
			return
		}

		if result == nil {
			// まだ相手が手を出していない
			continue
		}

		// 両者が手を出した: 結果を両プレイヤーへ送信
		var msgForA, msgForB outMessage
		if result.IsOver {
			msgForA = outMessage{
				Type:           "game_over",
				PlayerHP:       result.ForA.PlayerHP,
				NpcHP:          result.ForA.NpcHP,
				NpcHand:        string(result.ForA.NpcHand),
				Outcome:        result.ForA.Outcome,
				RankPointDelta: result.ForA.RankPointDelta,
				CoinReward:     result.ForA.CoinReward,
				StoneReward:    result.ForA.StoneReward,
			}
			msgForB = outMessage{
				Type:           "game_over",
				PlayerHP:       result.ForB.PlayerHP,
				NpcHP:          result.ForB.NpcHP,
				NpcHand:        string(result.ForB.NpcHand),
				Outcome:        result.ForB.Outcome,
				RankPointDelta: result.ForB.RankPointDelta,
				CoinReward:     result.ForB.CoinReward,
				StoneReward:    result.ForB.StoneReward,
			}
		} else {
			msgForA = outMessage{
				Type:     "round_result",
				PlayerHP: result.ForA.PlayerHP,
				NpcHP:    result.ForA.NpcHP,
				NpcHand:  string(result.ForA.NpcHand),
			}
			msgForB = outMessage{
				Type:     "round_result",
				PlayerHP: result.ForB.PlayerHP,
				NpcHP:    result.ForB.NpcHP,
				NpcHand:  string(result.ForB.NpcHand),
			}
		}

		room.WriterA.WriteJSON(msgForA) //nolint:errcheck
		room.WriterB.WriteJSON(msgForB) //nolint:errcheck

		if result.IsOver {
			return
		}
	}
}
