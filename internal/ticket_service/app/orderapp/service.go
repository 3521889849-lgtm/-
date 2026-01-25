// Package orderapp 承载“订单域”的应用层用例（下单/支付/退改/查询等）。
//
// 注意：订单域会依赖票务域的“座席分配”边界能力（ticketapp.AllocateSeats），
// 从而保证依赖方向清晰：order -> ticket，而不是把票务细节散落在订单实现里。
package orderapp

// Service 是订单域用例入口集合（下单/支付/退改/查询）。
type Service struct{}

// New 创建订单域应用服务。
func New() *Service {
	return &Service{}
}
