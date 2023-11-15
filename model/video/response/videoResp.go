package response

import (
	"awesomeProject/model/common"
	"awesomeProject/model/user"
	"awesomeProject/model/video"
)

// ResponseGetVideo 定义返回视频时的报文
type ResponseGetVideo struct {
	common.Response
	NextTime  int64         `json:"next_time"`  // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	VideoList []video.Video `json:"video_list"` // 视频列表
}

type ResponseGetPublishedVideo struct {
	common.Response
	VideoList []video.Video
}

type ResponseComment struct {
	CommentID int64     `json:"id"`
	User      user.User `json:"user"`
	Content   string    `json:"content"`
	Created   string    `json:"create_date"`
}

type ResponsePostComment struct {
	common.Response
	Comment *ResponseComment
}

type ResponseGetCommentList struct {
	common.Response
	CommentList []ResponseComment `json:"comment_list"`
}
