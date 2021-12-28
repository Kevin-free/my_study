package grpc

import (
	pb "call-server/api"

	"git.huoys.com/middle-end/kratos/pkg/conf/paladin"
	"git.huoys.com/middle-end/kratos/pkg/net/rpc/warden"
)

// New new a grpc server.
func New(svc pb.DemoServer) (ws *warden.Server, err error) {
	var (
		cfg *warden.ServerConfig
		ct paladin.TOML
	)

	if err = paladin.Get("grpc.toml").Unmarshal(&ct); err != nil {
        cfg = nil
    } else if err = ct.Get("Server").UnmarshalTOML(cfg); err != nil {
        cfg = nil
    }

	ws = warden.NewServer(cfg)
	pb.RegisterDemoServer(ws.Server(), svc)
	ws, err = ws.Start()
	return
}
