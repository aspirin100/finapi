package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aspirin100/finapi/internal/app"
	"github.com/aspirin100/finapi/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("current config:", cfg)

	application, err := app.New(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = application.Run()
	if err != nil {
		log.Fatal(err)
	}

}
