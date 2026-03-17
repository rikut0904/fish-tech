package router

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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
	if shouldRunAutoMigrate() {
		if err := gorm.AutoMigrateAll(db); err != nil {
			return nil, fmt.Errorf("DBマイグレーションに失敗しました: %w", err)
		}
	}

	// ユースケースの初期化
	helloUseCase := hello.NewHelloUseCase()
	adminRepository, err := repository.NewAdminRepository(db)
	if err != nil {
		return nil, fmt.Errorf("管理画面用テーブル初期化に失敗しました: %w", err)
	}
	photosClient := googlephotos.NewClientFromEnv()
	adminUseCase := admin.NewAdminUseCaseWithResolver(adminRepository, photosClient)
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
	adminHTTPHandler := adminHandler.NewAdminHandler(adminUseCase, photosClient)
	allowedAdminOrigins := parseAllowedOrigins(os.Getenv("ADMIN_ALLOWED_ORIGINS"), defaultAdminOrigin)
	allowedSwaggerOrigins := parseAllowedOrigins(
		os.Getenv("SWAGGER_ALLOWED_ORIGINS"),
		"http://localhost:3000",
	)
	allowedSwaggerRefererPaths := parseAllowedRefererPaths(
		os.Getenv("SWAGGER_ALLOWED_REFERER_PATHS"),
		[]string{"/swagger", "/swagger_ui.html"},
	)

	// ルーティング
	api := e.Group("/api")
	{
		api.GET("/hello", helloHandler.GetHello)
		api.GET("/fishes", publicFishHandler.ListFishes)
		api.GET("/pairs", publicFishHandler.ListPairs)
		api.GET("/places/recommendations", placeHandler.GetRecommendedPlaces)
		api.PATCH("/places/favorite", placeHandler.UpdatePlaceFavorite)

		adminGroup := api.Group("/admin")
		adminGroup.Use(RequireAdminOrigin(allowedAdminOrigins, allowedSwaggerOrigins, allowedSwaggerRefererPaths))
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

func parseAllowedRefererPaths(raw string, defaults []string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return defaults
	}

	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		normalized := normalizeRefererPath(part)
		if normalized == "" {
			continue
		}
		result = append(result, normalized)
	}

	if len(result) == 0 {
		return defaults
	}

	return result
}

// shouldRunAutoMigrate は起動時マイグレーションの有効状態を返します。
func shouldRunAutoMigrate() bool {
	enabled, err := strconv.ParseBool(strings.TrimSpace(os.Getenv("AUTO_MIGRATE")))
	if err != nil {
		return false
	}

	return enabled
}
