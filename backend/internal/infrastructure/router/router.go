package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"fish-tech/internal/interface/handler"
	"fish-tech/internal/usecase/hello"
)

// 全ルートが設定された新しいEchoルーターを作成する
func NewRouter() *echo.Echo {
	e := echo.New()

	// ミドルウェア
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// ユースケースの初期化
	helloUseCase := hello.NewHelloUseCase()

	// ハンドラーの初期化
	helloHandler := handler.NewHelloHandler(helloUseCase)

	// ルーティング
	api := e.Group("/api")
	{
		api.GET("/hello", helloHandler.GetHello)
	}

	return e
}
