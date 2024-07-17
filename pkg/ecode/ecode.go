package ecode

import "errors"

func Error(num int, msg string) *Code {
	return &Code{
		errMsg:  msg,
		errCode: num,
	}
}

type Code struct {
	errMsg  string
	errCode int
}

func (e *Code) Error() string {
	return e.errMsg
}

func (e *Code) Code() int {
	return e.errCode
}

func (e *Code) Message(msg string) *Code {
	e.errMsg = msg
	return e
}

// FromError 尝试将 error 转换为 *Code 类型
func FromError(err error) (*Code, bool) {
	if err == nil {
		return nil, false
	}
	// 类型断言
	var codeErr *Code
	if errors.As(err, &codeErr) {
		return codeErr, true
	}
	return nil, false
}
