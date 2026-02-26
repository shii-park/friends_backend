package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/shii-park/friends/internal/sqlc"
)

// Setup はデータベース接続をセットアップします
func Setup(dbDriver string, host string, port int, user string, password string, dbname string, sslmode string) (*sql.DB, error) {
	dsn := "host=" + host + " port=" + strconv.Itoa(port) + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=" + sslmode
	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, err
}

// InitSchema は schema.sql を実行してテーブルを作成します
func InitSchema(db *sql.DB) error {
	schema, err := os.ReadFile("internal/sqlc/schema/schema.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		// すでにテーブルがある場合のエラー（relation already exists）は無視するようにしたいが、
		// 簡易化のため一旦エラーを返さずログに出す運用にする
		log.Printf("Schema info: %v", err)
	}
	return nil
}

type CardData struct {
	CardName      string `json:"card_name"`
	CardKind      int16  `json:"card_kind"`
	CardIconURL   string `json:"card_icon_url"`
	Rarity        string `json:"rarity"`
	CardDetail    string `json:"card_detail"`
	HP            int32  `json:"hp"`
	ATK           int32  `json:"atk"`
	Tech          int32  `json:"tech"`
	InitHP        int32  `json:"init_hp"`
	InitATK       int32  `json:"init_atk"`
	InitTech      int32  `json:"init_tech"`
	MaxHP         int32  `json:"max_hp"`
	MaxATK        int32  `json:"max_atk"`
	MaxTECH       int32  `json:"max_tech"`
	SpecialType   string `json:"special_type"`
	BonusHP       int32  `json:"bonus_hp"`
	BonusATK      int32  `json:"bonus_atk"`
	BonusTech     int32  `json:"bonus_tech"`
	InitBonusHP   int32  `json:"init_bonus_hp"`
	InitBonusATK  int32  `json:"init_bonus_atk"`
	InitBonusTech int32  `json:"init_bonus_tech"`
	MaxBonusHP    int32  `json:"max_bonus_hp"`
	MaxBonusATK   int32  `json:"max_bonus_atk"`
	MaxBonusTECH  int32  `json:"max_bonus_tech"`
	BuffEffect    string `json:"buff_effect"`
}

// Seed はカードマスタデータを投入します
func Seed(db *sql.DB, queries *sqlc.Queries) error {
	ctx := context.Background()

	// 既にデータがあるか確認
	count, err := queries.CountCards(ctx)
	if err != nil {
		// テーブルがない場合はCountCardsが失敗するので、InitSchema後に呼び出す必要がある
		return err
	}
	if count > 0 {
		log.Println("Database already has card data. Skipping seed.")
		return nil
	}

	log.Println("Seeding card data...")
	data, err := os.ReadFile("internal/domain/cards.json")
	if err != nil {
		return err
	}

	var cards []CardData
	if err := json.Unmarshal(data, &cards); err != nil {
		return err
	}

	for _, c := range cards {
		cardID, err := queries.CreateCard(ctx, sqlc.CreateCardParams{
			CardName: c.CardName,
			CardKind: c.CardKind,
			CardIconUrl: sql.NullString{String: c.CardIconURL, Valid: c.CardIconURL != ""},
			Rarity:      sql.NullString{String: c.Rarity, Valid: c.Rarity != ""},
			CardDetail:  sql.NullString{String: c.CardDetail, Valid: c.CardDetail != ""},
		})
		if err != nil {
			log.Printf("failed to create card %s: %v", c.CardName, err)
			continue
		}

		if c.CardKind == 1 { // Character
			err = queries.CreateCharacter(ctx, sqlc.CreateCharacterParams{
				CharacterID: uuid.New(),
				CardID:      cardID,
				Hp:          c.HP,
				Atk:         c.ATK,
				Tech:        c.Tech,
				InitHp:      c.InitHP,
				InitAtk:     c.InitATK,
				InitTech:    c.InitTech,
				MaxHp:       c.MaxHP,
				MaxAtk:      c.MaxATK,
				MaxTech:     c.MaxTECH,
				SpecialType: sql.NullString{String: c.SpecialType, Valid: c.SpecialType != ""},
			})
		} else { // Equipment
			err = queries.CreateEquipment(ctx, sqlc.CreateEquipmentParams{
				EquipmentID:   uuid.New(),
				CardID:        cardID,
				BonusHp:       c.BonusHP,
				BonusAtk:      c.BonusATK,
				BonusTech:     c.BonusTech,
				InitBonusHp:   c.InitBonusHP,
				InitBonusAtk:  c.InitBonusATK,
				InitBonusTech: c.InitBonusTech,
				MaxBonusHp:    c.MaxBonusHP,
				MaxBonusAtk:   c.MaxBonusATK,
				MaxBonusTech:  c.MaxBonusTECH,
				BuffEffect:    sql.NullString{String: c.BuffEffect, Valid: c.BuffEffect != ""},
			})
		}
		if err != nil {
			log.Printf("failed to create detail for card %s: %v", c.CardName, err)
		}
	}

	log.Println("Seeding completed.")
	return nil
}
