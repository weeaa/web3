package groups

import (
	"github.com/gin-gonic/gin"
	"github.com/weeaa/nft/api/controllers"
)

func InitUserRoutes(router *gin.RouterGroup, controller *controllers.UserController) {
	user := router.Group("/user")
	{
		user.GET("/", controller.GetUser)
		user.PUT("/", controller.UpdateUser)
		user.POST("/", controller.InsertUser)
		user.DELETE("/", controller.RemoveUser)
	}
}
