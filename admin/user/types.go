package user

type LoginReq struct {
	// login form data
	Username string `form:"username"`
	Password string `form:"password"`
	Captcha  string `form:"captcha"`
	Sing     string `form:"sing"`

	// sso login data
	Ticket string `form:"ticket"`
	Key    string `form:"key"`
}

type LoginRes struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type MngVerifyRes struct {
	Code     int    `json:"code"`
	Username string `json:"username"`
	Message  string `json:"message"`
}
