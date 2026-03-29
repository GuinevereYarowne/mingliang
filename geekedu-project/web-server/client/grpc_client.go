package client

//gRPC 客户端，连接 Logic Server，提供调用接口

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"geekedu-project/proto"
)

var GeekEduClient proto.GeekEduServiceClient

// 初始化gRPC客户端
func InitGRPCClient() {
	// 连接Logic Server（容器名：logic-server，端口50051）
	conn, err := grpc.Dial("logic-server:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接gRPC服务失败：%v", err)
	}
	log.Println("gRPC客户端连接成功")

	// 创建客户端
	GeekEduClient = proto.NewGeekEduServiceClient(conn)
}
