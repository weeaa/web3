package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	repository_models "github.com/weeaa/nft/database/models"
	"github.com/weeaa/nft/internal/models"
	"github.com/weeaa/nft/internal/services"
	"github.com/weeaa/nft/pkg/utils"
	"net/http"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{userService: service}
}

func (t UserController) Insert(c *gin.Context) {
	response, err := utils.UnmarshalJSONToStruct[repository_models.FriendTechMonitor](c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	if err = t.userService.DB.Monitor.InsertUser(&response, context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, "ok")
}

func (t UserController) Get(c *gin.Context) {

}

func (t UserController) Update(c *gin.Context) {}

func (t UserController) Remove(c *gin.Context) {
	response, err := utils.UnmarshalJSONToStruct[models.TraderRemoveBody](c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
	}

	if err = t.userService.DB.Monitor.RemoveUser(response.BaseAddress, context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, "")
	}

	c.JSON(http.StatusOK, "ok")
}
