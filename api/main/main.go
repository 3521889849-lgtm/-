package main

import (
	"example_shop/api/client"
	"example_shop/api/router"
	_ "example_shop/common/init"
	"log"
)

func main() {
	// 初始化 RPC 客户端
	if err := client.InitHotelClient(); err != nil {
		log.Fatalf("初始化 RPC 客户端失败: %v", err)
	}

	// 创建路由
	r := router.SetupRouter()

	// 启动服务器
	port := ":8080"
	log.Printf("🚀 HTTP API 服务启动成功，监听端口: %s", port)

	if err := r.Run(port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
