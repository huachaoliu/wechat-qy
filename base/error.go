package base

import "fmt"

// 错误返回码
const (
	ErrCodeOk = 0
)

// Error 为 API 调用失败的响应内容
type Error struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("errcode: %d, errmsg: %s", e.ErrCode, e.ErrMsg)
}