package rakuten

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	categoryListEndpoint    = "https://openapi.rakuten.co.jp/recipems/api/Recipe/CategoryList/20170426"
	categoryRankingEndpoint = "https://openapi.rakuten.co.jp/recipems/api/Recipe/CategoryRanking/20170426"
)

// Category は楽天レシピカテゴリ情報です。
type Category struct {
	ID   string
	Name string
	Type string
}

// RankingRecipe は楽天レシピランキング情報です。
type RankingRecipe struct {
	ID              string
	Title           string
	ImageURL        *string
	RecipeURL       string
	CookingTime     *string
	Cost            *string
	Description     string
	Rank            int
	MatchedCategory string
}

// Client は楽天レシピAPIクライアントです。
type Client struct {
	httpClient  *http.Client
	appURL      string
	appID       string
	accessKey   string
	affiliateID string
}

// NewClientFromEnv は環境変数から楽天レシピAPIクライアントを生成します。
func NewClientFromEnv() *Client {
	return &Client{
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		appURL:      strings.TrimSpace(os.Getenv("RAKUTEN_APP_URL")),
		appID:       strings.TrimSpace(os.Getenv("RAKUTEN_APP_ID")),
		accessKey:   strings.TrimSpace(os.Getenv("RAKUTEN_APP_ACCESS_KEY")),
		affiliateID: strings.TrimSpace(os.Getenv("RAKUTEN_APP_AFILIATE_ID")),
	}
}

// Enabled は楽天レシピAPIの設定有無を返します。
func (c *Client) Enabled() bool {
	return c.appID != "" && c.accessKey != ""
}

// ListCategories は楽天レシピカテゴリ一覧を返します。
func (c *Client) ListCategories(ctx context.Context) ([]Category, error) {
	if !c.Enabled() {
		return nil, fmt.Errorf("楽天レシピAPI設定が不足しています")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, categoryListEndpoint, nil)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("applicationId", c.appID)
	query.Set("accessKey", c.accessKey)
	query.Set("format", "json")
	query.Set("formatVersion", "2")
	if c.affiliateID != "" {
		query.Set("affiliateId", c.affiliateID)
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Authorization", "Bearer "+c.accessKey)
	log.Printf("rakuten: カテゴリ一覧リクエストを送信します endpoint=%q", categoryListEndpoint)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("楽天レシピカテゴリ一覧の取得に失敗しました: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		snippet := readErrorBodySnippet(res)
		log.Printf("rakuten: カテゴリ一覧取得失敗 status=%d body=%q", res.StatusCode, snippet)
		return nil, fmt.Errorf("楽天レシピカテゴリ一覧取得エラー: status=%d", res.StatusCode)
	}

	var payload struct {
		Result struct {
			Large []struct {
				CategoryID   string `json:"categoryId"`
				CategoryName string `json:"categoryName"`
			} `json:"large"`
			Medium []struct {
				CategoryID   string `json:"categoryId"`
				CategoryName string `json:"categoryName"`
			} `json:"medium"`
			Small []struct {
				CategoryID   string `json:"categoryId"`
				CategoryName string `json:"categoryName"`
			} `json:"small"`
		} `json:"result"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	result := make([]Category, 0, len(payload.Result.Large)+len(payload.Result.Medium)+len(payload.Result.Small))
	for _, item := range payload.Result.Large {
		result = append(result, Category{ID: item.CategoryID, Name: item.CategoryName, Type: "large"})
	}
	for _, item := range payload.Result.Medium {
		result = append(result, Category{ID: item.CategoryID, Name: item.CategoryName, Type: "medium"})
	}
	for _, item := range payload.Result.Small {
		result = append(result, Category{ID: item.CategoryID, Name: item.CategoryName, Type: "small"})
	}

	return result, nil
}

// GetCategoryRanking はカテゴリ別ランキングを返します。
func (c *Client) GetCategoryRanking(ctx context.Context, categoryID string, categoryName string) ([]RankingRecipe, error) {
	if !c.Enabled() {
		return nil, fmt.Errorf("楽天レシピAPI設定が不足しています")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, categoryRankingEndpoint, nil)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("applicationId", c.appID)
	query.Set("accessKey", c.accessKey)
	query.Set("format", "json")
	query.Set("formatVersion", "2")
	query.Set("categoryId", categoryID)
	if c.affiliateID != "" {
		query.Set("affiliateId", c.affiliateID)
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Authorization", "Bearer "+c.accessKey)
	log.Printf("rakuten: ランキングリクエストを送信します category_id=%q category_name=%q", categoryID, categoryName)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("楽天レシピランキング取得に失敗しました: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		snippet := readErrorBodySnippet(res)
		log.Printf("rakuten: ランキング取得失敗 status=%d category_id=%q category_name=%q body=%q", res.StatusCode, categoryID, categoryName, snippet)
		return nil, fmt.Errorf("楽天レシピランキング取得エラー: status=%d", res.StatusCode)
	}

	var payload struct {
		Result []struct {
			RecipeID          int    `json:"recipeId"`
			RecipeTitle       string `json:"recipeTitle"`
			RecipeURL         string `json:"recipeUrl"`
			FoodImageURL      string `json:"foodImageUrl"`
			MediumImageURL    string `json:"mediumImageUrl"`
			RecipeDescription string `json:"recipeDescription"`
			RecipeIndication  string `json:"recipeIndication"`
			RecipeCost        string `json:"recipeCost"`
			Rank              int    `json:"rank"`
		} `json:"result"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	result := make([]RankingRecipe, 0, len(payload.Result))
	for _, item := range payload.Result {
		imageURL := strings.TrimSpace(item.FoodImageURL)
		if imageURL == "" {
			imageURL = strings.TrimSpace(item.MediumImageURL)
		}

		var imagePtr *string
		if imageURL != "" {
			imagePtr = &imageURL
		}

		var cookingTime *string
		if strings.TrimSpace(item.RecipeIndication) != "" && item.RecipeIndication != "指定なし" {
			value := item.RecipeIndication
			cookingTime = &value
		}

		var cost *string
		if strings.TrimSpace(item.RecipeCost) != "" && item.RecipeCost != "指定なし" {
			value := item.RecipeCost
			cost = &value
		}

		result = append(result, RankingRecipe{
			ID:              fmt.Sprintf("%d", item.RecipeID),
			Title:           item.RecipeTitle,
			ImageURL:        imagePtr,
			RecipeURL:       item.RecipeURL,
			CookingTime:     cookingTime,
			Cost:            cost,
			Description:     strings.TrimSpace(item.RecipeDescription),
			Rank:            item.Rank,
			MatchedCategory: categoryName,
		})
	}

	return result, nil
}

func readErrorBodySnippet(res *http.Response) string {
	body, err := io.ReadAll(io.LimitReader(res.Body, 512))
	if err != nil {
		return "response body read failed"
	}

	return strings.TrimSpace(string(body))
}
