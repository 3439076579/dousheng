package user

import (
	"time"
)

// User 这个结构体用来定义用户的整体信息，包括Id,名字，关注数，被关注数，以及是否关注
type User struct {
	Avatar          string `json:"avatar" gorm:"default:http://rqd659d3i.hn-bkt.clouddn.com/start_avator.png"` // 用户头像
	BackgroundImage string `json:"background_image"`                                                           // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count"`                                                             // 喜欢数
	FollowCount     int64  `json:"follow_count"`                                                               // 关注总数
	FollowerCount   int64  `json:"follower_count"`                                                             // 粉丝总数
	ID              int64  `json:"id"`                                                                         // 用户id
	IsFollow        bool   `json:"is_follow"`                                                                  // true-已关注，false-未关注
	Name            string `json:"name" gorm:"default:'初始用户名'"`                                                // 用户名称
	Signature       string `json:"signature" gorm:"default:'这个用户暂时没有签名'"`                                      // 个人简介
	TotalFavorited  string `json:"total_favorited" gorm:"default:'0'"`                                         // 获赞数量
	WorkCount       int64  `json:"work_count"`                                                                 // 作品数
	UserId          int64  `json:"-" `
}

// UserLogin 记录用户登录时的信息,Has one User
type UserLogin struct {
	ID        int64  `gorm:"primaryKey"`
	Username  string `gorm:"not null;type:varchar(32)"`
	Password  string `gorm:"not null;type:varchar(32)"`
	user      User   `gorm:"foreignKey:UserId;references:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
