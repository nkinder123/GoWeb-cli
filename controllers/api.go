package controllers

import (
	"github.com/gin-gonic/gin"
)

type ApiController struct {
}

func (con ApiController) ApiUse(context *gin.Context) {
	context.String(200, "api")
}
