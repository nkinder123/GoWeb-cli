package controllers

import "github.com/gin-gonic/gin"

type FrontController struct {
}

func (con FrontController) FrontUse(context *gin.Context) {
	context.String(200, "front")
}
