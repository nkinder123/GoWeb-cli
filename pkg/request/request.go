package request

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goWeb/middleware"
)

var ErrorUserNotLogin = errors.New("用户未登陆")

func GetUserId(context *gin.Context) (userId int64, err error) {
	user, ok := context.Get(middleware.CtxtUserIdKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userId, conerr := user.(int64)
	if !conerr {
		err = ErrorUserNotLogin
		return
	}
	return
}
