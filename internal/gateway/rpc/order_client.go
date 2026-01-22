package rpc

import (
	"example_shop/kitex_gen/orderapi/orderservice"

	kclient "github.com/cloudwego/kitex/client"
)

func NewOrderClient(addr string) (orderservice.Client, error) {
	return orderservice.NewClient("order_service", kclient.WithHostPorts(addr))
}

