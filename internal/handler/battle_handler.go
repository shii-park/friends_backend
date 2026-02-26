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
	// ready時のNPC情報
	NpcCharaName    string `json:"npcCharaName,omitempty"`
	NpcCharaRarity  string `json:"npcCharaRarity,omitempty"`
	NpcCharaIconURL string `json:"npcCharaIconURL,omitempty"`
	NpcEquipName    string `json:"npcEquipName,omitempty"`
	NpcEquipRarity  string `json:"npcEquipRarity,omitempty"`
	NpcEquipIconURL string `json:"npcEquipIconURL,omitempty"`
	NpcSpecialType  string `json:"npcSpecialType,omitempty"`
	// ready時の詳細ステータス
	NpcCharaHP      int `json:"npcCharaHP"`
	NpcCharaATK     int `json:"npcCharaATK"`
	NpcCharaTECH    int `json:"npcCharaTECH"`
	NpcEquipHP      int `json:"npcEquipHP"`
	NpcEquipATK     int `json:"npcEquipATK"`
	NpcEquipTECH    int `json:"npcEquipTECH"`
	PlayerCharaHP   int `json:"playerCharaHP"`
	PlayerCharaATK  int `json:"playerCharaATK"`
	PlayerCharaTECH int `json:"playerCharaTECH"`
	PlayerEquipHP   int `json:"playerEquipHP"`
	PlayerEquipATK  int `json:"playerEquipATK"`
	PlayerEquipTECH int `json:"playerEquipTECH"`
	// game_over時の報酬
	CoinReward  int `json:"coinReward,omitempty"`
	StoneReward int `json:"stoneReward,omitempty"`
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

			resp := outMessage{
				Type:     "ready",
				PlayerHP: session.Battle.PlayerHP,
				NpcHP:    session.Battle.OppoHP,
			}
			if session.Battle.PlayerCard != nil {
				resp.PlayerCharaHP = session.Battle.PlayerCard.HP
				resp.PlayerCharaATK = session.Battle.PlayerCard.ATK
				resp.PlayerCharaTECH = session.Battle.PlayerCard.TECH
			}
			if session.Battle.PlayerEquip != nil {
				resp.PlayerEquipHP = session.Battle.PlayerEquip.BonusHP
				resp.PlayerEquipATK = session.Battle.PlayerEquip.BonusATK
				resp.PlayerEquipTECH = session.Battle.PlayerEquip.BonusTECH
			}
			if session.NpcChara != nil {
				resp.NpcCharaName = session.NpcChara.Name
				resp.NpcCharaRarity = string(session.NpcChara.Rarity)
				resp.NpcCharaIconURL = session.NpcChara.Icon
				resp.NpcSpecialType = string(session.NpcChara.SpecialType)
				resp.NpcCharaHP = session.NpcChara.HP
				resp.NpcCharaATK = session.NpcChara.ATK
				resp.NpcCharaTECH = session.NpcChara.TECH
			}
			if session.NpcEquip != nil {
				resp.NpcEquipName = session.NpcEquip.Name
				resp.NpcEquipRarity = string(session.NpcEquip.Rarity)
				resp.NpcEquipIconURL = session.NpcEquip.Icon
				resp.NpcEquipHP = session.NpcEquip.BonusHP
				resp.NpcEquipATK = session.NpcEquip.BonusATK
				resp.NpcEquipTECH = session.NpcEquip.BonusTECH
			}
			conn.WriteJSON(resp)

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
					CoinReward:     result.CoinReward,
					StoneReward:    result.StoneReward,
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
