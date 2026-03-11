package main

import (
	"log"
	"time"

	"fish-tech/internal/infrastructure/router"
	"fish-tech/internal/shared/timeutil"
)

func main() {
	time.Local = timeutil.JSTLocation()

	e, err := router.NewRouter()
	if err != nil {
		log.Fatal(err)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
