package main

import (
	"example_shop/common/config"
	initpkg "example_shop/common/init"
	"example_shop/internal/gateway/http/handler"
	"example_shop/internal/gateway/http/router"
	"example_shop/internal/gateway/rpc"
	"example_shop/internal/ticket_service/handler/ticketapi"
	"example_shop/internal/ticket_service/job"
	"example_shop/internal/ticket_service/handler/userapi"
	"example_shop/internal/user_service/handler/orderapi"
	"example_shop/kitex_gen/orderapi/orderservice"
	"example_shop/kitex_gen/ticketapi/ticketservice"
	"example_shop/kitex_gen/userapi/userservice"
	"fmt"
	"log"
	"net"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	kserver "github.com/cloudwego/kitex/server"
)

func pickAvailable(host string, port int) string {
	for i := 0; i <= 20; i++ {
		addr := fmt.Sprintf("%s:%d", host, port+i)
		ln, err := net.Listen("tcp", addr)
		if err == nil {
			_ = ln.Close()
			return addr
		}
	}
	return fmt.Sprintf("%s:%d", host, port)
}

func main() {
	initpkg.Init()

	stop := make(chan struct{})
	job.StartOrderCleanup(stop)

	userListen := fmt.Sprintf("%s:%d", config.Cfg.Server.UserService.Host, config.Cfg.Server.UserService.Port)
	ticketListen := fmt.Sprintf("%s:%d", config.Cfg.Server.TicketService.Host, config.Cfg.Server.TicketService.Port)
	orderListen := fmt.Sprintf("%s:%d", config.Cfg.Server.OrderService.Host, config.Cfg.Server.OrderService.Port)
	gatewayListen := pickAvailable(config.Cfg.Server.Gateway.Host, config.Cfg.Server.Gateway.Port)

	userSvr := userservice.NewServer(
		userapi.NewUserServiceImpl(),
		kserver.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "user_service"}),
		kserver.WithServiceAddr(utils.NewNetAddr("tcp", userListen)),
	)
	ticketSvr := ticketservice.NewServer(
		ticketapi.NewTicketServiceImpl(),
		kserver.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "ticket_service"}),
		kserver.WithServiceAddr(utils.NewNetAddr("tcp", ticketListen)),
	)
	orderSvr := orderservice.NewServer(
		orderapi.NewOrderServiceImpl(),
		kserver.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "order_service"}),
		kserver.WithServiceAddr(utils.NewNetAddr("tcp", orderListen)),
	)

	go func() { log.Println("user_service started"); log.Println("user_service stopped:", userSvr.Run()) }()
	go func() { log.Println("ticket_service started"); log.Println("ticket_service stopped:", ticketSvr.Run()) }()
	go func() { log.Println("order_service started"); log.Println("order_service stopped:", orderSvr.Run()) }()

	userClient, err := rpc.NewUserClient(userListen)
	if err != nil {
		log.Fatalf("user_service client 初始化失败: %v", err)
	}
	ticketClient, err := rpc.NewTicketClient(ticketListen)
	if err != nil {
		log.Fatalf("ticket_service client 初始化失败: %v", err)
	}
	orderClient, err := rpc.NewOrderClient(orderListen)
	if err != nil {
		log.Fatalf("order_service client 初始化失败: %v", err)
	}

	h := server.Default(server.WithHostPorts(gatewayListen))
	app := handler.NewApp(userClient, ticketClient, orderClient)
	router.RegisterRoutes(h, app)
	log.Printf("dev all-in-one started at %s", gatewayListen)
	h.Spin()

	close(stop)
}
