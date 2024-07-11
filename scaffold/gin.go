package scaffold

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/daodao97/xgo/xlog"
	"github.com/pkg/errors"
	"github.com/spf13/cast"

	"github.com/daodao97/goadmin/internal/front"
	"github.com/daodao97/goadmin/pkg/cache"
	"github.com/daodao97/goadmin/pkg/db"
	"github.com/daodao97/goadmin/pkg/ecode"
	"github.com/daodao97/goadmin/pkg/sso"
	"github.com/daodao97/goadmin/pkg/util"
	"github.com/daodao97/goadmin/pkg/util/uploader"

	"github.com/gin-gonic/gin"
)

const (
	maxSize  = 1024 * 1024 * 20 // 20 MB
	typeSvga = "svga"
)

type GinRoute = func(e *gin.Engine)

func RegRoute(e *gin.Engine, routes []GinRoute) {
	for _, route := range routes {
		route(e)
	}
}

type EngineOption struct {
	Conf     *Conf
	Uploader uploader.Uploader
	Cache    cache.Cache
	Sso      *sso.Sso
}

func NewEngine(opt *EngineOption) *gin.Engine {
	if opt.Sso == nil {
		opt.Sso = &map[sso.Name]sso.SSO{}
	}
	if opt.Cache == nil {
		opt.Cache = cache.NewMemoryCache()
	}

	if opt.Uploader == nil {
		opt.Uploader = uploader.NewLocalUploader(
			"./uploads",
			"ok/{year}{month}{day}-{hour}{minute}{second}-{random}{.suffix}",
			"http://127.0.0.1:8001",
		)
	}

	us, _ := NewUserState(opt.Cache, opt.Conf)

	s := New(&Options{
		Conf:      opt.Conf,
		UserState: us,
		Cache:     opt.Cache,
		Sso:       opt.Sso,
	})

	c := opt.Conf
	user := us
	up := opt.Uploader

	if util.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	logger := xlog.GetLogger()

	// 创建一个自定义的 SlogWriter
	slogWriter := SlogWriter{logger: logger}

	// 设置 gin 使用自定义的日志记录器
	gin.DefaultWriter = slogWriter
	gin.DefaultErrorWriter = slogWriter

	e := gin.New()

	// 使用自定义的日志中间件
	e.Use(Cors(), func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		latency := end.Sub(start)

		// 使用 slog 记录结构化日志
		logger.Info("HTTP Request",
			slog.String("client_ip", c.ClientIP()),
			slog.String("time", end.Format(time.RFC1123)),
			slog.Int("status_code", c.Writer.Status()),
			slog.String("latency", latency.String()),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("error_message", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		)
	}, gin.Recovery())

	if c.HttpServer.WebPath != "" {
		e.Use(front.ModifyIndexHTML(c.HttpServer.WebPath, map[string]any{
			"basePath": c.HttpServer.BasePath,
			"webPath":  c.HttpServer.WebPath,
		}))
		e.StaticFS(c.HttpServer.WebPath, front.Static())
	}
	root := e.Group(c.HttpServer.BasePath).Use(Auth(c, user))
	root.GET("/ping", ping)
	root.POST("/upload", upload(up))
	root.Any("/proxy/:path", proxy)
	root.GET("/schema/*route", schemaByRoute(s))
	root.GET("/form_mutex", formMutex(s))

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

type SlogWriter struct {
	logger *slog.Logger
}

func (w SlogWriter) Write(p []byte) (n int, err error) {
	w.logger.Info(string(p))
	return len(p), nil
}
