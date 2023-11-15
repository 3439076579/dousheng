package common

// Response Response结构体，包括状态码以及返回信息的描述，如果状态码为0则可以省略返回状态描述
// 被用于其他返回报文的基本结构体
type Response struct {
	StatusCode int64  `json:"status_code"`          // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg,omitempty"` // 返回状态描述
}
