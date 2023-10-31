package scaffold

import (
	"fmt"
	"github.com/daodao97/goadmin/internal/front"
	"github.com/daodao97/goadmin/pkg/cache"
	"github.com/daodao97/goadmin/pkg/db"
	"github.com/daodao97/goadmin/pkg/ecode"
	"github.com/daodao97/goadmin/pkg/log"
	"github.com/daodao97/goadmin/pkg/util"
	"github.com/daodao97/goadmin/pkg/util/uploader"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	maxSize  = 1024 * 1024 * 20 // 20 MB
	typeSvga = "svga"
)

func NewEngine(c *Conf, up uploader.Uploader, user *UserState, s Scaffold) *gin.Engine {
	if util.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}
	e := gin.Default()
	if c.HttpServer.WebPath != "" {
		e.Use(front.ModifyIndexHTML(c.HttpServer.WebPath))
		e.StaticFS(c.HttpServer.WebPath, front.Static())
	}
	root := e.Group(c.HttpServer.BasePath)
	root.GET("/ping", ping)
	root.POST("/upload", upload(up))
	root.Any("/proxy/:path", proxy)
	root.GET("/schema/*route", schemaByRoute(s))
	root.GET("/form_mutex", formMutex(s))
	e.RouterGroup.Use(Auth(c, user))
	return e
}

func schemaByRoute(s Scaffold) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		route, ok := ctx.Params.Get("route")
		if !ok {
			Response(ctx, nil, ecode.RequestErr.Message("路由错误"))
			return
		}
		if !util.String(route).StartWith("/") {
			route = "/" + route
		}
		res, err := s.GetSchemaByRoute(ctx, route)
		Response(ctx, res, err)
	}
}

func formMutex(s Scaffold) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		path := ctx.Query("path")
		if path == "" {
			Response(ctx, nil, ecode.RequestErr.Message("缺少 path 参数"))
			return
		}
		user, err := s.User(ctx)
		if err != nil {
			Response(ctx, nil, err)
			return
		}

		key := fmt.Sprintf("form_mutex:%s", path)

		currentId, err := s.Cache.Get(ctx, key)
		if err != nil && !errors.Is(cache.ErrNotFound, err) {
			Response(ctx, nil, err)
			return
		}

		idStr := cast.ToString(user.Id)
		if currentId == "" || currentId == idStr {
			err = s.Cache.SetWithTTL(ctx, key, idStr, time.Second*5)
			Response(ctx, nil, err)
			return
		}

		u := NewUser().SelectOne(db.WhereEq("id", currentId))
		if u.Err != nil && !errors.Is(u.Err, db.ErrNotFound) {
			Response(ctx, nil, u.Err)
			return
		}

		name := u.GetString("name")
		if nickname := u.GetString("nickname"); nickname != "" {
			name = nickname
		}

		Response(ctx, fmt.Sprintf("[ %s ] 正在编辑当前数据, 您暂时不能操作", name), nil)
		return
	}
}

func upload(up uploader.Uploader) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if file.Size > maxSize {
			RenderErrMsg(ctx, ecode.RequestErr.Code(), "文件过大，不支持超过20M的文件")
			return
		}
		if err != nil {
			RenderErrMsg(ctx, ecode.RequestErr.Code(), "文件为空")
			return
		}
		// parse file, get type, size
		ftype := file.Header.Get("Content-Type")
		if ftype == "application/octet-stream" && strings.Contains(file.Filename, typeSvga) {
			ftype = typeSvga
		}
		allowed := util.ArrStr{"image/jpeg", "image/png", "image/webp", "image/gif", "text/csv", typeSvga}
		if !allowed.Has(ftype) {
			log.Error("filetype not allow file type(%s)", ftype)
			RenderErrMsg(ctx, ecode.RequestErr.Code(), "不支持的文件类型 "+ftype)
			return
		}

		src, err := up.Upload(file.Filename, file)
		if err != nil {
			RenderErrMsg(ctx, ecode.RequestErr.Code(), fmt.Sprintf("文件上传失败 %s", err.Error()))
			return
		}
		Response(ctx, map[string]string{"url": src}, nil)
	}
}

func proxy(ctx *gin.Context) {
	path, exist := ctx.Params.Get("path")
	if !exist {
		RenderErrMsg(ctx, ecode.RequestErr.Code(), "必须指定转发路径")
		return
	}
	err := util.HttpProxy(ctx.Writer, ctx.Request, path)
	if err != nil {
		RenderErrMsg(ctx, ecode.RequestErr.Code(), fmt.Sprintf("接口转发失败 %v", err))
		return
	}
}

func ping(ctx *gin.Context) {
	Response(ctx, "pong", nil)
}
