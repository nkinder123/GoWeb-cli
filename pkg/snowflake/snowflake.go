package snowflake

import (
	//给包取别名
	sf "github.com/bwmarrin/snowflake"
	"time"
)

var node *sf.Node

// 将相对日期、机器的ID传入初始化雪花算法
func Init(startTime string, mechineId int64) (err error) {
	var str time.Time
	str, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	//设置开始值--雪花算法：time+mechine+id（创建的时候id1，id2.。。 ）
	sf.Epoch = str.UnixNano()
	node, err = sf.NewNode(mechineId)
	return
}

// 将雪花算法生成的值转换为Int64的格式
func GenID() int64 {
	return node.Generate().Int64()
}
