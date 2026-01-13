package main

import (
	_ "example_shop/common/init"
	coupon "example_shop/kitex_gen/coupon/couponservice"
	"log"
)

func main() {
	svr := coupon.NewServer(new(CouponServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
