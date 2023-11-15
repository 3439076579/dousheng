package request

// RequestFavourite 用于接收赞请求的请求结构体
type RequestFavourite struct {
	VideoID     int64
	ActionType  string
	UserLoginID int64
}

type RequestGetFavouriteList struct {
	UserLoginID int64
}

// RequestPostComment 用于接收评论操作的请求结构体
type RequestPostComment struct {
	UserID     int64
	VideoID    int64
	ActionType string
	Comment    string // only use in action_type=1
	CommentID  int64  // only use in action_type=2
}

type RequestGetComment struct {
	VideoID int64
	UserID  int64
}
