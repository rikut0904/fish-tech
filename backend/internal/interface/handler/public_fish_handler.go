package handler

import (
	"net/http"

	adminUseCase "fish-tech/internal/usecase/admin"

	"github.com/labstack/echo/v4"
)

type publicFishResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	LinkURL     string `json:"linkUrl"`
}

type publicPairResponse struct {
	ID      string `json:"id"`
	FishIDa string `json:"fishIdA"`
	FishIDb string `json:"fishIdB"`
	Score   int    `json:"score"`
	Memo    string `json:"memo"`
}

// PublicFishHandler は一般公開向けの魚データHTTPハンドラーです。
type PublicFishHandler struct {
	useCase adminUseCase.UseCase
}

// NewPublicFishHandler は一般公開向けの魚データハンドラーを作成します。
func NewPublicFishHandler(useCase adminUseCase.UseCase) *PublicFishHandler {
	return &PublicFishHandler{useCase: useCase}
}

// ListFishes は魚一覧を返します。
func (h *PublicFishHandler) ListFishes(c echo.Context) error {
	fishes, err := h.useCase.ListFishes(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "魚一覧の取得に失敗しました"})
	}

	responses := make([]publicFishResponse, 0, len(fishes))
	for _, fish := range fishes {
		responses = append(responses, publicFishResponse{
			ID:          fish.ID,
			Name:        fish.Name,
			Category:    fish.Category,
			Description: fish.Description,
			ImageURL:    fish.ImageURL,
			LinkURL:     fish.LinkURL,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{"items": responses})
}

// ListPairs は魚相性一覧を返します。
func (h *PublicFishHandler) ListPairs(c echo.Context) error {
	pairs, err := h.useCase.ListPairs(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "相性一覧の取得に失敗しました"})
	}

	responses := make([]publicPairResponse, 0, len(pairs))
	for _, pair := range pairs {
		responses = append(responses, publicPairResponse{
			ID:      pair.ID,
			FishIDa: pair.FishIDa,
			FishIDb: pair.FishIDb,
			Score:   pair.Score,
			Memo:    pair.Memo,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{"items": responses})
}
