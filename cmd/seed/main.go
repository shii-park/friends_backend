package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/shii-park/friends/internal/db"
	"github.com/shii-park/friends/internal/sqlc"
)

type CardData struct {
	CardName      string `json:"card_name"`
	CardKind      int16  `json:"card_kind"`
	CardIconURL   string `json:"card_icon_url"`
	Rarity        string `json:"rarity"`
	CardDetail    string `json:"card_detail"`
	// Character specific
	HP          int32  `json:"hp"`
	ATK         int32  `json:"atk"`
	Tech        int32  `json:"tech"`
	InitHP      int32  `json:"init_hp"`
	InitATK     int32  `json:"init_atk"`
	InitTech    int32  `json:"init_tech"`
	MaxHP       int32  `json:"max_hp"`
	MaxATK      int32  `json:"max_atk"`
	MaxTech     int32  `json:"max_tech"`
	SpecialType string `json:"special_type"`
	// Equipment specific
	BonusHP       int32  `json:"bonus_hp"`
	BonusATK      int32  `json:"bonus_atk"`
	BonusTech     int32  `json:"bonus_tech"`
	InitBonusHP   int32  `json:"init_bonus_hp"`
	InitBonusATK  int32  `json:"init_bonus_atk"`
	InitBonusTech int32  `json:"init_bonus_tech"`
	MaxBonusHP    int32  `json:"max_bonus_hp"`
	MaxBonusATK   int32  `json:"max_bonus_atk"`
	MaxBonusTech  int32  `json:"max_bonus_tech"`
	BuffEffect    string `json:"buff_effect"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	dbConn, err := db.Setup(
		os.Getenv("DB_DRIVER"),
		os.Getenv("DB_HOST"),
		port,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("SSLMODE"),
	)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer dbConn.Close()

	queries := sqlc.New(dbConn)
	ctx := context.Background()

	// Read cards.json
	data, err := os.ReadFile("internal/domain/cards.json")
	if err != nil {
		log.Fatalf("failed to read cards.json: %v", err)
	}

	var cards []CardData
	if err := json.Unmarshal(data, &cards); err != nil {
		log.Fatalf("failed to unmarshal cards.json: %v", err)
	}

	for _, c := range cards {
		cardID, err := queries.CreateCard(ctx, sqlc.CreateCardParams{
			CardName: c.CardName,
			CardKind: c.CardKind,
			CardIconUrl: sql.NullString{
				String: c.CardIconURL,
				Valid:  c.CardIconURL != "",
			},
			Rarity: sql.NullString{
				String: c.Rarity,
				Valid:  c.Rarity != "",
			},
			CardDetail: sql.NullString{
				String: c.CardDetail,
				Valid:  c.CardDetail != "",
			},
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
				MaxTech:     c.MaxTech,
				SpecialType: sql.NullString{
					String: c.SpecialType,
					Valid:  c.SpecialType != "",
				},
			})
			if err != nil {
				log.Printf("failed to create character for card %s: %v", c.CardName, err)
			}
		} else if c.CardKind == 0 { // Equipment
			err = queries.CreateEquipment(ctx, sqlc.CreateEquipmentParams{
				EquipmentID: uuid.New(),
				CardID:      cardID,
				BonusHp:     c.BonusHP,
				BonusAtk:    c.BonusATK,
				BonusTech:   c.BonusTech,
				InitBonusHp: c.InitBonusHP,
				InitBonusAtk: c.InitBonusATK,
				InitBonusTech: c.InitBonusTech,
				MaxBonusHp:  c.MaxBonusHP,
				MaxBonusAtk: c.MaxBonusATK,
				MaxBonusTech: c.MaxBonusTech,
				BuffEffect: sql.NullString{
					String: c.BuffEffect,
					Valid:  c.BuffEffect != "",
				},
			})
			if err != nil {
				log.Printf("failed to create equipment for card %s: %v", c.CardName, err)
			}
		}
	}

	fmt.Println("Seeding completed successfully")
}
