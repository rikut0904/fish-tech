package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"your_project/internal/domain/fish"
)

type RakutenRecipeClient struct {
	AppID string
}

func NewRakutenRecipeClient() *RakutenRecipeClient {
	return &RakutenRecipeClient{
		AppID: os.Getenv("RAKUTEN_APP_ID"),
	}
}

// 楽天レシピカテゴリ順位APIのレスポンス構造
type RakutenRankingResponse struct {
	Result []struct {
		RecipeId         int    `json:"recipeId"`
		RecipeTitle      string `json:"recipeTitle"`
		RecipeUrl        string `json:"recipeUrl"`
		FoodImageUrl     string `json:"foodImageUrl"`
		RecipeCost       string `json:"recipeCost"`
		RecipeIndication string `json:"recipeIndication"`
	} `json:"result"`
}

func (c *RakutenRecipeClient) FetchRankingByCategoryId(categoryID string) ([]fish.Recipe, error) {
	// 参照記事に基づいたエンドポイント
	url := fmt.Sprintf("https://app.rakuten.co.jp/services/api/Recipe/CategoryRanking/20170426?applicationId=%s&categoryId=%s", c.AppID, categoryID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rakutenResp RakutenRankingResponse
	if err := json.NewDecoder(resp.Body).Decode(&rakutenResp); err != nil {
		return nil, err
	}

	var recipes []fish.Recipe
	for _, r := range rakutenResp.Result {
		recipes = append(recipes, fish.Recipe{
			ID:           fmt.Sprintf("%d", r.RecipeId),
			Title:        r.RecipeTitle,
			RecipeURL:    r.RecipeUrl,
			ImageURL:     r.FoodImageUrl,
			CookingTime:  r.RecipeIndication,
			Cost:         r.RecipeCost,
		})
	}
	return recipes, nil
}