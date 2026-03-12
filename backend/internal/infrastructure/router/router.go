package router

import (
	"context"
	"fmt"
	"os"

	"fish-tech/internal/infrastructure/googlephotos"
	"fish-tech/internal/infrastructure/persistence/gorm"
	"fish-tech/internal/infrastructure/persistence/gorm/repository"
	"fish-tech/internal/infrastructure/rakuten"
	"fish-tech/internal/interface/handler"
	adminHandler "fish-tech/internal/interface/handler/admin"
	"fish-tech/internal/usecase/admin"
	"fish-tech/internal/usecase/hello"
	recipeUseCase "fish-tech/internal/usecase/recipe"

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

	// ユースケースの初期化
	helloUseCase := hello.NewHelloUseCase()
	adminRepository, err := repository.NewAdminRepository(db)
	if err != nil {
		return nil, fmt.Errorf("管理画面用テーブル初期化に失敗しました: %w", err)
	}
	photosClient := googlephotos.NewClientFromEnv()
	adminUseCase := admin.NewAdminUseCaseWithResolver(adminRepository, photosClient)
	rakutenClient := rakuten.NewClientFromEnv()
	recipeRepository, err := repository.NewRecipeRepository(db)
	if err != nil {
		return nil, fmt.Errorf("レシピ用テーブル初期化に失敗しました: %w", err)
	}
	recipeUC := recipeUseCase.NewRecipeUseCase(recipeRepository, newRakutenClientAdapter(rakutenClient))

	// ハンドラーの初期化
	helloHandler := handler.NewHelloHandler(helloUseCase)
	publicFishHandler := handler.NewPublicFishHandler(adminUseCase)
	recipeHandler := handler.NewRecipeHandler(recipeUC)
	adminHTTPHandler := adminHandler.NewAdminHandler(adminUseCase, photosClient)
	allowedAdminOrigins := parseAllowedOrigins(os.Getenv("ADMIN_ALLOWED_ORIGINS"), defaultAdminOrigin)

	// ルーティング
	api := e.Group("/api")
	{
		api.GET("/hello", helloHandler.GetHello)
		api.GET("/fishes", publicFishHandler.ListFishes)
		api.GET("/pairs", publicFishHandler.ListPairs)
		api.GET("/recipes", recipeHandler.SearchRecipes)
		api.GET("/recipes/seasonal", recipeHandler.GetSeasonalRecipes)
		api.PATCH("/recipes/favorite", recipeHandler.UpdateRecipeFavorite)

		adminGroup := api.Group("/admin")
		adminGroup.Use(RequireAdminOrigin(allowedAdminOrigins))
		{
			adminGroup.POST("/fishes/upload-image", adminHTTPHandler.UploadFishImage)
			adminGroup.POST("/fishes", adminHTTPHandler.CreateFish)
			adminGroup.PATCH("/fishes/:id/seasons", adminHTTPHandler.UpdateFishSeasons)
			adminGroup.DELETE("/fishes/:id", adminHTTPHandler.DeleteFish)
			adminGroup.POST("/pairs", adminHTTPHandler.CreatePair)
			adminGroup.DELETE("/pairs/:id", adminHTTPHandler.DeletePair)
		}
	}

	return e, nil
}

type rakutenClientAdapter struct {
	client *rakuten.Client
}

func newRakutenClientAdapter(client *rakuten.Client) recipeUseCase.RakutenClient {
	if client == nil {
		return nil
	}

	return &rakutenClientAdapter{client: client}
}

func (a *rakutenClientAdapter) Enabled() bool {
	return a.client != nil && a.client.Enabled()
}

func (a *rakutenClientAdapter) ListCategories(ctx context.Context) ([]recipeUseCase.RakutenCategory, error) {
	categories, err := a.client.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]recipeUseCase.RakutenCategory, 0, len(categories))
	for _, category := range categories {
		result = append(result, recipeUseCase.RakutenCategory{
			ID:   category.ID,
			Name: category.Name,
			Type: category.Type,
		})
	}

	return result, nil
}

func (a *rakutenClientAdapter) GetCategoryRanking(ctx context.Context, categoryID string, categoryName string) ([]recipeUseCase.RakutenRecipe, error) {
	ranking, err := a.client.GetCategoryRanking(ctx, categoryID, categoryName)
	if err != nil {
		return nil, err
	}

	result := make([]recipeUseCase.RakutenRecipe, 0, len(ranking))
	for _, item := range ranking {
		result = append(result, recipeUseCase.RakutenRecipe{
			ID:              item.ID,
			Title:           item.Title,
			ImageURL:        item.ImageURL,
			RecipeURL:       item.RecipeURL,
			CookingTime:     item.CookingTime,
			Cost:            item.Cost,
			Description:     item.Description,
			Rank:            item.Rank,
			MatchedCategory: item.MatchedCategory,
		})
	}

	return result, nil
}
