package ecode

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
