package main

import (
	"log"
	"time"

	"fish-tech/internal/infrastructure/router"
	"fish-tech/internal/shared/timeutil"
)

// @title Fish-Tech API
// @version 1.0
// @description Fish-Tech バックエンド API 仕様です。
// @BasePath /api
func main() {
	time.Local = timeutil.JSTLocation()

	e, err := router.NewRouter()
	if err != nil {
		log.Fatal(err)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
