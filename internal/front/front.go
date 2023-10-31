package front

import (
	"embed"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed static/*
var static embed.FS

func Static() http.FileSystem {
	assets, _ := fs.Sub(static, "static")
	return http.FS(assets)
}

func ModifyIndexHTML(basePath string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path != basePath+"/" {
			return
		}
		// 读取 index.html 文件
		indexContent, err := static.ReadFile("static/index.html")
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to read index.html")
			ctx.Abort()
			return
		}

		indexContentStr := string(indexContent)
		indexContentStr = strings.Replace(indexContentStr, "{{ base_path }}", basePath, -1)

		// 设置修改后的内容回响应体
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.String(http.StatusOK, indexContentStr)
		ctx.Abort()
	}
}
