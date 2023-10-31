package ecode

var (
	TokenErr   = Error(401, "接口签名错误")
	RequestErr = Error(500, "请求参数错误")
	ServerErr  = Error(500, "服务器错误")
)
