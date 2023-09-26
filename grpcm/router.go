/*
File Name:  router.py
Description:
Author:      Chenghu
Date:       2023/8/29 22:27
Change Activity:
*/
package grpcm

import (
	"google.golang.org/grpc"
)

var RpcRouter = make([]Router, 0)

type Router struct {
	Server *grpc.ServiceDesc
	Imp    any
}

func RpcRegister(s *grpc.ServiceDesc, i any) {
	RpcRouter = append(RpcRouter, Router{
		Server: s,
		Imp:    i,
	})
}

func InitRpcRouter(s *grpc.Server) {
	for _, r := range RpcRouter {
		s.RegisterService(r.Server, r.Imp)
	}
}
