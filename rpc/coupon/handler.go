package coupon

import (
	"context"
	"example_shop/kitex_gen/coupon"
)

type CouponService struct{}

func (s *CouponService) AddCoupon(ctx context.Context, req *coupon.AddCouponReq) (*coupon.BaseResp, error) {
	return &coupon.BaseResp{Code: 200, Msg: "添加成功"}, nil
}
