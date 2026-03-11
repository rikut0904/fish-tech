package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"fish-tech/internal/usecase/hello"
)

// HelloResponse は動作確認レスポンスです。
type HelloResponse struct {
	Message string `json:"message"`
}

// hello関連のHTTPリクエストを処理する
type HelloHandler struct {
	useCase hello.UseCase
}

// 新しいhelloハンドラーを作成する
func NewHelloHandler(useCase hello.UseCase) *HelloHandler {
	return &HelloHandler{
		useCase: useCase,
	}
}

// GET /helloリクエストを処理する
// @Summary 動作確認
// @Description API の動作確認メッセージを返します。
// @Tags Public
// @Produce json
// @Success 200 {object} HelloResponse
// @Router /hello [get]
func (h *HelloHandler) GetHello(c echo.Context) error {
	result := h.useCase.GetHello()
	return c.JSON(http.StatusOK, HelloResponse{Message: result.Message})
}
