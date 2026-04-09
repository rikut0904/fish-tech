package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	recipeDomain "fish-tech/internal/domain/recipe"
	recipeUseCase "fish-tech/internal/usecase/recipe"

	"github.com/labstack/echo/v4"
)

// RecipeHandler はレシピ機能のHTTPハンドラーです。
type RecipeHandler struct {
	useCase recipeUseCase.UseCase
}

type recipeFavoriteRequest struct {
	UserID   string `json:"userId"`
	RecipeID string `json:"recipeId"`
	IsLikes  *bool  `json:"isLikes"`
	Favorite *bool  `json:"favorite"`
}

type recipeFavoriteResponse struct {
	UserID   string `json:"userId"`
	RecipeID string `json:"recipeId"`
	IsLikes  bool   `json:"isLikes"`
}

type recipeSearchResponse struct {
	Page    int                                 `json:"page"`
	PerPage int                                 `json:"perPage"`
	Count   int                                 `json:"count"`
	Total   int                                 `json:"total"`
	Items   []recipeDomain.RecipeRecommendation `json:"items"`
}

const (
	defaultRecipePage  = 1
	defaultRecipeCount = 10
	maxRecipeCount     = 100
)

var allowedRecipeQueries = map[string]struct{}{
	"fishName": {},
	"keyword":  {},
	"userId":   {},
	"favorite": {},
	"count":    {},
	"page":     {},
}

var allowedSeasonalRecipeQueries = map[string]struct{}{
	"fishName": {},
	"userId":   {},
	"favorite": {},
	"count":    {},
	"page":     {},
}

var allowedRecipeFavoriteQueries = map[string]struct{}{
	"userId":   {},
	"recipeId": {},
	"isLikes":  {},
	"favorite": {},
}

// NewRecipeHandler はレシピハンドラーを生成します。
func NewRecipeHandler(useCase recipeUseCase.UseCase) *RecipeHandler {
	return &RecipeHandler{useCase: useCase}
}

// GetSeasonalRecipes は旬の魚レシピ一覧を返します。
func (h *RecipeHandler) GetSeasonalRecipes(c echo.Context) error {
	for key := range c.QueryParams() {
		if _, ok := allowedSeasonalRecipeQueries[key]; !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "未対応のクエリパラメータです: " + key})
		}
	}

	condition, err := parseRecipeSearchCondition(c, false)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	response, err := h.useCase.GetSeasonalRecipes(c.Request().Context(), condition)
	if err != nil {
		switch {
		case errors.Is(err, recipeUseCase.ErrFavoriteUserIDRequired):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, recipeUseCase.ErrFishNotInSeason):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, recipeUseCase.ErrRakutenRecipesUnavailable):
			return c.JSON(http.StatusBadGateway, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "旬の魚レシピ取得に失敗しました"})
		}
	}

	return c.JSON(http.StatusOK, response)
}

// SearchRecipes はレシピ一覧を検索します。
func (h *RecipeHandler) SearchRecipes(c echo.Context) error {
	for key := range c.QueryParams() {
		if _, ok := allowedRecipeQueries[key]; !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "未対応のクエリパラメータです: " + key})
		}
	}

	condition, err := parseRecipeSearchCondition(c, true)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	response, err := h.useCase.SearchRecipes(c.Request().Context(), condition)
	if err != nil {
		if errors.Is(err, recipeUseCase.ErrFavoriteUserIDRequired) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		if errors.Is(err, recipeUseCase.ErrRakutenRecipesUnavailable) {
			return c.JSON(http.StatusBadGateway, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "レシピ検索に失敗しました"})
	}

	return c.JSON(http.StatusOK, recipeSearchResponse{
		Page:    response.Page,
		PerPage: response.PerPage,
		Count:   response.Count,
		Total:   response.Total,
		Items:   response.Items,
	})
}

