package response

import (
	"awesomeProject/model/common"
	"awesomeProject/model/user"
)

// ResponseUserLogin
// 用于/douyin/user/login/ 和 /douyin/user/register的返回报文
type ResponseUserLogin struct {
	Response common.Response
	Token    string `json:"token"`             // 用户鉴权token
	UserID   int64  `json:"user_id,omitempty"` // 用户id
}

// ResponseGetUserInfo 请求用户信息时返回的JSON字符串结构体
type ResponseGetUserInfo struct {
	Response common.Response
	User     user.User
}
