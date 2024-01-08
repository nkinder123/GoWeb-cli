package routers

import (
	"github.com/gin-gonic/gin"
	"goWeb/controllers"
)

func InitAdminRouters(r *gin.Engine) {
	adminRouters := r.Group("/admin")
	{
		adminRouters.GET("/", controllers.AdminController{}.AdminList)
	}
}
