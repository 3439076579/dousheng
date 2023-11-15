package utils

import "errors"

var (
	RecordNotFound        = errors.New("record not found")                   //查询不到记录时的错误
	EmailIllegal          = errors.New("email is incorrect")                 //邮箱不合法
	PasswordIllegal       = errors.New("password is incorrect")              //密码不正确
	PassWordOrUserNameErr = errors.New("password or username has been used") //密码或账号已被使用
	ParameterErr          = errors.New("too many parameters")
	CommentError          = errors.New("comment has been removed")
	CancelFavoriteErr     = errors.New("cannot cancel favorite cause has not favorite")
	DuplicateAddFavorite  = errors.New("cannot add favorite cause has been favorite")
	InvalidFavoriteAction = errors.New("invalid favorite action")
)
