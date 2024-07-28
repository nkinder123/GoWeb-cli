package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goWeb/logger"
	"goWeb/pkg/valitor"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	if err := valitor.InitTrans("zh"); err != nil {
		fmt.Printf("init validator falied,err:%v", err)
	}
	r := gin.Default()
	r.Use(logger.GinLogger())
	//, logger.GinRecovery(true)
	InitFrontRoters(r)
	InitAdminRouters(r)
	InitApiRouters(r)
	return r
}
