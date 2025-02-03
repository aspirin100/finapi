package main

import (
	"log"

	"github.com/aspirin100/finapi/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	println(cfg)
}
