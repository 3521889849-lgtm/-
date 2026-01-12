package main

import (
	"context"
	coupon "example_shop/kitex_gen/coupon"
)

// CouponServiceImpl implements the last service interface defined in the IDL.
type CouponServiceImpl struct{}

// Test implements the CouponServiceImpl interface.
func (s *CouponServiceImpl) Test(ctx context.Context, req *coupon.EmptyReq) (resp *coupon.BaseResp, err error) {
	// TODO: Your code here...
	return
}
