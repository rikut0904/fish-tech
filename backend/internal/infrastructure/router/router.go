package router

import (
	"fmt"
	"os"

	"fish-tech/internal/infrastructure/googlephotos"
	"fish-tech/internal/infrastructure/hotpepper"
	"fish-tech/internal/infrastructure/persistence/gorm"
	"fish-tech/internal/infrastructure/persistence/gorm/repository"
	"fish-tech/internal/interface/handler"
	adminHandler "fish-tech/internal/interface/handler/admin"
	"fish-tech/internal/usecase/admin"
	"fish-tech/internal/usecase/hello"
	placeUseCasePkg "fish-tech/internal/usecase/place"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	// defaultAdminOrigin は管理画面オリジンのデフォルト値です。
	defaultAdminOrigin = "http://localhost:3001"
)

// NewRouter は全ルートが設定された新しいEchoルーターを作成します。
func NewRouter() (*echo.Echo, error) {
	e := echo.New()

	// ミドルウェア
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	db, err := gorm.NewPostgresDB()
	if err != nil {
		return nil, fmt.Errorf("DB初期化に失敗しました: %w", err)
	}
	if err := gorm.AutoMigrateAll(db); err != nil {
		return nil, fmt.Errorf("DBマイグレーションに失敗しました: %w", err)
	}

	// ユースケースの初期化
	helloUseCase := hello.NewHelloUseCase()
	adminRepository, err := repository.NewAdminRepository(db)
	if err != nil {
		return nil, fmt.Errorf("管理画面用テーブル初期化に失敗しました: %w", err)
	}
	adminUseCase := admin.NewAdminUseCase(adminRepository)
	hotpepperClient := hotpepper.NewClientFromEnv()
	placeRepository, err := repository.NewPlaceRepository(db, hotpepperClient)
	if err != nil {
		return nil, fmt.Errorf("店舗キャッシュテーブル初期化に失敗しました: %w", err)
	}
	placeUseCase := placeUseCasePkg.NewPlaceUseCase(placeRepository)

	// ハンドラーの初期化
	helloHandler := handler.NewHelloHandler(helloUseCase)
	publicFishHandler := handler.NewPublicFishHandler(adminUseCase)
	placeHandler := handler.NewPlaceHandler(placeUseCase)
	photosClient := googlephotos.NewClientFromEnv()
	adminHTTPHandler := adminHandler.NewAdminHandler(adminUseCase, photosClient)
	allowedAdminOrigins := parseAllowedOrigins(os.Getenv("ADMIN_ALLOWED_ORIGINS"), defaultAdminOrigin)

	// ルーティング
	api := e.Group("/api")
	{
		api.GET("/hello", helloHandler.GetHello)
		api.GET("/fishes", publicFishHandler.ListFishes)
		api.GET("/pairs", publicFishHandler.ListPairs)
		api.GET("/places/recommendations", placeHandler.GetRecommendedPlaces)
		api.PATCH("/places/favorite", placeHandler.UpdatePlaceFavorite)

		adminGroup := api.Group("/admin")
		adminGroup.Use(RequireAdminOrigin(allowedAdminOrigins))
		{
			adminGroup.POST("/fishes/upload-image", adminHTTPHandler.UploadFishImage)
			adminGroup.POST("/fishes", adminHTTPHandler.CreateFish)
			adminGroup.DELETE("/fishes/:id", adminHTTPHandler.DeleteFish)
			adminGroup.POST("/pairs", adminHTTPHandler.CreatePair)
			adminGroup.DELETE("/pairs/:id", adminHTTPHandler.DeletePair)
		}
	}

	return e, nil
}
