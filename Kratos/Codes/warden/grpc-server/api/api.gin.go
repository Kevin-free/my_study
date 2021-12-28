// Code generated by protoc-gen-gin v0.7, DO NOT EDIT.
// source: api.proto

/*
Package api is a generated gin stub package.
This code was generated with kratos/tool/protobuf/protoc-gen-gin v0.7.

package 命名使用 {appid}.{version} 的方式, version 形如 v1, v2 ..

It is generated from these files:
	api.proto
*/
package api

import (
	"context"

	kgin "git.huoys.com/middle-end/kratos/pkg/net/http/gin"
	"github.com/gin-gonic/gin"
	"net/http"
)
import types "github.com/gogo/protobuf/types"

// to suppressed 'imported but not used warning'
var _ *gin.Context
var _ context.Context

var PathDemoPing = "/demo.service.v1.Demo/Ping"
var PathDemoSayHello = "/demo.service.v1.Demo/SayHello"
var PathDemoSayHelloURL = "/kratos-demo/say_hello"

// DemoGinServer is the server API for Demo service.
type DemoGinServer interface {
	Ping(ctx context.Context, req *types.Empty) (resp *types.Empty, err error)

	SayHello(ctx context.Context, req *HelloReq) (resp *types.Empty, err error)

	SayHelloURL(ctx context.Context, req *HelloReq) (resp *HelloResp, err error)
}

var DemoSvc DemoGinServer

func demoPing(c *gin.Context) {
	p := new(types.Empty)
	if err := c.Bind(p); err != nil {
		return
	}
	resp, err := DemoSvc.Ping(c, p)
	res := kgin.TOJSON(resp, err)
	c.JSON(http.StatusOK, res)
}

func demoSayHello(c *gin.Context) {
	p := new(HelloReq)
	if err := c.Bind(p); err != nil {
		return
	}
	resp, err := DemoSvc.SayHello(c, p)
	res := kgin.TOJSON(resp, err)
	c.JSON(http.StatusOK, res)
}

func demoSayHelloURL(c *gin.Context) {
	p := new(HelloReq)
	if err := c.Bind(p); err != nil {
		return
	}
	resp, err := DemoSvc.SayHelloURL(c, p)
	res := kgin.TOJSON(resp, err)
	c.JSON(http.StatusOK, res)
}

// RegisterDemoGinServer Register the gin route
func RegisterDemoGinServer(e *gin.Engine, server DemoGinServer) {
	DemoSvc = server
	e.GET("/demo.service.v1.Demo/Ping", demoPing)
	e.GET("/demo.service.v1.Demo/SayHello", demoSayHello)
	e.GET("/kratos-demo/say_hello", demoSayHelloURL)
}
