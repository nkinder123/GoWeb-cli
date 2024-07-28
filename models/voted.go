package models

type Voted struct {
	PostId    int64 `json:"post_id,string" binding:"required"`
	Direction int8  `json:"direction" binding:"oneof=1 0 -1"`
}
