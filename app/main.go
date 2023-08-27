package main

import (
	"dynamic-user-segmentation-service/interfaces/db"
	"dynamic-user-segmentation-service/settings"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

func main() {
	cfg := settings.MustLoad()

	userStorage, err := db.NewUserRepository(cfg.Database)
	if err != nil {
		log.Println("Error while init User Storage:", err)
		os.Exit(1)
	}

	_ = userStorage

	user, err := userStorage.GetUserByUUID(uuid.MustParse("00112233-4455-6677-8899-aabbccddeeff"))
	fmt.Println(user)
	// slugStorage, err := db.NewSlugRepository(cfg.Database)
}
