package routers

import (
	"github.com/gin-gonic/gin"
	"goWeb/controllers"
	"goWeb/middleware"
)

func InitFrontRoters(r *gin.Engine) {
	//注册路由
	r.POST("/signup", controllers.FrontController{}.SignUp)

	r.POST("/login", controllers.FrontController{}.Login)
	frontRouters := r.Group("/", middleware.JWTAuthMiddleWare())
	{

		//携带token验证
		frontRouters.GET("/ping", controllers.FrontController{}.Ping)
		//前端展示社区接口
		frontRouters.GET("/community", controllers.FrontController{}.GetCommunity)
		//前端社区展示信息
		frontRouters.GET("/community/:id", controllers.FrontController{}.GetCommunityDetail)
		//新建帖子
		frontRouters.POST("/post", controllers.FrontController{}.NewPost)
		//帖子详情
		frontRouters.GET("/post/:id", controllers.FrontController{}.PostDetail)
		//实现枫分页操作
		frontRouters.GET("/posts", controllers.FrontController{}.GetPostList)
		//实现投票的功能
		frontRouters.POST("/voted", controllers.FrontController{}.VotedPost)
	}
}
