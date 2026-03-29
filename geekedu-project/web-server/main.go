package main

//主程序入口，初始化数据库连接，启动 gRPC 服务

import (
	"log"

	"geekedu-project/web-server/client"
	"geekedu-project/web-server/router"
)

func main() {
	// 1. 初始化gRPC客户端
	client.InitGRPCClient()

	// 2. 初始化路由
	r := router.InitRouter()

	// 3. 启动HTTP服务
	log.Println("Web服务启动成功，端口：8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("启动Web服务失败：%v", err)
	}
}
