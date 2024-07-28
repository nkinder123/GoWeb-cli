package logic

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"goWeb/dao/mysql"
	"goWeb/dao/redisCon"
	"goWeb/models"
	"goWeb/pkg/md5"
	"goWeb/pkg/snowflake"
	"math"
	"time"
)

var (
	ErrorUserExit    = errors.New("用户已经存在")
	ErrorUserNotExit = errors.New("用户名不存在")
	ErrorPassword    = errors.New("密码错误")
	ErrorPostExpired = errors.New("帖子过期")
	ErrorUpdatePost  = errors.New("更改帖子score值出错")
)

const (
	ExpireTime   = 7 * 24 * 3600
	ScorePerVote = 432
)

// 注册逻辑
func SignUp(user *models.UserConfirm) (err error) {
	//1.判断用户是否存在
	userinfo := models.Users{}
	if infoerr := mysql.GetDB().Where("username=?", user.Username).First(&userinfo).Error; infoerr != nil {
		//2.不存在生成uid
		userId := snowflake.GenID()
		//3.密码加密
		hdpssword := md5.EncodeMd5(user.Password)
		userinfo = models.Users{
			Username:   user.Username,
			UserId:     userId,
			Password:   hdpssword,
			CreateTime: int(time.Now().Unix()),
		}
		//3.写入数据库
		if dberr := mysql.GetDB().Create(&userinfo).Error; dberr != nil {
			mysql.GetDB().Rollback()
			fmt.Printf("create data has err ,err:%v", dberr)
			return
		}
		return
	}
	zap.L().Info("用户存在")
	return ErrorUserExit
}

// 登陆逻辑
func Login(user *models.UserLogin) (reuser *models.Users, err error) {
	//判断有没有用户名
	var userinfo models.Users
	if usererr := mysql.GetDB().Where("username=?", user.Username).First(&userinfo).Error; usererr != nil {
		zap.L().Error("用户名不存在")
		return nil, ErrorUserNotExit
	}
	fmt.Sprintf("userinfo:%v", userinfo)
	//判断用户名密码是否正确
	if userinfo.Password != md5.EncodeMd5(user.Password) {
		zap.L().Error("密码错误")
		return nil, ErrorPassword
	}
	return &userinfo, nil
}

// 前端展示逻辑
func CommunityList() (*[]models.CommunityList, error) {

	var dbcommunity []models.Community
	if err := mysql.GetDB().Find(&dbcommunity).Error; err != nil {
		zap.L().Error("get data has error", zap.Error(err))
		return nil, err
	}
	var result = make([]models.CommunityList, len(dbcommunity))
	for key, item := range dbcommunity {
		result[key].Id = item.CommunityId
		result[key].CommunityName = item.CommunityName
	}
	return &result, nil
}

// 获取社区detail
func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	var communityItem models.Community
	if err := mysql.GetDB().Where("community_id=?", id).First(&communityItem).Error; err != nil {
		zap.L().Error("get data has err", zap.Error(err))
		return nil, err
	}

	var communityDetail = models.CommunityDetail{}
	communityDetail.Id = communityItem.CommunityId
	communityDetail.CommunityName = communityItem.CommunityName
	communityDetail.Introduction = communityItem.Introduction
	communityDetail.CreateTime = communityItem.CreateTime
	return &communityDetail, nil
}

// 新建帖子
func NewPost(postpara *models.PostPara) error {
	//生成帖子id,生成potmodel
	var post models.Post
	post.PostId = snowflake.GenID()
	post.CreateTime = time.Now()
	post.Title = postpara.Title
	post.Content = postpara.Content
	post.AuthorId = postpara.AuthorId
	post.CommunityId = postpara.CommunityId
	post.UpdateTime = time.Now()
	if mysqlerr := mysql.GetDB().Create(&post).Error; mysqlerr != nil {
		zap.L().Error("mysql创建错误", zap.Error(mysqlerr))
		return mysqlerr
	}
	data, redisErr := redisCon.GetRedis().ZAdd(models.GetKey(models.KeyPostTime), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: post.PostId,
	}).Result()
	fmt.Println("data:", data)
	if redisErr != nil {
		zap.L().Error("redis创建数据错误", zap.Error(redisErr))
		return redisErr
	}
	//返回值
	return nil
}

