// Package ticket 是网关内的“票务业务模块”（Application/UseCase）。
//
// 目标：
// - 将票务查询类逻辑（站点联想、车次查询、缓存/熔断策略）从 HTTP Handler 下沉出来
// - 让 handler 只做“参数绑定 -> 调用用例 -> 写响应”，避免业务代码散落在网关各处
package ticket

// Service 是网关“票务模块”的用例入口集合（站点联想、车次查询等）。
type Service struct{}

// New 创建票务业务模块实例。
func New() *Service {
	return &Service{}
}

// Result 是网关业务模块返回给 HTTP Handler 的统一结果。
// Status 为 HTTP 状态码，Body 为要输出的 JSON 结构体。
type Result struct {
	Status int
	Body   any
}
