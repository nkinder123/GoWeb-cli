package models

const (
	KeyPrefix          = "blubell:"    //zset:使用命名空间命名
	KeyPostTime        = "post:time"   //zset：帖子的发表时间
	KeyPostScore       = "post:score"  //zset：帖子的分数
	KeyPostVotedPrefix = "post:voted:" //zset：帖子的支持量；后接post_id
)

func GetKey(str string) string {
	return KeyPrefix + str
}
