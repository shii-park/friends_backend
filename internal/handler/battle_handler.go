package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/shii-park/friends/internal/domain"
	"github.com/shii-park/friends/internal/service"
)

var upgrader = websocket.Upgrader{
	// TODO: CORS設定
	CheckOrigin: func(r *http.Request) bool { return true },
}

type inMessage struct {
	Type    string `json:"type"`    // "start" | "round"
	CharaID string `json:"charaID"` // start時
	EquipID string `json:"equipID"` // start時
	Hand    string `json:"hand"`    // round時: "rock"|"paper"|"scissors"
}

type outMessage struct {
	Type           string `json:"type"`
	PlayerHP       int    `json:"playerHP,omitempty"`
	NpcHP          int    `json:"npcHP,omitempty"`
	NpcHand        string `json:"npcHand,omitempty"`
	Outcome        string `json:"outcome,omitempty"`
	RankPointDelta int    `json:"rankPointDelta,omitempty"`
	Error          string `json:"error,omitempty"`
}

type BattleGinHandler struct {
	svc *service.BattleService
}

func NewBattleGinHandler(svc *service.BattleService) *BattleGinHandler {
	return &BattleGinHandler{svc: svc}
}

// GET /battle/ws
func (h *BattleGinHandler) WS(c *gin.Context) {
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

	var session *service.BattleSession

	for {
		var msg inMessage
		if err := conn.ReadJSON(&msg); err != nil {
			// クライアント切断
			break
		}

		switch msg.Type {
		case "start":
			charaID, err := uuid.Parse(msg.CharaID)
			if err != nil {
				conn.WriteJSON(outMessage{Type: "error", Error: "無効なキャラクターIDです"})
				return
			}
			equipID, err := uuid.Parse(msg.EquipID)
			if err != nil {
				conn.WriteJSON(outMessage{Type: "error", Error: "無効な装備IDです"})
				return
			}

			session, err = h.svc.StartBattle(c.Request.Context(), userID, charaID, equipID)
			if err != nil {
				conn.WriteJSON(outMessage{Type: "error", Error: err.Error()})
				return
			}

			conn.WriteJSON(outMessage{
				Type:     "ready",
				PlayerHP: session.Battle.PlayerHP,
				NpcHP:    session.Battle.OppoHP,
			})

		case "round":
			if session == nil {
				conn.WriteJSON(outMessage{Type: "error", Error: "バトルが開始されていません"})
				continue
			}

			hand := domain.AttackType(msg.Hand)
			if hand != domain.Rock && hand != domain.Paper && hand != domain.Scissors {
				conn.WriteJSON(outMessage{Type: "error", Error: "無効な手です"})
				continue
			}

			result, err := h.svc.RoundBattle(c.Request.Context(), session, hand)
			if err != nil {
				conn.WriteJSON(outMessage{Type: "error", Error: err.Error()})
				return
			}

			if result.IsOver {
				conn.WriteJSON(outMessage{
					Type:           "game_over",
					PlayerHP:       result.PlayerHP,
					NpcHP:          result.NpcHP,
					NpcHand:        string(result.NpcHand),
					Outcome:        result.Outcome,
					RankPointDelta: result.RankPointDelta,
				})
				return
			}

			conn.WriteJSON(outMessage{
				Type:     "round_result",
				PlayerHP: result.PlayerHP,
				NpcHP:    result.NpcHP,
				NpcHand:  string(result.NpcHand),
			})

		default:
			conn.WriteJSON(outMessage{Type: "error", Error: "不明なメッセージタイプです"})
		}
	}
}
