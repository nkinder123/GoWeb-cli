package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"goWeb/logic"
	"goWeb/models"
	"goWeb/pkg/jwt"
	"goWeb/pkg/request"
	"goWeb/pkg/response"
	"goWeb/pkg/valitor"
	"net/http"
	"strconv"
)

type FrontController struct {
}

func (con FrontController) FrontUse(context *gin.Context) {
	context.String(200, "front")
}

// 注册
func (con FrontController) SignUp(context *gin.Context) {
	//1.接受参数、参数校验
	var userSign models.UserConfirm
	if err := context.ShouldBindJSON(&userSign); err != nil {
		zap.L().Error("param is vailid", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ResponseErrorWithMsg(context, response.CodeInvalidParam, err)
			return
		}
		response.ResponseErrorWithMsg(
			context,
			response.CodeServerBusy,
			valitor.RemoveTopStruct(errs.Translate(valitor.GetTrans())))
		return
	}
	//2.业务处理
	if logicerr := logic.SignUp(&userSign); logicerr != nil {
		zap.L().Error("logic has error", zap.Error(logicerr))
		if errors.Is(logicerr, logic.ErrorPassword) {
			response.ResponseErrorWithMsg(context, response.CodeInvalidPassword, logicerr.Error())
		}
		if errors.Is(logicerr, logic.ErrorUserExit) {
			response.ResponseErrorWithMsg(context, response.CodeUserExit, logicerr.Error())
		}
		return
	}
	//3.返回响应
	response.ResponseSuccess(context, nil)
}

// Login校验
func (con FrontController) Login(context *gin.Context) {
	var userLogin models.UserLogin
	//参数校验
	if err := context.ShouldBindJSON(&userLogin); err != nil {
		zap.L().Error("param is vailid", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ResponseError(context, response.CodeInvalidParam)
		}
		response.ResponseErrorWithMsg(context, response.CodeServerBusy, valitor.RemoveTopStruct(errs.Translate(valitor.GetTrans())))
		return
	}
	//逻辑处理
	//var loginuser *models.Users
	//var logicerr error
	if _, logicerr := logic.Login(&userLogin); logicerr != nil {
		zap.L().Error("login falied")
		response.ResponseError(context, response.CodeInvalidPassword)
		return
	}
	//生成token
	loginuser, _ := logic.Login(&userLogin)
	token, tokenerr := jwt.GenToken(loginuser.Username, loginuser.UserId)
	fmt.Printf("err", tokenerr)
	if tokenerr != nil {
		response.ResponseErrorWithMsg(context, response.CodeServerBusy, zap.Error(tokenerr))
		return
	}
	zap.L().Info("login success")
	response.ResponseSuccess(context, token)
	return
}

// 携带token验证
func (con FrontController) Ping(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"msg": "登陆成功",
	})
}

// 前端展示社区列表
func (con FrontController) GetCommunity(context *gin.Context) {
	//调用logic层返回数据
	data, err := logic.CommunityList()
	if err != nil {
		zap.L().Error("communityList logic has error", zap.Error(err))
		response.ResponseError(context, response.CodeServerBusy)
	}
	response.ResponseSuccess(context, data)
}

// 前端展示社区detail
func (con FrontController) GetCommunityDetail(context *gin.Context) {
	idstr := context.Param("id")
	id, strpinterr := strconv.ParseInt(idstr, 10, 64)
	if strpinterr != nil {
		response.ResponseError(context, response.CodeServerBusy)
	}
	item, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("communitydDetail has error", zap.Error(err))
		response.ResponseError(context, response.CodeInvalidPassword)
		return
	}
	response.ResponseSuccess(context, item)
	return
}

