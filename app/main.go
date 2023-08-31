package main

import (
	"dynamic-user-segmentation-service/infrastructure/database"
	webframework "dynamic-user-segmentation-service/infrastructure/web_framework"
	"dynamic-user-segmentation-service/settings"
	"log"
	"os"
)

func main() {
	// Config init
	cfg := settings.MustLoad()

	// Database init
	db_connection, err := database.GetConnectionToDB(cfg.Database)
	userStorage, err := database.NewUserRepository(db_connection)
	if err != nil {
		log.Println("Error while init User Storage:", err)
		os.Exit(1)
	}
	slugStorage, err := database.NewSlugRepository(db_connection)
	if err != nil {
		log.Println("Error while init Slug Storage:", err)
		os.Exit(1)
	}
	Storage := &database.Repository{
		S: slugStorage,
		U: userStorage,
	}

	// Gin init
	router := webframework.InitGinRouter(Storage)
	router.Run(cfg.Server.SERVER_ADDRESS)

}
