package routers

import (
	"github.com/gin-gonic/gin"
	"goWeb/logger"
)

func Setup() *gin.Engine {
	r := gin.Default()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	InitFrontRoters(r)
	InitAdminRouters(r)
	InitApiRouters(r)
	return r
}
