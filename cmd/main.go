package main

import (
	"log"

	"summer/practice2022/internal/config"
	"summer/practice2022/internal/db"
	"summer/practice2022/internal/logic"
	"summer/practice2022/internal/server"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	logic := logic.NewLogic(cfg, db)

	server := server.NewServer(cfg, logic)

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}
