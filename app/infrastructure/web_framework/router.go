package webframework

import (
	"dynamic-user-segmentation-service/infrastructure/database"
	api_handlers "dynamic-user-segmentation-service/interfaces"

	"github.com/gin-gonic/gin"
)

func InitGinRouter(db_interactor *database.Repository) *gin.Engine {
	router := gin.Default()
	api_v1_router := router.Group("api/v1")
	initUserEndpoints(api_v1_router, db_interactor)
	initSlugEndpoints(api_v1_router, db_interactor)
	// router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}

func initUserEndpoints(router *gin.RouterGroup, db_interactor *database.Repository) {
	userController := api_handlers.NewUserController(db_interactor)
	userRoutes := router.Group("/user")
	{
		userRoutes.GET("/:uuid", userController.SelectUserByUUID)
		userRoutes.GET("/:uuid/slugs", userController.SelectUserSlugsByUUID)
		userRoutes.PUT("/:uuid", userController.ChangeUserSlugs)
		userRoutes.GET("/", userController.GetAllUsers)
		// userRoutes.POST("/", createUser)
	}
}

func initSlugEndpoints(router *gin.RouterGroup, db_interactor *database.Repository) {
	slugController := api_handlers.NewSlugController(db_interactor)

	slugRoutes := router.Group("/slug")
	{
		slugRoutes.POST("/", slugController.CreateSlug)
		slugRoutes.DELETE("/:name", slugController.DeleteSlug)
	}
}

// -----------------------------------
// import (
// 	_ "swag-gin-demo/docs"

// 	"github.com/gin-gonic/gin"
// 	swaggerFiles "github.com/swaggo/files"
// 	ginSwagger "github.com/swaggo/gin-swagger"
// )

// // todo represents data about a task in the todo list
// type todo struct {
// 	ID   string `json:"id"`
// 	Task string `json:"task"`
// }

// // message represents request response with a message
// type message struct {
// 	Message string `json:"message"`
// }

// // todo slice to seed todo list data
// var todoList = []todo{
// 	{"1", "Learn Go"},
// 	{"2", "Build an API with Go"},
// 	{"3", "Document the API with swag"},
// }

// // @title Go + Gin Todo API
// // @version 1.0
// // @description This is a sample server todo server. You can visit the GitHub repository at https://github.com/LordGhostX/swag-gin-demo

// // @contact.name API Support
// // @contact.url http://www.swagger.io/support
// // @contact.email support@swagger.io

// // @license.name MIT
// // @license.url https://opensource.org/licenses/MIT

// // @host localhost:8080
// // @BasePath /
// // @query.collection.format multi
