package main

import (
	"fish-tech/internal/infrastructure/router"
)

func main() {
	e := router.NewRouter()
	e.Logger.Fatal(e.Start(":8080"))
}
