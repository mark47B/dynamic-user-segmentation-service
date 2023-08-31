package api_handlers

import (
	"dynamic-user-segmentation-service/core"
	infrastructure "dynamic-user-segmentation-service/infrastructure/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	DbInteractor *infrastructure.Repository
}

func NewUserController(repoInteractor *infrastructure.Repository) *UserController {
	return &UserController{
		DbInteractor: repoInteractor,
	}
}

func (controller *UserController) SelectUserByUUID(c *gin.Context) {
	user_uuid, _ := uuid.Parse(c.Param("uuid"))
	user, err := controller.DbInteractor.U.GetUserByUUID(user_uuid)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
	}
	c.IndentedJSON(http.StatusOK, user)
	return
}

func (controller *UserController) SelectUserSlugsByUUID(c *gin.Context) {
	user_uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid user uuid"})
		return
	}
	user_slugs, err := controller.DbInteractor.U.SelectUserSlugsByUUID(user_uuid)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal error"})
		return
	}
	c.IndentedJSON(http.StatusOK, user_slugs)
	return
}

func (controller *UserController) CreateUser(c *gin.Context) {
	var userData core.UserRequestCreate

	if err := c.BindJSON(&userData); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error while binding JSON"})
		return
	}

	user_uuid, err := controller.DbInteractor.U.CreateUser(userData)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error while creating User"})
		return
	}
	c.IndentedJSON(http.StatusCreated, user_uuid)
	return

}
func (controller *UserController) GetAllUsers(c *gin.Context) {
	users := *new([]core.User)
	users, err := controller.DbInteractor.U.GetAll()
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "users not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, users)
	return
}

func (controller *UserController) ChangeUserSlugs(c *gin.Context) {
	user_uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid user uuid"})
		return
	}

	var userChangeInfo core.UserPut
	if err := c.BindJSON(&userChangeInfo); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error while binding JSON"})
		return
	}
	if len(userChangeInfo.Delete_slugs) != 0 {
		err = controller.DbInteractor.DeleteSlugsForUser(user_uuid, userChangeInfo.Delete_slugs)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error while deleteing slugs for user"})
			return
		}
	}
	if len(userChangeInfo.Add_slugs) != 0 {
		err = controller.DbInteractor.AddSlugToUser(user_uuid, userChangeInfo.Add_slugs)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error while adding slugs for user"})
			return
		}
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully deleted and added slugs"})
	return
}
