package routers

import (
	"github.com/gin-gonic/gin"
	"goWeb/controllers"
)

func InitApiRouters(r *gin.Engine) {
	apiRouters := r.Group("/api")
	{
		apiRouters.GET("/", controllers.ApiController{}.ApiUse)

	}
}
