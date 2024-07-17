package scaffold

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/daodao97/goadmin/pkg/ecode"
)

func RenderErrMsg(c *gin.Context, code int, msg string) {
	data := map[string]interface{}{
		"code":    code,
		"message": msg,
	}
	c.JSON(http.StatusOK, data)
}

func Response(ctx *gin.Context, res interface{}, err error) {
	if err != nil {
		if code, ok := ecode.FromError(err); ok {
			RenderErrMsg(ctx, code.Code(), err.Error())
			return
		}
		RenderErrMsg(ctx, ecode.ServerErr.Code(), err.Error())
		return
	}
	data := map[string]interface{}{
		"code": 0,
		"data": res,
	}
	ctx.JSON(http.StatusOK, data)
}