func (con FrontController) NewPost(context *gin.Context) {
	//1.参数获取、参数校验
	var post models.Post
	if err := context.ShouldBindJSON(&post); err != nil {
		zap.L().Error("param is vailid", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ResponseError(context, response.CodeInvalidParam)
		}
		response.ResponseErrorWithMsg(context, response.CodeServerBusy, valitor.RemoveTopStruct(errs.Translate(valitor.GetTrans())))
		return
	}
	//2.获取context中authorId
	id, getIdErr := request.GetUserId(context)
	if getIdErr != nil {
		zap.L().Error("用户需要登陆")
		response.ResponseError(context, response.CodeNeedLogin)
		return
	}
	var postpara models.PostPara
	postpara.AuthorId = id
	postpara.CommunityId = post.CommunityId
	postpara.Title = post.Title
	postpara.Content = post.Content
	//3.logic判断
	if logicerr := logic.NewPost(&postpara); logicerr != nil {
		zap.L().Error("新建帖子逻辑错误", zap.Error(logicerr))
		response.ResponseErrorWithMsg(context, response.CodeServerBusy, logicerr)
		return
	}
	//4.返回响应
	response.ResponseSuccess(context, nil)
}

func (con FrontController) PostDetail(context *gin.Context) {
	//获取帖子id值  from url
	pid := context.Param("id")
	id, err := strconv.ParseInt(pid, 10, 64)
	if err != nil {
		response.ResponseError(context, response.CodeInvalidParam)
		return
	}
	if len(pid) == 0 {
		zap.L().Error("id参数值传递错误")
		response.ResponseError(context, response.CodeInvalidParam)
		return
	}
	//logic
	if logicerr, _, _, _ := logic.GetPostDetail(id); logicerr != nil {
		zap.L().Error("getPostDetaild logic has error")
		response.ResponseError(context, response.CodeServerBusy)
		return
	}
	_, authorname, community, post := logic.GetPostDetail(id)
	var postDetail = &models.PostDetail{
		authorname, community, post,
	}
	response.ResponseSuccess(context, postDetail)
	return
}

// 实现分页操作
func (con FrontController) GetPostList(context *gin.Context) {
	//获取参数，参数校验
	pageStr := context.Query("page")
	sizeStr := context.Query("size")
	page, perr := strconv.ParseInt(pageStr, 10, 64)
	if perr != nil {
		page = 1
	}
	size, serr := strconv.ParseInt(sizeStr, 10, 64)
	if serr != nil {
		size = 10
	}
	data, logicerr := logic.GetPostList(page, size)
	if logicerr != nil {
		zap.L().Error("postList logic has error", zap.Error(logicerr))
		response.ResponseError(context, response.CodeServerBusy)
		return
	}
	response.ResponseSuccess(context, data)
	return
}

//投票的功能需求
/*
	1.用户投票的数据「1，帖子的Id   2.赞成还是反对[-1,0,1]」
	2.拿到当前用户的Id
	3.投一票432分，一个帖子获得200张赞成票就可以续一天
	4.投票的几种情况：direction=1   之前没有投票过
								 之前投反对票，现在改成赞成票
					direction=0  之前投反对票 ---》取消
								 之前赞成---》取消
					direction=-1 之前没投票--》反对
								 之前赞成现在反对票
	5.投票限制：每个帖子发表之日起，一个星期内可以投票，超过一个星期就不允许投票
				到期之后：帖子的赞成和反对就将其存储到mysql中，删除那个keyPostVotedPrefix中
*/

func (con FrontController) VotedPost(context *gin.Context) {
	var vote models.Voted
	//传输绑定多个变量
	if err := context.ShouldBindJSON(&vote); err != nil {
		zap.L().Error("param is vailid", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ResponseError(context, response.CodeInvalidParam)
		}
		response.ResponseErrorWithMsg(context, response.CodeServerBusy, valitor.RemoveTopStruct(errs.Translate(valitor.GetTrans())))
		return
	}
	current_id, userErr := request.GetUserId(context)
	if userErr != nil {
		response.ResponseError(context, response.CodeNeedLogin)
		return
	}
	logicerr := logic.PostVoted(strconv.Itoa(int(current_id)), strconv.Itoa(int(vote.PostId)), float64(vote.Direction))
	if logicerr != nil {
		response.ResponseErrorWithMsg(context, response.CodeServerBusy, logicerr)
		return
	}
	response.ResponseSuccess(context, nil)
	return
}
