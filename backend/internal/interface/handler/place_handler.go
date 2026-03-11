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

// RecommendedPlaceResponse はおすすめ店舗レスポンスです。
type RecommendedPlaceResponse struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Lat     string `json:"lat"`
	Lng     string `json:"lng"`
	Coupon  string `json:"coupon"`
	Genre   string `json:"genre"`
	Card    string `json:"card"`
	Logo    string `json:"logo"`
}

// PlaceRecommendationsResponse はおすすめ店舗一覧レスポンスです。
type PlaceRecommendationsResponse struct {
	Page    int                        `json:"page"`
	PerPage int                        `json:"perPage"`
	Count   int                        `json:"count"`
	Items   []RecommendedPlaceResponse `json:"items"`
}

// UpdateFavoriteRequest は店舗お気に入り更新リクエストです。
type UpdateFavoriteRequest struct {
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
// @Summary おすすめ店舗を取得
// @Description 条件に応じておすすめ店舗一覧を返します。
// @Tags Place
// @Produce json
// @Param fishName query string false "魚名"
// @Param keyword query string false "検索キーワード"
// @Param cityCode query string false "HotPepper small_area コード"
// @Param userId query string false "ユーザーID"
// @Param favorite query boolean false "お気に入りのみ取得"
// @Param count query int false "取得件数"
// @Param page query int false "ページ番号"
// @Success 200 {object} PlaceRecommendationsResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /places/recommendations [get]
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

	items := make([]RecommendedPlaceResponse, 0, len(places))
	for _, place := range places {
		items = append(items, RecommendedPlaceResponse{
			Name:    place.Name,
			Address: place.Address,
			Lat:     place.Lat,
			Lng:     place.Lng,
			Coupon:  place.Coupon,
			Genre:   place.Genre,
			Card:    place.Card,
			Logo:    place.Logo,
		})
	}

	return c.JSON(http.StatusOK, PlaceRecommendationsResponse{
		Page:    condition.Page,
		PerPage: condition.Count,
		Count:   len(places),
		Items:   items,
	})
}

// UpdatePlaceFavorite は店舗のお気に入り状態を更新します。
// @Summary 店舗のお気に入り状態を更新
// @Description 店舗のお気に入り状態を更新します。
// @Tags Place
// @Accept json
// @Produce json
// @Param userId query string false "ユーザーID"
// @Param placeId query string false "店舗ID"
// @Param favorite query boolean false "お気に入り状態"
// @Param request body UpdateFavoriteRequest false "後方互換用リクエストボディ"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /places/favorite [patch]
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

	var req UpdateFavoriteRequest
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
