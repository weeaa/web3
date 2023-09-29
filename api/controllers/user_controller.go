package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/weeaa/nft/internal/models"
	"github.com/weeaa/nft/internal/services"
	"github.com/weeaa/nft/pkg/utils"
	"net/http"
)

type UserController struct {
	userService services.UserService
}

// todo finir les fonctions + func db
func NewUserController(service services.UserService) *UserController {
	return &UserController{userService: service}
}

func (t UserController) InsertUser(c *gin.Context) {
	u, err := utils.UnmarshalJSONToStruct[models.TraderAddBody](c.Request.Body)
	if err != nil {

	}
	c.JSON(http.StatusOK, u)
}

func (t UserController) GetUser(c *gin.Context) {

}

func (t UserController) UpdateUser(c *gin.Context) {}

func (t UserController) RemoveUser(c *gin.Context) {
	u, err := utils.UnmarshalJSONToStruct[models.TraderRemoveBody](c.Request.Body)
	if err != nil {

	}
	c.JSON(http.StatusOK, u)
}
