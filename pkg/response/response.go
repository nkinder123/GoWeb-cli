package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// 返回错误信息
func ResponseError(c *gin.Context, code ResCode) {
	rd := &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}

// 返回带有信息的错误信息
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	rd := &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}

// 返回成功信息
func ResponseSuccess(c *gin.Context, data interface{}) {
	rd := ResponseData{
		Code: CodeSucess,
		Msg:  CodeSucess.Msg(),
		Data: data,
	}
	c.JSON(http.StatusOK, rd)
}