// UpdateRecipeFavorite はユーザーのレシピお気に入り状態を更新します。
func (h *RecipeHandler) UpdateRecipeFavorite(c echo.Context) error {
	for key := range c.QueryParams() {
		if _, ok := allowedRecipeFavoriteQueries[key]; !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "未対応のクエリパラメータです: " + key})
		}
	}

	request := recipeFavoriteRequest{
		UserID:   strings.TrimSpace(c.QueryParam("userId")),
		RecipeID: strings.TrimSpace(c.QueryParam("recipeId")),
	}

	if rawIsLikes := strings.TrimSpace(c.QueryParam("isLikes")); rawIsLikes != "" {
		value, err := strconv.ParseBool(rawIsLikes)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "isLikesはtrueかfalseで指定してください"})
		}
		request.IsLikes = &value
	}
	if request.IsLikes == nil {
		if rawFavorite := strings.TrimSpace(c.QueryParam("favorite")); rawFavorite != "" {
			value, err := strconv.ParseBool(rawFavorite)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "favoriteはtrueかfalseで指定してください"})
			}
			request.Favorite = &value
		}
	}

	if request.UserID == "" || request.RecipeID == "" || (request.IsLikes == nil && request.Favorite == nil) {
		var body recipeFavoriteRequest
		if err := c.Bind(&body); err == nil {
			if request.UserID == "" {
				request.UserID = strings.TrimSpace(body.UserID)
			}
			if request.RecipeID == "" {
				request.RecipeID = strings.TrimSpace(body.RecipeID)
			}
			if request.IsLikes == nil && body.IsLikes != nil {
				request.IsLikes = body.IsLikes
			}
			if request.Favorite == nil && body.Favorite != nil {
				request.Favorite = body.Favorite
			}
		}
	}

	isLikes := false
	switch {
	case request.IsLikes != nil:
		isLikes = *request.IsLikes
	case request.Favorite != nil:
		isLikes = *request.Favorite
	}

	if err := h.useCase.UpdateRecipeFavorite(c.Request().Context(), request.UserID, request.RecipeID, isLikes); err != nil {
		switch {
		case errors.Is(err, recipeUseCase.ErrInvalidUserID), errors.Is(err, recipeUseCase.ErrInvalidRecipeID):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, recipeUseCase.ErrUserNotFound), errors.Is(err, recipeUseCase.ErrRecipeNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "レシピお気に入りの更新に失敗しました"})
		}
	}

	return c.JSON(http.StatusOK, recipeFavoriteResponse{
		UserID:   request.UserID,
		RecipeID: request.RecipeID,
		IsLikes:  isLikes,
	})
}

func parseRecipeSearchCondition(c echo.Context, allowKeyword bool) (recipeDomain.SearchCondition, error) {
	condition := recipeDomain.SearchCondition{
		FishName: strings.TrimSpace(c.QueryParam("fishName")),
		UserID:   strings.TrimSpace(c.QueryParam("userId")),
		Count:    defaultRecipeCount,
		Page:     defaultRecipePage,
	}
	if allowKeyword {
		condition.Keyword = strings.TrimSpace(c.QueryParam("keyword"))
	}

	if rawCount := strings.TrimSpace(c.QueryParam("count")); rawCount != "" {
		parsedCount, err := strconv.Atoi(rawCount)
		if err != nil {
			return recipeDomain.SearchCondition{}, errors.New("countは数値で指定してください")
		}
		condition.Count = parsedCount
	}
	if rawPage := strings.TrimSpace(c.QueryParam("page")); rawPage != "" {
		parsedPage, err := strconv.Atoi(rawPage)
		if err != nil {
			return recipeDomain.SearchCondition{}, errors.New("pageは数値で指定してください")
		}
		condition.Page = parsedPage
	}
	if condition.Count <= 0 {
		condition.Count = defaultRecipeCount
	}
	if condition.Count > maxRecipeCount {
		condition.Count = maxRecipeCount
	}
	if condition.Page <= 0 {
		condition.Page = defaultRecipePage
	}

	if rawFavorite := strings.TrimSpace(c.QueryParam("favorite")); rawFavorite != "" {
		parsedFavorite, err := strconv.ParseBool(rawFavorite)
		if err != nil {
			return recipeDomain.SearchCondition{}, errors.New("favoriteはtrue/falseで指定してください")
		}
		condition.Favorite = &parsedFavorite
	}

	return condition, nil
}
