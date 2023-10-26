package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	repository_models "github.com/weeaa/nft/database/models"
	"github.com/weeaa/nft/internal/models"
	"github.com/weeaa/nft/internal/services"
	"net/http"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{userService: service}
}

func (t UserController) Insert(c *gin.Context) {
	var req repository_models.FriendTechMonitor
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := t.userService.DB.Monitor.InsertUser(&req, context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func (t UserController) Get(c *gin.Context) {
	var req map[string]any
	var buf bytes.Buffer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := t.userService.DB.Monitor.GetUserByAddress(req["address"].(string), context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = json.NewEncoder(&buf).Encode(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (t UserController) Update(c *gin.Context) {}

func (t UserController) Remove(c *gin.Context) {
	var req models.UserRemoveBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := t.userService.DB.Monitor.RemoveUser(req.BaseAddress, context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
