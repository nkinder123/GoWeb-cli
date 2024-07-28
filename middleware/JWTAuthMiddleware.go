package middleware

import (
	"github.com/gin-gonic/gin"
	"goWeb/pkg/jwt"
	"goWeb/pkg/response"
	"strings"
)

/*
客户端携带token有三种方式：1.请求头 2.请求体  3.url中
假设token放在header中的Authorization中并且使用bearer开头
Authorization：Bearer xxx.xxxx.xxxxx
*/
const CtxtUserIdKey = "userId"

func JWTAuthMiddleWare() func(context *gin.Context) {
	return func(context *gin.Context) {
		//获取Authorization的内容

		//请求头数据格式：Authorization：Bearer xxx.xxxx.xxxxx
		atu := context.Request.Header.Get("Authorization")
		if atu == "" {
			response.ResponseError(context, response.CodeNeedLogin)
			//不执行后面的中间件了，直接将当前，以及在此之前的中间件进行返回
			context.Abort()
			return
		}
		//处理Authorization中的内容获取token
		//在SplitN中2代表返回的切片个数
		parts := strings.SplitN(atu, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.ResponseError(context, response.CodeInvalidToken)
			context.Abort()
			return
		}
		//解析token
		claim, parTokeErr := jwt.ParasToken(parts[1])
		if parTokeErr != nil {
			response.ResponseError(context, response.CodeInvalidToken)
			context.Abort()
			return
		}
		//c.Next()是执行下一个中间件
		context.Set(CtxtUserIdKey, claim.UserId)
		context.Next()
	}
}