// 获取社区detail值
func GetPostDetail(pid int64) (err error, author_name string, community *models.Community, post *models.Post) {
	//获取帖子
	var postItem models.Post
	if err := mysql.GetDB().Where("post_id=?", pid).First(&postItem).Error; err != nil {
		zap.L().Error("post数据库获取错误", zap.Error(err))
		return err, "", nil, nil
	}
	authorId := postItem.AuthorId
	communityId := postItem.CommunityId
	var comm models.Community
	var author models.Users
	//获取社区
	if commerr := mysql.GetDB().Where("community_id=?", communityId).First(&comm).Error; commerr != nil {
		zap.L().Error("社区数据获取错误", zap.Error(commerr))
		return err, "", nil, nil
	}
	//获取作者
	if authorerr := mysql.GetDB().Where("user_id=?", authorId).First(&author).Error; authorerr != nil {
		zap.L().Error("用户数据获取错误", zap.Error(authorerr))
		return err, "", nil, nil
	}
	return nil, author.Username, &comm, &postItem
}

// 帖子分页
func GetPostList(page, size int64) ([]*models.PostDetail, error) {
	var data []*models.PostDetail
	var posts []models.Post
	if err := mysql.GetDB().Find(&posts).Error; err != nil {
		zap.L().Error("数据库查询出错")
		return nil, err
	}
	if listerr := mysql.GetDB().Limit(int(size)).Offset(int(page - 1)).Find(&posts).Error; listerr != nil {
		zap.L().Error("数据分页错误")
		return nil, listerr
	}
	for _, post := range posts {
		var postitem models.Post
		var author models.Users
		var community models.Community
		postitem = post
		if usererr := mysql.GetDB().Where("user_id=?", post.AuthorId).First(&author).Error; usererr != nil {
			zap.L().Error("获取用户数据出错")
			continue
		}
		if commerr := mysql.GetDB().Where("community_id=?", post.CommunityId).First(&community).Error; commerr != nil {
			zap.L().Error("获取社区数据出错")
			continue
		}
		var postDetail = &models.PostDetail{author.Username, &community, &postitem}
		data = append(data, postDetail)
	}

	return data, nil
}

// 帖子投票logic
func PostVoted(userId, postId string, direction float64) (err error) {
	//首先判断帖子是否过期
	cTime, _ := redisCon.GetRedis().ZScore(models.GetKey(models.KeyPostTime), postId).Result()
	if float64(time.Now().Unix())-cTime > ExpireTime {
		return ErrorPostExpired
	}
	//查找
	oldval, _ := redisCon.GetRedis().ZScore(models.GetKey(models.KeyPostVotedPrefix+postId), userId).Result()
	diff := math.Abs(oldval - direction)
	var dir float64
	if oldval > direction {
		dir = -1
	} else {
		dir = 1
	}
	//对post进行修改
	postscroreerr := redisCon.GetRedis().ZIncrBy(models.GetKey(models.KeyPostScore), dir*diff*ScorePerVote, postId).Err()
	if postscroreerr != nil {
		return ErrorUpdatePost
	}
	if direction == 0 {
		_, reErr := redisCon.GetRedis().ZRem(models.GetKey(models.GetKey(models.KeyPostVotedPrefix+postId)), userId).Result()
		if reErr != nil {
			return reErr
		}
	} else {
		//对post-user 的score进行修改
		_, addErr := redisCon.GetRedis().ZAdd(models.GetKey(models.KeyPostVotedPrefix+postId), redis.Z{
			Score:  direction,
			Member: userId,
		}).Result()
		if addErr != nil {
			return addErr
		}
	}
	return err
}
