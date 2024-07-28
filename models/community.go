package models

import "time"

type Community struct {
	Id            int       `json:"id"`
	CommunityId   int       `json:"community_id"`
	CommunityName string    `json:"community_name"`
	Introduction  string    `json:"introduction"`
	CreateTime    time.Time `json:"create_time"`
	UpdateTime    time.Time `json:"update_time"`
}

func (Community) TableName() string {
	return "community"
}

type CommunityList struct {
	Id            int    `json:"id"`
	CommunityName string `json:"community_name"`
}

type CommunityDetail struct {
	Id            int       `json:"id"`
	CommunityName string    `json:"community_name"`
	Introduction  string    `json:"introduction"`
	CreateTime    time.Time `json:"create_time"`
}
