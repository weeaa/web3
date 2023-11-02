package groups

import (
	"github.com/gin-gonic/gin"
	"github.com/weeaa/nft/api/controllers"
)

func InitUserRoutes(router *gin.RouterGroup, controller *controllers.UserController) {
	user := router.Group("/user")

	{
		user.GET("", controller.Get)
		user.PUT("", controller.Update)
		user.POST("", controller.Insert)
		user.DELETE("", controller.Remove)
	}
}
