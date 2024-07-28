package response

/*响应格式
{
	"code":10001
	"msg":xx,
	"data":{}
}
*/

// 新建返回值code类型
type ResCode int

// 错误状态码
const (
	CodeSucess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExit
	CodeUserNotExit
	CodeInvalidPassword
	CodeServerBusy

	CodeInvalidToken
	CodeNeedLogin
)

// 状态码信息
var CodeMsg = map[ResCode]string{
	CodeSucess:          "success",
	CodeInvalidParam:    "参数错误",
	CodeUserExit:        "用户存在",
	CodeUserNotExit:     "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "网络繁忙",
	CodeInvalidToken:    "无效token",
	CodeNeedLogin:       "需要登陆",
}

// 获取状态码信息
func (r ResCode) Msg() string {
	msg, ok := CodeMsg[r]
	if !ok {
		return CodeMsg[CodeServerBusy]
	}
	return msg
}
