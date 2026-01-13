package coupon

import (
	"example_shop/kitex_gen/coupon"
)

type CouponService struct{}

func (s *CouponService) AddCoupon() (*coupon.BaseResp, error) {
	return &coupon.BaseResp{Code: 200, Msg: "添加成功"}, nil
}
