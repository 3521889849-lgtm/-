package rpc

import (
	"example_shop/kitex_gen/ticketapi/ticketservice"

	kclient "github.com/cloudwego/kitex/client"
)

func NewTicketClient(addr string) (ticketservice.Client, error) {
	return ticketservice.NewClient("ticket_service", kclient.WithHostPorts(addr))
}

