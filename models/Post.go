package models

import "time"

// 新建帖子model
type Post struct {
	Id          int64     `json:"id"`
	PostId      int64     `json:"postId"`
	Title       string    `json:"title" binding:"required"`
	Content     string    `json:"content" gorm:"type:text" binding:"required"`
	AuthorId    int64     `json:"author_id"`
	CommunityId int       `json:"community_id" binding:"required"`
	Status      int       `json:"status" gorm:"type:tinyint;default:0"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time" gorm:"type:datetime(0)"`
}

func (Post) TableName() string {
	return "post"
}

// 帖子传入参数
type PostPara struct {
	Title       string `json:"title"`
	Content     string `json:"content" gorm:"type:text"`
	CommunityId int    `json:"community_id"`
	AuthorId    int64  `json:"author_id"`
}

// 获取帖子的参数
type PostDetail struct {
	AuthorName string `json:"author_name"`
	*Community `json:"community"`
	*Post
}
