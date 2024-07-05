package app

import (
	"admin/internal/api"
	"admin/internal/conf"
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/daodao97/goadmin/pkg/util/uploader"
	"github.com/daodao97/goadmin/scaffold"
)

type HttpOptions struct {
	Conf         *conf.Conf
	ScaffoldConf *scaffold.Conf
	Ctrls        []scaffold.Ctrl
	Uploader     uploader.Uploader
	UserState    *scaffold.UserState
	Scaffold     scaffold.Scaffold
}

func NewHttpServer(opt *HttpOptions) (*HttpServer, func(), error) {
	h := &HttpServer{options: opt}
	err := h.Start()
	if err != nil {
		return nil, nil, err
	}
	return h, func() {
		h.Stop()
	}, nil
}

type HttpServer struct {
	options *HttpOptions
	server  *http.Server
}

func (h *HttpServer) Start() error {
	e := scaffold.NewEngine(h.options.ScaffoldConf, h.options.Uploader, h.options.UserState, h.options.Scaffold)

	api.Route(e, h.options.Conf)

	// 后台业务 路由注册
	for _, v := range h.options.Ctrls {
		v.Route(e)
	}
	h.server = &http.Server{
		Addr:           ":",
		Handler:        e,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := copier.Copy(h.server, h.options.ScaffoldConf.HttpServer)
	if err != nil {
		return err
	}
	if h.server.Addr[0] == ':' {
		h.server.Addr = "127.0.0.1" + h.server.Addr
	}

	go func() {
		if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	go func() {
		for {
			if _, err := net.Dial("tcp", h.server.Addr); err == nil {
				fmt.Printf("GoAdmin Server Start @ http://%s%s\n", h.server.Addr, h.options.ScaffoldConf.HttpServer.WebPath)
				break
			}
			time.Sleep(time.Second)
		}
	}()

	return nil
}

func (h *HttpServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	now := time.Now()
	err := h.server.Shutdown(ctx)
	if err != nil {
		return
	}
	fmt.Println("http server shutdown, use time ", time.Since(now))
}
