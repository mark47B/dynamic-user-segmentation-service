package main

import (
	"dynamic-user-segmentation-service/infrastructure/database"
	webframework "dynamic-user-segmentation-service/infrastructure/web_framework"
	"dynamic-user-segmentation-service/settings"
	"log"
	"os"
)

func main() {
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
	_ = Storage

	router := webframework.InitGinRouter(Storage)
	router.Run("localhost:8080")

	// // Testing User Storage
	// user_uuid, err := userStorage.CreateUser(core.UserRequestCreate{Username: "Stepan"})
	// user, err := userStorage.GetUserByUUID(user_uuid)
	// fmt.Println(user)

	// // Testing common storage
	// err = Storage.AddSlugToUser(user_uuid, []string{"AVITO_DISCOUNT_10", "AVITO_PERFORMANCE_VAS"})

	// // Testing Slug Storage
	// slug_id, err := slugStorage.CreateSlug(&core.SlugRequestAdd{Name: "AVITO_DISCOUNT_50"})
	// fmt.Println(slug_id)
	// err = slugStorage.DeleteSlugByName("AVITO_DISCOUNT_50")
	// fmt.Println(slug_id)
	// user, err := userStorage.GetUserByUUID(user_uuid)
	// fmt.Println(user)

}
