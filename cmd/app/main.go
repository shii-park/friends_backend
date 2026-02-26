package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shii-park/friends/internal/db"
	"github.com/shii-park/friends/internal/handler"
	"github.com/shii-park/friends/internal/middleware"
	"github.com/shii-park/friends/internal/repository"
	"github.com/shii-park/friends/internal/service"
	"github.com/shii-park/friends/internal/sqlc"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// DBのセットアップ
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
		log.Fatalf("データベースのセットアップに失敗しました: %v", err)
	}
	log.Println("データベースのセットアップが正常に完了しました")
	defer dbConn.Close()

	// データベースの初期化
	if err := db.InitSchema(dbConn); err != nil {
		log.Printf("スキーマの初期化に失敗しました: %v", err)
	}

	// sqlcセットアップ
	queries := sqlc.New(dbConn)

	// シードデータの投入
	if err := db.Seed(dbConn, queries); err != nil {
		log.Printf("シードデータの投入に失敗しました: %v", err)
	}

	// serviceセットアップ
	registerSvc := service.NewRegisterService(queries)
	// Storage
	storageRepo := repository.NewStorageRepository(queries)
	storageSvc := service.NewStorageService(storageRepo)
	storageHandler := handler.NewStorageGinHandler(storageSvc)
	battleSvc := service.NewBattleService(queries)
	battleHandler := handler.NewBattleGinHandler(battleSvc)

	// Enhancement
	enhRepo := repository.NewEnhancementRepository(dbConn, queries)
	enhSvc := service.NewEnhancementService(enhRepo)
	enhHandler := handler.NewEnhancementGinHandler(enhSvc)

	// User
	getUserSvc := service.NewGetUserService(queries)
	userRepo := repository.NewUserRepository(queries)

	coinSvc := service.NewCoinService(userRepo)
	coinHandler := handler.NewCoinGinHandler(coinSvc)

	gachaStoneSvc := service.NewGachaStoneService(userRepo)
	gachaStoneHandler := handler.NewGachaStoneGinHandler(gachaStoneSvc)

	gachaRepo := repository.NewGachaRepository(dbConn, queries, storageRepo)
	gachaSvc := service.NewGachaService(dbConn, queries, gachaRepo, storageRepo)
	gachaHandler := handler.NewGachaGinHandler(gachaSvc)

	collectionRepo := repository.NewCollectionRepository(queries)
	collectionSvc := service.NewCollectionService(collectionRepo)
	collectionHandler := handler.NewCollectionGinHandler(collectionSvc)

	r := gin.Default()
	// CORS設定
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))
	// TODO: クッキーの秘密鍵や名前の変更
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// 認証なしエンドポイント
	// 新規登録
	r.POST("/register", handler.RegisterHandler(registerSvc, storageSvc))

	// 認証が必要なエンドポイント
	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		auth.GET("/battle/ws", battleHandler.WS)
		auth.GET("/gacha/characters", gachaHandler.ListCharacters)
		auth.GET("/gacha/equipments", gachaHandler.ListEquipments)
		auth.GET("/gacha/lineup", gachaHandler.Lineup)
		auth.POST("/gacha/draw", gachaHandler.Draw)

		auth.GET("/storage", storageHandler.ListCards)
		auth.POST("/storage/cards", storageHandler.AddCard)
		auth.DELETE("/storage/cards/:instanceID", storageHandler.RemoveCard)
		auth.GET("/storage/detail", storageHandler.ListCardDetails)

		auth.GET("/user/:userID/get", handler.GetUserHandler(getUserSvc))

		auth.POST("/user/:userID/coin/add", coinHandler.Add)
		auth.POST("/user/:userID/gacha-stone/add", gachaStoneHandler.Add)

		auth.POST("/card/:instanceID/enhancement", enhHandler.Enhance)
		auth.GET("/card/:instanceID/detail", storageHandler.Detail)
		auth.GET("/collection", collectionHandler.List)
		// auth.POST("/matching", testHandler)
	}

	// ポート8080番でリッスン
	r.Run()
}

// TODO: すべてのハンドラーが完成したら以下を削除
// func testHandler(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"message": "ok"})
// }
