package client

import (
	"example_shop/kitex_gen/hotel"
	"example_shop/kitex_gen/hotel/hotelservice"
	"log"

	"github.com/cloudwego/kitex/client"
)

var HotelClient hotelservice.Client

// InitHotelClient 初始化酒店服务 RPC 客户端
func InitHotelClient() error {
	var err error
	// 使用服务发现或直接指定地址
	// 这里使用默认配置，实际项目中应该从配置中心获取
	HotelClient, err = hotelservice.NewClient(
		"hotel_service",
		client.WithHostPorts("127.0.0.1:8888"), // RPC 服务地址
	)
	if err != nil {
		log.Printf("初始化酒店服务客户端失败: %v", err)
		return err
	}
	log.Println("✅ 酒店服务 RPC 客户端初始化成功")
	return nil
}

// 辅助函数：转换错误响应
func handleRPCError(resp *hotel.BaseResp, err error) (int, string) {
	if err != nil {
		return 500, "RPC 调用失败: " + err.Error()
	}
	if resp.Code != 200 {
		return int(resp.Code), resp.Msg
	}
	return 200, resp.Msg
}
