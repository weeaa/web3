package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/weeaa/nft/api/controllers"
	"github.com/weeaa/nft/api/groups"
	"github.com/weeaa/nft/internal/services"
	"net/http"
	"os"
	"time"
)

var credentials = map[string]string{
	os.Getenv("BASIC_USERNAME"): os.Getenv("BASIC_PASSWORD"),
}

func InitRoutes(router *gin.Engine, traderService *services.UserService) {
	apiGroup := router.Group("/v1", gin.BasicAuth(credentials))
	{
		groups.InitUserRoutes(apiGroup, wireTraderHandler(traderService))
	}
	apiGroup.POST("/ping", pong)
}

func pong(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}

func configureRouter(router *gin.Engine) {
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))
}

func wireTraderHandler(service *services.UserService) *controllers.UserController {
	return controllers.NewUserController(*service)
}
