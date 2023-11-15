package request

// RequestUserLogin 用于用户登录和用户注册时的请求结构体
type RequestUserLogin struct {
	UserName string
	PassWord string
}

// RequestGetUserInfo 获取用户信息需要token和user_id
// token交给middleware去处理即可,所以该requestGetUserInfo中只有user_id一个field
type RequestGetUserInfo struct {
	UserId int64
}
