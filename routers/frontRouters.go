package routers

import (
	"github.com/gin-gonic/gin"
	"goWeb/controllers"
)

func InitFrontRoters(r *gin.Engine) {
	frontRouters := r.Group("/")
	{
		frontRouters.GET("/", controllers.FrontController{}.FrontUse)
	}
}
