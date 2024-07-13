package front

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed static/*
var static embed.FS

func Static() http.FileSystem {
	assets, _ := fs.Sub(static, "static")
	return http.FS(assets)
}

func ModifyIndexHTML(webPath string, data map[string]any) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path != webPath+"/" {
			return
		}
		indexHtml, err := renderIndex(data)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to render index.html")
			ctx.Abort()
			return
		}
		// 设置修改后的内容回响应体
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.String(http.StatusOK, string(indexHtml))
		ctx.Abort()
	}
}

func renderIndex(data map[string]any) ([]byte, error) {
	fileData, err := fs.ReadFile(static, "static/index.html")
	if err != nil {
		return nil, err
	}

	tmpl, _ := template.New("index").Funcs(template.FuncMap{
		"js": func(s string) template.JS {
			return template.JS(s)
		},
	}).Parse(string(fileData))

	var buf1 bytes.Buffer
	if err := tmpl.Execute(&buf1, data); err != nil {
		return nil, err
	}
	return buf1.Bytes(), nil
}
