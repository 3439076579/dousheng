package response

import (
	"awesomeProject/model/common"
	"awesomeProject/model/user"
)

type ResponseComment struct {
	CommentID int64     `json:"id"`
	User      user.User `json:"user"`
	Content   string    `json:"content"`
	Created   string    `json:"create_date"`
}

type ResponsePostComment struct {
	common.Response
	Comment ResponseComment
}
