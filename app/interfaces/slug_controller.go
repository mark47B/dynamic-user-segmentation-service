package api_handlers

import (
	"dynamic-user-segmentation-service/core"
	infrastructure "dynamic-user-segmentation-service/infrastructure/database"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SlugController struct {
	DbInteractor *infrastructure.Repository
}

func NewSlugController(repoInteractor *infrastructure.Repository) *SlugController {
	return &SlugController{
		DbInteractor: repoInteractor,
	}
}

func (controller *SlugController) CreateSlug(c *gin.Context) {
	op := "interfaces.slug_controller.CreateSlug"

	var slugCreate core.SlugRequestAdd
	if err := c.BindJSON(&slugCreate); err != nil {
		log.Println("Error in CreateSlug", op)
		c.Error(fmt.Errorf("Can't serialize your JSON"))
		return
	}
	Id, err := controller.DbInteractor.S.CreateSlug(&slugCreate)
	if err != nil {
		log.Println("Error while create slug in database", op)
		c.Error(fmt.Errorf("Error while create slug in database"))
		return
	}
	c.JSON(http.StatusCreated, Id)
	return
}
func (controller *SlugController) DeleteSlug(c *gin.Context) {
	op := "interfaces.slug_controller.DeleteSlug"

	slugName := c.Param("name")
	fmt.Println(slugName)

	err := controller.DbInteractor.S.DeleteSlugByName(slugName)
	if err != nil {
		log.Println("Error while deleting slug from database", op)
		c.Error(fmt.Errorf("Error while deleting slug from database"))
		return
	}
	return

}
