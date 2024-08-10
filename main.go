package main

import (
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/util/log"
	"meowcloud-action/common/config"
	"meowcloud-action/controller"
	"net"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"github.com/xh-polaris/gopkg/kitex/middleware"
	logx "github.com/xh-polaris/gopkg/util/log"
	action "github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowcloud/action/actionservice"
)

func main() {

	config.Init()

	klog.SetLogger(logx.NewKlogLogger())

	addr, err := net.ResolveTCPAddr("tcp", config.Get().ListenOn)

	if err != nil {
		panic(err)
	}
	svr := action.NewServer(
		controller.NewActionController(),
		server.WithServiceAddr(addr),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.Get().Name}),
		server.WithMiddleware(middleware.LogMiddleware(config.Get().Name)),
	)

	err = svr.Run()

	if err != nil {
		log.Error(err.Error())
	}
}
