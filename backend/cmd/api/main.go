package main

import (
	"log"

	"fish-tech/internal/infrastructure/router"
)

func main() {
	e, err := router.NewRouter()
	if err != nil {
		log.Fatal(err)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
