package scaffold

import (
	"fmt"
	"github.com/daodao97/goadmin/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

const tokenErrCode = 401
const authErrCode = 403

func Auth(conf *Conf, userState *UserState) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		inWhite := false
		for _, v := range conf.Jwt.OpenApi {
			if util.WildcardPatternMatch(ctx.Request.RequestURI, v) {
				inWhite = true
			}
		}
		if inWhite {
			ctx.Next()
			return
		}
		token := ctx.Request.Header.Get("x-token")
		if token == "" {
			RenderErrMsg(ctx, tokenErrCode, "缺少访问凭证 X-Token")
			ctx.Abort()
			return
		}

		t := NewToken(&JwtConf{
			Secret:      conf.Jwt.Secret,
			TokenExpire: conf.Jwt.TokenExpire,
		})

		_t, err := t.ParseToken(token)
		if err != nil {
			RenderErrMsg(ctx, tokenErrCode, fmt.Sprintf("登录信息已失效 %s", err.Error()))
			ctx.Abort()
			return
		}

		exp, _ := _t.GetExpirationTime()
		if exp.Before(time.Now()) {
			RenderErrMsg(ctx, tokenErrCode, "登录信息已过期")
			ctx.Abort()
			return
		}

		isLogin, _ := userState.IsLogin(ctx, _t.UserID)
		if !isLogin {
			RenderErrMsg(ctx, tokenErrCode, "登录信息已失效")
			ctx.Abort()
			return
		}

		havePermission := false
		for _, v := range conf.Jwt.PublicApi {
			if util.WildcardPatternMatch(ctx.Request.RequestURI, v) {
				havePermission = true
			}
		}
		if havePermission {
			ctx.Next()
			return
		}

		if !userState.HavePermission(ctx, cast.ToInt(_t.UserID), ctx.Request) {
			RenderErrMsg(ctx, authErrCode, "权限不足")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
