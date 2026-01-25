package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"fish-tech/internal/usecase/hello"
)

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
func (h *HelloHandler) GetHello(c echo.Context) error {
	result := h.useCase.GetHello()
	return c.JSON(http.StatusOK, result)
}
