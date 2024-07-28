package models

type Users struct {
	Id         int    `json:"id"`
	UserId     int64  `json:"user_id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Gender     int    `json:"gender"`
	CreateTime int    `json:"create_time"`
	UpdateTime int    `json:"update_time"`
}

func (Users) Tablename() string { return "bluebell.users" }

// 注册校验验证
type UserConfirm struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// 登陆校验验证
type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
