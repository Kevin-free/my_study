package di

import (
	"context"
	"time"

	"grpc-server/internal/service"

	"git.huoys.com/middle-end/kratos/pkg/log"
	kgin "git.huoys.com/middle-end/kratos/pkg/net/http/gin"
	"github.com/gin-gonic/gin"
	"git.huoys.com/middle-end/kratos/pkg/net/rpc/warden"
)

//go:generate kratos tool wire
type App struct {
	svc *service.Service
	http *gin.Engine
	grpc *warden.Server
}

func NewApp(svc *service.Service, h *gin.Engine, g *warden.Server) (app *App, closeFunc func(), err error){
	app = &App{
		svc: svc,
		http: h,
		grpc: g,
	}
	closeFunc = func() {
		ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
		if err := g.Shutdown(ctx); err != nil {
			log.Error("grpcSrv.Shutdown error(%v)", err)
		}
		if err := kgin.Shutdown(h, ctx); err != nil {
            log.Error("httpSrv.Shutdown error(%v)", err)
        }
        log.Info("httpSrv.Shutdown")
		cancel()
	}
	return
}
