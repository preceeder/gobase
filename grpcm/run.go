/*
File Name:  run.go
Description:
Author:      Chenghu
Date:       2023/8/31 09:57
Change Activity:
*/
package grpcm

import (
	"fmt"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/preceeder/gobase/utils"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"os"
	"time"
)

type RpcConfig struct {
	Addr string `json:"addr"`
	Name string `json:"name"`
}

func redaRpcConfig(c viper.Viper) (rc *RpcConfig) {
	rpcConfig := c.Sub("rpc")

	if rpcConfig == nil {
		fmt.Printf("rpc config is nil")
		os.Exit(1)
	}
	err := rpcConfig.Unmarshal(&rc)
	if err != nil {
		fmt.Printf("rpc config read error: " + err.Error())
		os.Exit(1)
	}
	return
}

//var rpcListner net.Listener

var rcpConfig *RpcConfig

func InitRpc(config viper.Viper) {
	ur := RpcConfig{}
	utils.ReadViperConfig(config, "rpc", &ur)
	rcpConfig = &ur
}

func Server(server *grpc.Server) {
	// 创建 Tcp 连接
	if rcpConfig == nil {
		slog.Error("rpc 还未初始化 请先调用 gobase.Init(rpc:true)")
	}
	rpcLi, err := net.Listen("tcp", rcpConfig.Addr)
	if err != nil {
		slog.Error("监听失败: %v", "error", err.Error())
	}
	slog.Info("开启监听： ", "addr", rcpConfig.Addr)
	//开启信号监听
	c := utils.StartSignalLister()

	//开启信号处理
	go utils.SignalHandler(c, func() {
		//平滑关闭
		server.GracefulStop()
	})

	//初始化 注册路由
	InitRpcRouter(server)
	// 在 gRPC 服务上注册反射服务
	reflection.Register(server)

	err = server.Serve(rpcLi)
	if err != nil {
		slog.Error("failed to serve: %v", err)
	}

}

func Client(interceptor ...grpc.UnaryClientInterceptor) *grpc.ClientConn {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  time.Second * 1,
				Multiplier: 1.6,
				MaxDelay:   time.Second * 15,
			},
			MinConnectTimeout: time.Second * 15,
		}),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
			grpc_retry.WithMax(10),
			grpc_retry.WithBackoff(grpc_retry.BackoffExponential(1*time.Second)),
		)),
		grpc.WithChainUnaryInterceptor(interceptor...),
	}
	clt, err := grpc.Dial(rcpConfig.Addr, opts...)
	if err != nil {
		slog.Error("连接 gPRC 服务失败", "error", err)
	}
	return clt
}
