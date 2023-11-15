package video

import (
	"awesomeProject/model/user"
	"time"
)

// Video 定义视频类型,Belong to Author
type Video struct {
	Author         user.User    `json:"author" gorm:"foreignKey:AuthorId;references:ID"`
	AuthorId       int64        `json:"-" `
	CommentCount   int64        `json:"comment_count"`         // 视频的评论总数
	CoverURL       string       `json:"cover_url"`             // 视频封面地址
	FavoriteCount  int64        `json:"favorite_count"`        // 视频的点赞总数
	ID             int64        `json:"id" gorm:"primary_key"` // 视频唯一标识
	IsFavorite     bool         `json:"is_favorite"`           // true-已点赞，false-未点赞
	PlayURL        string       `json:"play_url"`              // 视频播放地址
	Title          string       `json:"title"`                 // 视频标题
	CreatedAt      time.Time    `json:"-"`                     // 创建视频的时间
	UpdatedAt      time.Time    `json:"-"`                     // 更新视频信息的时间
	FavouriteUsers []*user.User `json:"-" gorm:"many2many:favourite_video_user"`
}
