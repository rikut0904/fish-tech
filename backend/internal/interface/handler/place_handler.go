package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	placeDomain "fish-tech/internal/domain/place"
	placeUseCase "fish-tech/internal/usecase/place"
)

// PlaceHandler は店舗検索のHTTPハンドラーです。
type PlaceHandler struct {
	useCase placeUseCase.UseCase
}

type placeRecommendationsResponse struct {
	Page    int                            `json:"page"`
	PerPage int                            `json:"perPage"`
	Count   int                            `json:"count"`
	Items   []placeDomain.RecommendedPlace `json:"items"`
}

type updateFavoriteRequest struct {
	UserID   string `json:"userId"`
	PlaceID  string `json:"placeId"`
	Favorite bool   `json:"favorite"`
}

const (
	defaultPage  = 1
	defaultCount = 10
	maxCount     = 100
)

var allowedPlaceQueries = map[string]struct{}{
	"fishName": {},
	"keyword":  {},
	"cityCode": {},
	"userId":   {},
	"favorite": {},
	"count":    {},
	"page":     {},
}

var allowedPlaceFavoriteQueries = map[string]struct{}{
	"userId":   {},
	"placeId":  {},
	"favorite": {},
}

// NewPlaceHandler は新しい店舗検索ハンドラーを作成します。
func NewPlaceHandler(useCase placeUseCase.UseCase) *PlaceHandler {
	return &PlaceHandler{useCase: useCase}
}

// GetRecommendedPlaces はおすすめ店舗一覧を返します。
func (h *PlaceHandler) GetRecommendedPlaces(c echo.Context) error {
	for key := range c.QueryParams() {
		if _, ok := allowedPlaceQueries[key]; !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "未対応のクエリパラメータです: " + key})
		}
	}

	countValue := 0
	if rawCount := c.QueryParam("count"); rawCount != "" {
		parsedCount, err := strconv.Atoi(rawCount)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "countは数値で指定してください"})
		}
		countValue = parsedCount
	}
	pageValue := 0
	if rawPage := c.QueryParam("page"); rawPage != "" {
		parsedPage, err := strconv.Atoi(rawPage)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "pageは数値で指定してください"})
		}
		pageValue = parsedPage
	}
	if pageValue <= 0 {
		pageValue = defaultPage
	}
	if countValue <= 0 {
		countValue = defaultCount
	}
	if countValue > maxCount {
		countValue = maxCount
	}
	favorite := false
	if rawFavorite := c.QueryParam("favorite"); rawFavorite != "" {
		parsedFavorite, err := strconv.ParseBool(rawFavorite)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "favoriteはtrue/falseで指定してください"})
		}
		favorite = parsedFavorite
	}
	userID := c.QueryParam("userId")
	if favorite && userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "favorite=true の場合は userId が必要です"})
	}

	condition := placeDomain.SearchCondition{
		FishName: c.QueryParam("fishName"),
		Keyword:  c.QueryParam("keyword"),
		CityCode: c.QueryParam("cityCode"),
		UserID:   userID,
		Favorite: favorite,
		Count:    countValue,
		Page:     pageValue,
	}

	places, err := h.useCase.SearchRecommendedPlaces(c.Request().Context(), condition)
	if err != nil {
		if errors.Is(err, placeUseCase.ErrSmallAreaNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, placeRecommendationsResponse{
		Page:    condition.Page,
		PerPage: condition.Count,
		Count:   len(places),
		Items:   places,
	})
}

// UpdatePlaceFavorite は店舗のお気に入り状態を更新します。
func (h *PlaceHandler) UpdatePlaceFavorite(c echo.Context) error {
	for key := range c.QueryParams() {
		if _, ok := allowedPlaceFavoriteQueries[key]; !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "未対応のクエリパラメータです: " + key})
		}
	}

	queryUserID := c.QueryParam("userId")
	queryPlaceID := c.QueryParam("placeId")
	queryFavoriteRaw := c.QueryParam("favorite")
	queryFavorite := false
	queryFavoriteSet := false
	if queryFavoriteRaw != "" {
		parsedFavorite, err := strconv.ParseBool(queryFavoriteRaw)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "favoriteはtrue/falseで指定してください"})
		}
		queryFavorite = parsedFavorite
		queryFavoriteSet = true
	}

	var req updateFavoriteRequest
	if queryUserID == "" || queryPlaceID == "" || !queryFavoriteSet {
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストが不正です"})
		}
	}

	if queryUserID != "" {
		req.UserID = queryUserID
	}
	if queryPlaceID != "" {
		req.PlaceID = queryPlaceID
	}
	if queryFavoriteSet {
		req.Favorite = queryFavorite
	}

	err := h.useCase.UpdatePlaceFavorite(c.Request().Context(), req.UserID, req.PlaceID, req.Favorite)
	if err != nil {
		switch {
		case errors.Is(err, placeUseCase.ErrFavoriteUserIDRequired),
			errors.Is(err, placeUseCase.ErrFavoritePlaceIDRequired):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, placeUseCase.ErrFavoriteUserNotFound),
			errors.Is(err, placeUseCase.ErrFavoritePlaceNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "お気に入り更新に失敗しました"})
		}
	}

	return c.NoContent(http.StatusNoContent)
}
