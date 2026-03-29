package main

//主程序入口，初始化数据库连接，启动 gRPC 服务

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"geekedu-project/logic-server/db"
	"geekedu-project/logic-server/service"
	"geekedu-project/proto"
)

func main() {
	// 1. 初始化数据库
	if err := db.InitDB(); err != nil {
		log.Fatalf("初始化数据库失败：%v", err)
	}
	log.Println("数据库初始化成功")

	// 2. 启动gRPC服务
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("监听端口失败：%v", err)
	}
	grpcServer := grpc.NewServer()

	// 注册服务
	proto.RegisterGeekEduServiceServer(grpcServer, &service.GeekEduService{})

	log.Println("gRPC服务启动成功，端口：50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("启动gRPC服务失败：%v", err)
	}
}
