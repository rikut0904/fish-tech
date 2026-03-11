package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"

	adminUseCase "fish-tech/internal/usecase/admin"
)

// AdminHandler は管理画面向けHTTPハンドラーです。
type AdminHandler struct {
	useCase  adminUseCase.UseCase
	uploader ImageUploader
}

// ImageUploader は画像アップロード機能を表すインターフェースです。
type ImageUploader interface {
	Enabled() bool
	Upload(ctx context.Context, filename string, data []byte) (string, string, error)
}

// NewAdminHandler は新しい管理画面ハンドラーを作成します。
func NewAdminHandler(useCase adminUseCase.UseCase, uploader ImageUploader) *AdminHandler {
	return &AdminHandler{useCase: useCase, uploader: uploader}
}

type fishResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	Description  string `json:"description"`
	ImageURL     string `json:"imageUrl"`
	ImageMediaID string `json:"imageMediaId,omitempty"`
	LinkURL      string `json:"linkUrl"`
}

type pairResponse struct {
	ID      string `json:"id"`
	FishIDa string `json:"fishIdA"`
	FishIDb string `json:"fishIdB"`
	Score   int    `json:"score"`
	Memo    string `json:"memo"`
}

type uploadImageResponse struct {
	ImageURL     string `json:"imageUrl"`
	ImageMediaID string `json:"imageMediaId,omitempty"`
}

// CreateFish は魚を登録します。
func (h *AdminHandler) CreateFish(c echo.Context) error {
	fish, err := h.useCase.CreateFish(
		c.Request().Context(),
		c.QueryParam("name"),
		c.QueryParam("category"),
		c.QueryParam("description"),
		c.QueryParam("imageUrl"),
		c.QueryParam("imageMediaId"),
		c.QueryParam("linkUrl"),
	)
	if err != nil {
		if errors.Is(err, adminUseCase.ErrInvalidFishName) || errors.Is(err, adminUseCase.ErrInvalidFishCategory) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "魚の登録に失敗しました"})
	}

	return c.JSON(http.StatusCreated, fishResponse{
		ID:           fish.ID,
		Name:         fish.Name,
		Category:     fish.Category,
		Description:  fish.Description,
		ImageURL:     fish.ImageURL,
		ImageMediaID: fish.ImageMediaID,
		LinkURL:      fish.LinkURL,
	})
}

// UploadFishImage は魚画像をGoogle Photosへアップロードします。
func (h *AdminHandler) UploadFishImage(c echo.Context) error {
	if h.uploader == nil || !h.uploader.Enabled() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Google Photos設定が不足しています"})
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "画像ファイルが指定されていません"})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "画像ファイルを開けません"})
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "画像ファイルの読み込みに失敗しました"})
	}
	if len(data) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "空の画像はアップロードできません"})
	}

	name := filepath.Base(fileHeader.Filename)
	if name == "." || name == "/" || name == "" {
		name = "upload.jpg"
	}

	imageURL, imageMediaID, err := h.uploader.Upload(c.Request().Context(), name, data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, uploadImageResponse{
		ImageURL:     imageURL,
		ImageMediaID: imageMediaID,
	})
}

// DeleteFish は魚を削除します。
func (h *AdminHandler) DeleteFish(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "idは必須です"})
	}

	if err := h.useCase.DeleteFish(c.Request().Context(), id); err != nil {
		if errors.Is(err, adminUseCase.ErrFishNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "魚の削除に失敗しました"})
	}

	return c.NoContent(http.StatusNoContent)
}

// CreatePair は魚相性を登録します。
func (h *AdminHandler) CreatePair(c echo.Context) error {
	score, err := strconv.Atoi(c.QueryParam("score"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "scoreは数値で指定してください"})
	}

	pair, err := h.useCase.CreatePair(
		c.Request().Context(),
		c.QueryParam("fishIdA"),
		c.QueryParam("fishIdB"),
		score,
		c.QueryParam("memo"),
	)
	if err != nil {
		switch {
		case errors.Is(err, adminUseCase.ErrInvalidPair),
			errors.Is(err, adminUseCase.ErrPairAlreadyExists),
			errors.Is(err, adminUseCase.ErrFishNotFound):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "相性の登録に失敗しました"})
		}
	}

	return c.JSON(http.StatusCreated, pairResponse{
		ID:      pair.ID,
		FishIDa: pair.FishIDa,
		FishIDb: pair.FishIDb,
		Score:   pair.Score,
		Memo:    pair.Memo,
	})
}

// DeletePair は魚相性を削除します。
func (h *AdminHandler) DeletePair(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "idは必須です"})
	}

	if err := h.useCase.DeletePair(c.Request().Context(), id); err != nil {
		if errors.Is(err, adminUseCase.ErrPairNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "相性データが見つかりません"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "相性の削除に失敗しました"})
	}

	return c.NoContent(http.StatusNoContent)
}
