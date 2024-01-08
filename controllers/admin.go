package controllers

import "github.com/gin-gonic/gin"

type AdminController struct {
}

func (con AdminController) AdminList(context *gin.Context) {
	context.String(200, "admin")
}
