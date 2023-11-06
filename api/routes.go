package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/weeaa/nft/api/controllers"
	"github.com/weeaa/nft/api/groups"
	"github.com/weeaa/nft/internal/services"
	"net/http"
	"os"
	"time"
)

func InitRoutes(router *gin.Engine, userService *services.UserService) {
	configureRouter(router)

	credentials, ok := isBasicValid()
	if !ok {
		log.Error().Bool("basic auth set", ok)
		return
	}

	apiGroup := router.Group("/v1", gin.BasicAuth(credentials))

	{
		groups.InitUserRoutes(apiGroup, wireUserHandler(userService))
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

func isBasicValid() (map[string]string, bool) {
	credentials := make(map[string]string)
	username, ok := os.LookupEnv("BASIC_USERNAME")
	if !ok {
		return nil, false
	}
	password, ok := os.LookupEnv("BASIC_PASSWORD")
	credentials[username] = password
	return credentials, ok
}

func wireUserHandler(service *services.UserService) *controllers.UserController {
	return controllers.NewUserController(*service)
}
