package http

import (
	"net/http"

	pb "grpc-server/api"
	"grpc-server/internal/model"
	"git.huoys.com/middle-end/kratos/pkg/conf/paladin"
	"git.huoys.com/middle-end/kratos/pkg/log"
	cfg "git.huoys.com/middle-end/kratos/pkg/net/http/config"
	kgin "git.huoys.com/middle-end/kratos/pkg/net/http/gin"
	"github.com/gin-gonic/gin"
)

var svc pb.DemoServer

// New new a bm server.
func New(s pb.DemoServer) (engine *gin.Engine, err error) {
	var (
		cfg cfg.ServerConfig
		ct paladin.TOML
	)
	if err = paladin.Get("http.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Server").UnmarshalTOML(&cfg); err != nil {
		return
	}
	svc = s
	engine = kgin.DefaultServer(&cfg)
	pb.RegisterDemoGinServer(engine, s)
	initRouter(engine)
	kgin.Start(engine)
	return
}

func initRouter(e *gin.Engine) {
	e.GET("/ping", ping)
	g := e.Group("/grpc-server")
	{
		g.GET("/start", howToStart)
	}
}

func ping(ctx *gin.Context) {
	if _, err := svc.Ping(ctx, nil); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// example for http request handler.
func howToStart(c *gin.Context) {
	k := &model.Kratos{
		Hello: "Golang 大法好 !!!",
	}
	c.JSON(200, k)
}