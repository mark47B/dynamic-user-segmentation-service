package api_handlers

import (
	"dynamic-user-segmentation-service/core"
	infrastructure "dynamic-user-segmentation-service/infrastructure/database"
	"fmt"
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
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
	}
	c.IndentedJSON(http.StatusOK, user)
	return
}

func (controller *UserController) SelectUserSlugsByUUID(c *gin.Context) {
	user_uuid, _ := uuid.Parse(c.Param("uuid"))
	user_slugs, err := controller.DbInteractor.U.SelectUserSlugsByUUID(user_uuid)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "slugs not found"})
	}
	c.IndentedJSON(http.StatusOK, user_slugs)
	return
}

func (controller *UserController) GetAllUsers(c *gin.Context) {
	users := *new([]core.User)
	users, err := controller.DbInteractor.U.GetAll()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "users not found"})
	}
	c.IndentedJSON(http.StatusOK, users)
	return
}

func (controller *UserController) ChangeUserSlugs(c *gin.Context) {
	op := "interfaces.api.ChangeUserSlugs CONTROLLER"

	user_uuid, _ := uuid.Parse(c.Param("uuid"))
	var userChangeInfo core.UserPut

	if err := c.BindJSON(&userChangeInfo); err != nil {
		c.Error(fmt.Errorf("Can't serialize your JSON"))
		return
	}
	err := controller.DbInteractor.DeleteSlugsForUser(user_uuid, userChangeInfo.Delete_slugs)
	err = controller.DbInteractor.AddSlugToUser(user_uuid, userChangeInfo.Add_slugs)
	if err != nil {
		log.Println("Error in adding slugs", op)
		c.Error(fmt.Errorf("Can't change slugs"))
		return
	}
	res := make(map[string]interface{}, 0)
	res["status"] = 200
	res["healthy"] = "OK"
	c.JSON(http.StatusOK, res)
	return
}
