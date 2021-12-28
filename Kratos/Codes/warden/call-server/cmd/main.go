package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"call-server/internal/di"
	"git.huoys.com/middle-end/kratos/pkg/conf/paladin"
	"git.huoys.com/middle-end/kratos/pkg/log"
    //"git.huoys.com/middle-end/kratos/pkg/naming/kubernetes"
    //"git.huoys.com/middle-end/kratos/pkg/net/rpc/warden/resolver"
)

func main() {
	flag.Parse()
	log.Init(nil) // debug flag: log.dir={path}
	defer log.Close()
	log.Info("call-server start")
	paladin.Init()

	//resolver.Register(kubernetes.Builder())

	_, closeFunc, err := di.InitApp()
	if err != nil {
		panic(err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeFunc()
			log.Info("call-server exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
