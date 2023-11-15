package utils

import (
	"awesomeProject/model/interactor"
	"awesomeProject/model/user"
	"bytes"
	"runtime"
	"strconv"
	"strings"
	"time"
)

/* key的组装
UserID,CreatedAt,Content
*/

func ParseStringSliceToCommentList(cmd []string) []interactor.Comment {
	var CommentList []interactor.Comment

	for i := 0; i < len(cmd); i++ {

		var comment interactor.Comment
		StringSlice := strings.Split(cmd[i], ",")
		comment.UserID, _ = strconv.ParseInt(StringSlice[0], 10, 64)
		comment.Content = StringSlice[2]

		location, _ := time.LoadLocation("Local")
		comment.CreatedAt, _ = time.ParseInLocation("2006-01-02 15:04:05", StringSlice[1], location)

		CommentList = append(CommentList, comment)

	}
	return CommentList
}

/*
	key:UserLogin:UserName
	val:UserID:Password

*/

func ParseStringToUserLogin(value string) user.UserLogin {
	var userLoginModel user.UserLogin

	stringSlice := strings.Split(value, ":")

	userLoginModel.ID, _ = strconv.ParseInt(stringSlice[0], 10, 64)
	userLoginModel.Password = stringSlice[1]

	return userLoginModel

}

func GetGoID() int64 {

	GoInfo := make([]byte, 64)

	_ = runtime.Stack(GoInfo, false)

	GoInfo = bytes.TrimPrefix(GoInfo, []byte("goroutine "))
	GoInfo = GoInfo[:bytes.IndexByte(GoInfo, ' ')]
	n, _ := strconv.ParseInt(string(GoInfo), 10, 64)

	return n
}
