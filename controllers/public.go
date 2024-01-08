package controllers

import "github.com/gin-gonic/gin"

type PublicController struct {
}

func (con PublicController) PublicUse(context *gin.Context) {
	context.String(200, "public")
}
